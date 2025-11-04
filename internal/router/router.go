package router

import (
	"net/http"

	"github.com/gptlv/re-partners/packs/internal/api"
)

// New constructs the HTTP mux wiring handlers to routes.
func New(handler *api.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/packs", handler.GetPacks)
	mux.HandleFunc("/api/calc", handler.CalculateJSON)
	mux.HandleFunc("/api/sizes", handler.CreateSize)
	mux.HandleFunc("/api/sizes/", handler.DeleteSize)
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	return mux
}
