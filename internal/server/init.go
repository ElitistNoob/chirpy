package server

import (
	"log"
	"net/http"
)

func initServer() error {
	const (
		filepathRoot = "."
		port         = "8080"
	)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, s.Addr)
	return s.ListenAndServe()
}
