package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/gptlv/re-partners/packs/internal/api"
	"github.com/gptlv/re-partners/packs/internal/app"
	"github.com/gptlv/re-partners/packs/internal/repository"
	"github.com/gptlv/re-partners/packs/internal/router"
	"github.com/gptlv/re-partners/packs/pkg/db"
)

func main() {
	databaseURL, err := buildDatabaseURL()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Open(databaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
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

func buildDatabaseURL() (string, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return "", fmt.Errorf("POSTGRES_HOST is required")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		return "", fmt.Errorf("POSTGRES_USER is required")
	}

	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		return "", fmt.Errorf("POSTGRES_DB is required")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	u := url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + dbName,
	}
	if password != "" {
		u.User = url.UserPassword(user, password)
	} else {
		u.User = url.User(user)
	}

	query := u.Query()
	query.Set("sslmode", sslmode)
	u.RawQuery = query.Encode()

	return u.String(), nil
}
