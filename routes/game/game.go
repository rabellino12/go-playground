package game

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	game "github.com/rabellino12/go-playground/db/collections"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handlers is a Struct that contains handler methods and shared server data
type Handlers struct {
	logger      *log.Logger
	gameHandler *game.Handler
}

func (h *Handlers) gateway(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
	case "POST":
		h.post(w, r)
	}
}

// get is a handler funciton for get single game route "/"
func (h *Handlers) get(w http.ResponseWriter, r *http.Request) {
	game := game.Game{Name: "gametest", Players: []string{"player1", "player2"}}

	js, err := json.Marshal(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}

//post is a handler function to insert new game
func (h *Handlers) post(w http.ResponseWriter, r *http.Request) {
	var game game.Game

	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.gameHandler.Insert(game)
	js, err := json.Marshal(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}

// Logger is the game logging middleware
func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.logger.Println("Processing game request")
		defer h.logger.Printf("Request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

// SetupRoutes creates all game related routes
func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	h.logger.Println("game route setup")
	mux.HandleFunc("/game", h.Logger(h.gateway))
}

// NewHandlers returns a game page handlers struct
func NewHandlers(logger *log.Logger, client *mongo.Client) *Handlers {
	return &Handlers{
		logger:      logger,
		gameHandler: game.NewHandler(client),
	}
}
