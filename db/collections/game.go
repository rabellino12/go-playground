package game

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Game is the game collection struct
type Game struct {
	Players []string `json:"players"`
	Name    string   `json:"name"`
}

// Handler is the Game collection handler
type Handler struct {
	client *mongo.Client
}

func (h *Handler) getCollection() *mongo.Collection {
	collection := h.client.Database("main").Collection("games")
	return collection
}

// Get returns a game from the database
func (h *Handler) Get(gameID string) (Game, error) {
	players := []string{"player1", "player2"}
	return Game{Players: players}, nil
}

// Insert creates a new game instance in the database
func (h *Handler) Insert(game Game) (*mongo.InsertOneResult, error) {
	collection := h.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	doc, _ := bson.Marshal(game)
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// NewHandler creates a new game collection handler instance
func NewHandler(client *mongo.Client) *Handler {
	handler := &Handler{client}
	return handler
}
