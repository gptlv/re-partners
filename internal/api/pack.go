package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gptlv/re-partners/packs/internal/app"
	"github.com/gptlv/re-partners/packs/pkg/calculate"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetPacks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, http.MethodGet)
		return
	}

	sizes, err := h.service.Sizes(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "database error")
		return
	}

	resp := PackSizesResponse{
		Packs: toPackSizeResponses(sizes),
	}

	h.writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) CreateSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, http.MethodPost)
		return
	}

	var req CreateSizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Size <= 0 {
		h.respondError(w, http.StatusBadRequest, "pack size must be greater than 0")
		return
	}

	created, err := h.service.AddSize(r.Context(), req.Size)
	if err != nil {
		switch {
		case errors.Is(err, app.ErrSizeExists):
			h.respondError(w, http.StatusConflict, err.Error())
			return
		default:
			h.respondError(w, http.StatusInternalServerError, "database error")
			return
		}
	}

	resp := CreateSizeResponse{
		ID:   created.ID,
		Size: created.Size,
	}

	h.writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) DeleteSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.methodNotAllowed(w, http.MethodDelete)
		return
	}

	const prefix = "/api/sizes/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		h.respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	rawID := strings.TrimPrefix(r.URL.Path, "/api/sizes/")
	if rawID == "" || strings.Contains(rawID, "/") {
		h.respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		h.respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.DeleteSize(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, app.ErrLastSize):
			h.respondError(w, http.StatusUnprocessableEntity, err.Error())
			return
		case errors.Is(err, app.ErrSizeNotFound):
			h.respondError(w, http.StatusNotFound, err.Error())
			return
		default:
			h.respondError(w, http.StatusInternalServerError, "database error")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CalculateJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, http.MethodPost)
		return
	}

	var req CalculateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Amount <= 0 {
		h.respondError(w, http.StatusBadRequest, "amount must be a positive integer")
		return
	}

	packs, err := h.service.CalculatePackages(r.Context(), req.Amount)
	if err != nil {
		if errors.Is(err, calculate.ErrCannotFulfill) {
			h.respondError(w, http.StatusUnprocessableEntity, "cannot fulfill order")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "calculation failed")
		return
	}

	resp := CalculateResponse{
		Amount: req.Amount,
		Packs:  toAPIPacks(packs),
	}

	h.writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (h *Handler) methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
