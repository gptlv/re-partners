package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gptlv/re-partners/packs/internal/api"
	"github.com/gptlv/re-partners/packs/internal/app"
	"github.com/gptlv/re-partners/packs/internal/db"
	"github.com/gptlv/re-partners/packs/internal/repository"
	"github.com/gptlv/re-partners/packs/internal/router"
	"github.com/gptlv/re-partners/packs/internal/view"
)

func main() {
	database, err := db.Open("file:packs.db?_journal=WAL")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))

	repo := repository.NewPackRepository(database)
	service := app.NewService(repo)
	renderer := view.NewTemplateRenderer(tmpl)
	handler := api.NewHandler(service, renderer)
	mux := router.New(handler)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
