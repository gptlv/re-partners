package router

import (
	"net/http"

	"github.com/gptlv/re-partners/packs/internal/api"
)

// New constructs the HTTP mux wiring handlers to routes.
func New(handler *api.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/ui/calc", handler.CalculateHTML)
	mux.HandleFunc("/api/packs", handler.GetPacks)
	mux.HandleFunc("/api/calc", handler.CalculateJSON)
	return mux
}
