package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupDB initializes de mongodb connection
func SetupDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))
	return conn, err
}
