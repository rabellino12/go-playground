package routes

import (
	"log"
	"net/http"

	"github.com/rabellino12/go-playground/iohttp"
	"github.com/rabellino12/go-playground/routes/auth"
	"github.com/rabellino12/go-playground/routes/game"
	"github.com/rabellino12/go-playground/routes/home"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetRoutes is the main function to setup all the server routes
func SetRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	client *mongo.Client,
	ioh *iohttp.Client,
) {
	homeHandler := home.NewHandlers(logger)
	homeHandler.SetupRoutes(mux)
	gameHandler := game.NewHandlers(logger, client)
	gameHandler.SetupRoutes(mux)
	authHandler := auth.NewHandlers(logger, ioh)
	authHandler.SetupRoutes(mux)

}
