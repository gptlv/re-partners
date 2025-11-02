package router

import (
	"net/http"

	"github.com/gptlv/re-partners/packs/internal/api"
)

// New constructs the HTTP mux wiring handlers to routes.
func New(handler *api.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	assets := http.FileServer(http.Dir("web/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/ui/calc", handler.CalculateHTML)
	mux.HandleFunc("/api/packs", handler.GetPacks)
	mux.HandleFunc("/api/calc", handler.CalculateJSON)
	return mux
}
