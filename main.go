package main

import (
	"log"
	"net/http"
	"os"

	redis "github.com/rabellino12/go-playground/cache"
	mongodb "github.com/rabellino12/go-playground/db"
	"github.com/rabellino12/go-playground/ioclient"
	"github.com/rabellino12/go-playground/iohttp"
	"github.com/rabellino12/go-playground/loop"
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
	initialize(mux, logger)
	srv := server.NewServer(mux, getServerAddress())
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

func initialize(mux *http.ServeMux, logger *log.Logger) {
	mongoClient, _ := mongodb.SetupDB()
	ioh := iohttp.Init(logger)
	r := redis.NewClient(logger)
	go ioclient.Connect(ioh.Client, r, logger)
	loop.Initialize(ioh, logger, r, mongoClient)
	routes.SetRoutes(mux, logger, mongoClient, ioh)
}

func getServerAddress() string {
	if serverAddr != "" {
		return "0.0.0.0:" + serverAddr
	}
	return "0.0.0.0:8080"
}
