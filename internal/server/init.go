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
	"github.com/ElitistNoob/chirpy/internal/handlers/chirps"
	"github.com/ElitistNoob/chirpy/internal/handlers/users"
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

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET missing from .env")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	dbQueries := database.New(db)
	appState := &app.App{
		Queries:  dbQueries,
		Platform: platform,
		Secret:   secret,
	}

	mux := http.NewServeMux()
	fileServer := http.FileServerFS(serverFS)

	handler := http.StripPrefix("/app/", fileServer)
	mux.Handle("/app/", h.MiddlewareMetricsInc(appState)(handler))

	// healtz endpoints
	mux.HandleFunc("GET /api/healthz", h.HealthHandler)

	// users endpoints
	mux.HandleFunc("POST /api/users", users.CreateHandler(appState))
	mux.HandleFunc("POST /api/login", users.LoginHandler(appState))
	mux.HandleFunc("POST /api/refresh", users.RefreshHandler(appState))
	mux.HandleFunc("POST /api/revoke", users.RevokeHandler(appState))

	// chirps endpoints
	mux.HandleFunc("POST /api/chirps", chirps.CreateHandler(appState))
	mux.HandleFunc("GET /api/chirps", chirps.GetAllHandler(appState))
	mux.HandleFunc("GET /api/chirps/{chirpID}", chirps.GetHandler(appState))

	// admin endpoints
	mux.HandleFunc("GET /admin/metrics", h.MetricsHandler(appState))
	mux.HandleFunc("POST /admin/reset", h.ResetHandler(appState))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files on port %s\n", s.Addr)
	return s.ListenAndServe()
}
