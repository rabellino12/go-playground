package game

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Body is the game request body struct
type Body struct {
	Players []Player `bson:"players" json:"players"`
}

// Game is the game collection struct
type Game struct {
	Body `bson:",inline"`
	// ObjectId() or objectid.ObjectID is deprecated--use primitive instead
	ID primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
}

// Position is a x,y axis point
type Position struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

// Player is the game player structure, with player id, index and position "x,y"
type Player struct {
	Index    int      `bson:"index, omitempty" json:"index"`
	Position Position `bson:"position" json:"position"`
	ID       string   `bson:"id, omitempty" json:"id"`
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
	collection := h.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		return Game{}, err
	}
	newGame := collection.FindOne(ctx, bson.M{"_id": id})
	var game Game
	err = newGame.Decode(&game)
	if err != nil {
		return game, err
	}
	return game, nil
}

// Insert creates a new game instance in the database
func (h *Handler) Insert(game *Body) (Game, error) {
	collection := h.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	doc, _ := bson.Marshal(game)
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return Game{}, err
	}
	newGame := collection.FindOne(ctx, bson.M{"_id": res.InsertedID.(primitive.ObjectID)})
	var resGame Game
	err = newGame.Decode(&resGame)
	log.Printf("Found a single document: %+v\n", resGame)
	if err != nil {
		return resGame, err
	}
	return resGame, nil
}

// NewHandler creates a new game collection handler instance
func NewHandler(client *mongo.Client) *Handler {
	handler := &Handler{client}
	return handler
}
