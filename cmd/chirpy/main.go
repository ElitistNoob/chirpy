package main

import (
	"github.com/ElitistNoob/chirpy/internal/server"
	_ "github.com/lib/pq"
)

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed loading .env: %s", err)
	}
	server.RunServer()
}
