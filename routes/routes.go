package routes

import (
	"log"
	"net/http"

	"github.com/rabellino12/go-playground/ioclient/iohttp"
	"github.com/rabellino12/go-playground/routes/game"
	"github.com/rabellino12/go-playground/routes/home"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetRoutes is the main function to setup all the server routes
func SetRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	client *mongo.Client,
	iohttp *iohttp.IoHTTP,
) {
	homeHandler := home.NewHandlers(logger)
	homeHandler.SetupRoutes(mux)
	gameHandler := game.NewHandlers(logger, client)
	gameHandler.SetupRoutes(mux)
}
