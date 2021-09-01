package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

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
	srv.ListenAndServe()

	gracefullShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefullShutdown, os.Interrupt)
	<-gracefullShutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	defer srv.Close(ctx)
}
