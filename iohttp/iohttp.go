package iohttp

import (
	"context"
	"log"

	"github.com/centrifugal/gocent"
	"github.com/dgrijalva/jwt-go"
)

// Client is a handler for IO actions
type Client struct {
	// Client is a gocent.Client instance
	Client  *gocent.Client
	context context.Context
	logger  *log.Logger
}

// CentrifugeAuthCustomClaims struct for jwt
type CentrifugeAuthCustomClaims struct {
	Client  string `json:"client"`
	Channel string `json:"channel"`
	jwt.StandardClaims
}

// Init initializes gocent http library connection
func Init(logger *log.Logger) *Client {

	c := gocent.New(gocent.Config{
		Addr: "http://centrifugo:9000",
		Key:  "some-centrifugo-api-key",
	})

	ctx := context.Background()

	return &Client{c, ctx, logger}
}

//Publish sends a message to the specified channel
func (io *Client) Publish(ch string, message []byte) {
	err := io.Client.Publish(io.context, ch, message)
	if err != nil {
		log.Fatalf("Error calling publish: %v", err)
	}
}

// Presence is the method to Get channel's users
func (io *Client) Presence(ch string) (map[string]gocent.ClientInfo, error) {
	users, err := io.Client.Presence(io.context, ch)
	if err != nil {
		return nil, err
	}
	return users.Presence, err
}
