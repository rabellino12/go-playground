package mongodb

import (
	"context"
	"time"

	"github.com/rabellino12/go-playground/helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbURL = helper.GoDotEnvVariable("DB_URL")

// SetupDB initializes de mongodb connection
func SetupDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	return conn, err
}
