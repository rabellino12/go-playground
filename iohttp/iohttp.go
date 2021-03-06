package iohttp

import (
	"context"
	"log"

	"github.com/centrifugal/gocent"
	"github.com/dgrijalva/jwt-go"
	"github.com/rabellino12/go-playground/helper"
)

var centrifugoURL = helper.GoDotEnvVariable("CENTRIFUGO_PRIVATE_URL")

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
		Addr: centrifugoURL,
		Key:  "some-centrifugo-api-key",
	})

	ctx := context.Background()

	return &Client{c, ctx, logger}
}

//Publish sends a message to the specified channel
func (io *Client) Publish(ch string, message []byte) error {
	err := io.Client.Publish(io.context, ch, message)
	return err
}

// Presence is the method to Get channel's users
func (io *Client) Presence(ch string) (map[string]gocent.ClientInfo, error) {
	users, err := io.Client.Presence(io.context, ch)
	if err != nil {
		return nil, err
	}
	return users.Presence, err
}

// History is the method to Get channel's message history
func (io *Client) History(ch string) (gocent.HistoryResult, error) {
	history, err := io.Client.History(io.context, ch)
	if err != nil {
		return gocent.HistoryResult{}, err
	}
	return history, err
}
