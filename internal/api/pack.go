package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gptlv/re-partners/packs/internal/app"
	"github.com/gptlv/re-partners/packs/internal/view"
	"github.com/gptlv/re-partners/packs/pkg/calculate"
)

var ErrInvalidAmount = errors.New("invalid amount")

type Handler struct {
	service  *app.Service
	renderer view.Renderer
}

func NewHandler(service *app.Service, renderer view.Renderer) *Handler {
	return &Handler{
		service:  service,
		renderer: renderer,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.service.Sizes(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "service error")
		return
	}

	payload := IndexViewModel{
		Sizes: sizes,
	}

	err = h.renderer.Render(w, "index", payload)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "template error")
	}
}

func (h *Handler) CalculateHTML(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid form")
		return
	}

	amount, err := parseAmount(r.FormValue("amount"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "amount must be a positive integer")
		return
	}

	packs, err := h.service.CalculatePackages(r.Context(), amount)
	if err != nil {
		if errors.Is(err, calculate.ErrCannotFulfill) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`<div>Cannot fulfill the order with current packs.</div>`))
			return
		}
		h.respondError(w, http.StatusInternalServerError, "calculation failed")
		return
	}

	payload := CalculateHTMLResponse{
		Amount: amount,
		Packs:  toAPIPacks(packs),
	}

	err = h.renderer.Render(w, "result", payload)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "template error")
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

func parseAmount(raw string) (int64, error) {
	amount, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || amount <= 0 {
		return 0, ErrInvalidAmount
	}
	return amount, nil
}
