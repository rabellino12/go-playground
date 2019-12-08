package main

import (
	"log"
	"net/http"
	"os"

	mongodb "github.com/rabellino12/go-playground/db"
	"github.com/rabellino12/go-playground/routes"
	"github.com/rabellino12/go-playground/server"
)

var (
	serverAddr = os.Getenv("SERVER_ADDRESS")
)

func main() {
	logger := log.New(os.Stdout, "gophercon-tutorial", log.LstdFlags|log.Lshortfile)
	logger.Println("server address", serverAddr)
	mux := http.NewServeMux()
	srv := server.NewServer(mux, getServerAddress())
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

func initialize(mux *http.ServeMux, logger *log.Logger) {
	client, _ := mongodb.SetupDB()
	routes.SetRoutes(mux, logger, client)

}

func getServerAddress() string {
	if serverAddr != "" {
		return "0.0.0.0:" + serverAddr
	}
	return "0.0.0.0:8080"
}
