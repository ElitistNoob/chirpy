package server

import (
	"embed"
	"log"
	"net/http"

	"github.com/ElitistNoob/chirpy/internal/healthz"
	m "github.com/ElitistNoob/chirpy/internal/metrics"
)

//go:embed index.html assets/*
var serverFS embed.FS

func initServer() error {
	const port = "8080"

	mux := http.NewServeMux()
	fileServer := http.FileServerFS(serverFS)
	metrics := m.New()

	mux.Handle("/app/", metrics.MiddlewareMetricsInt(http.StripPrefix("/app/", fileServer)))
	mux.HandleFunc("GET /api/healthz", healthz.Handler)
	mux.HandleFunc("GET /admin/metrics", metrics.Handler)
	mux.HandleFunc("POST /admin/reset", metrics.ResetHandler)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files on port %s\n", s.Addr)
	return s.ListenAndServe()
}
