package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rabellino12/go-playground/routes/home"
	"github.com/rabellino12/go-playground/server"
)

var (
	serverAddr = os.Getenv("SERVER_ADDRESS")
)

func main() {
	logger := log.New(os.Stdout, "gophercon-tutorial", log.LstdFlags|log.Lshortfile)
	homeHandler := home.NewHandlers(logger)
	mux := http.NewServeMux()
	homeHandler.SetupRoutes(mux)
	srv := server.NewServer(mux, serverAddr)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
