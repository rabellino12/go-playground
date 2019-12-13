package iohttp

import (
	"context"
	"log"

	"github.com/centrifugal/gocent"
	"github.com/dgrijalva/jwt-go"
	"github.com/rabellino12/go-playground/ioclient"
)

// Client is a handler for IO actions
type Client struct {
	// Client is a gocent.Client instance
	Client  *gocent.Client
	context context.Context
	logger  *log.Logger
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
func (io *Client) Publish(ch string) {
	err := io.Client.Publish(io.context, ch, []byte(`{"input": "test"}`))
	if err != nil {
		log.Fatalf("Error calling publish: %v", err)
	}
	log.Printf("Publish into channel %s successful", ch)
}

// GetJWT receives a username and returns a jwt
func (io *Client) GetJWT(userName string) (string, error) {
	signingKey := []byte(ioclient.CentrifugoSecret)

	// Create the Claims
	claims := &jwt.StandardClaims{
		// ExpiresAt: 150000,
		Issuer: userName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
