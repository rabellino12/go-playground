package main

import (
	"log"
	"net/http"
	"os"

	redis "github.com/rabellino12/go-playground/cache"
	controller "github.com/rabellino12/go-playground/controllers"
	mongodb "github.com/rabellino12/go-playground/db"
	game "github.com/rabellino12/go-playground/db/collections"
	"github.com/rabellino12/go-playground/helper"
	"github.com/rabellino12/go-playground/ioclient"
	"github.com/rabellino12/go-playground/iohttp"
	"github.com/rabellino12/go-playground/loop"
	"github.com/rabellino12/go-playground/routes"
	"github.com/rabellino12/go-playground/server"
)

func main() {
	var serverAddr = helper.GoDotEnvVariable("SERVER_ADDRESS")
	logger := log.New(os.Stdout, "gophercon-tutorial", log.LstdFlags|log.Lshortfile)
	logger.Println("server address", serverAddr)
	mux := http.NewServeMux()
	initialize(mux, logger)
	srv := server.NewServer(mux, serverAddr)
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
	gameHandler := game.NewHandler(mongoClient)
	lobbyController := &controller.Lobby{
		IO:          ioh,
		Logger:      logger,
		Redis:       r,
		GameHandler: gameHandler,
	}
	loop.Initialize(lobbyController, 20)
	routes.SetRoutes(mux, logger, mongoClient, ioh)
}
