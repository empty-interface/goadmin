package main

import (
	"log"

	"github.com/kalitheniks/goadmin/server"
)

func main() {
	config := server.ServerConfig{
		Port: "5000",
		Host: "localhost",
	}
	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Could not setup a server: %v", err)
	}
	log.Fatalf("Server shutdown: %v", srv.ListenAndServe())
}
