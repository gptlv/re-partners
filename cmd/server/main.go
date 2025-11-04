package main

import (
	"log"
	"net/http"

	"github.com/gptlv/re-partners/packs/internal/api"
	"github.com/gptlv/re-partners/packs/internal/app"
	"github.com/gptlv/re-partners/packs/internal/db"
	"github.com/gptlv/re-partners/packs/internal/repository"
	"github.com/gptlv/re-partners/packs/internal/router"
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

	repo := repository.NewPackRepository(database)
	service := app.NewService(repo)
	handler := api.NewHandler(service)
	mux := router.New(handler)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
