package server

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/database"
	h "github.com/ElitistNoob/chirpy/internal/handlers"
)

//go:embed index.html assets/*
var serverFS embed.FS

func initServer() error {
	const port = "8080"
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL missing from .env")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM missing from .env")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	dbQueries := database.New(db)
	appState := app.NewApp()
	appState.Queries = dbQueries
	appState.Platform = platform

	mux := http.NewServeMux()
	fileServer := http.FileServerFS(serverFS)

	handler := http.StripPrefix("/app/", fileServer)
	mux.Handle("/app/", h.MiddlewareMetricsInc(appState)(handler))

	// api
	mux.HandleFunc("GET /api/healthz", h.HealthHandler)
	mux.HandleFunc("POST /api/users", h.CreateUserHandler(appState))
	mux.HandleFunc("POST /api/chirps", h.ChirpsHandler(appState))

	// admin
	mux.HandleFunc("GET /admin/metrics", h.MetricsHandler(appState))
	mux.HandleFunc("POST /admin/reset", h.ResetHandler(appState))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files on port %s\n", s.Addr)
	return s.ListenAndServe()
}
