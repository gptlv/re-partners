package api

import (
	"encoding/json"
	"errors"
	"net/http"

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
	sizes, err := h.service.Sizes(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "database error")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]any{"packs": sizes})
}

func (h *Handler) CalculateJSON(w http.ResponseWriter, r *http.Request) {
	var req CalculateJSONRequest

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

	resp := CalculateJSONResponse{
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
