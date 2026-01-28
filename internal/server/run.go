package server

import "log"

func RunServer() {
	if err := initServer(); err != nil {
		log.Fatal(err)
	}
}
