package iohttp

import (
	"context"
	"log"

	"github.com/centrifugal/gocent"
)

// IoHTTP is a handler for IO actions
type IoHTTP struct {
	// Client is a gocent.Client instance
	Client  *gocent.Client
	context context.Context
	logger  *log.Logger
}

// Init initializes gocent http library connection
func Init(logger *log.Logger) *IoHTTP {

	c := gocent.New(gocent.Config{
		Addr: "http://centrifugo:9000",
		Key:  "some-centrifugo-api-key",
	})

	ctx := context.Background()

	return &IoHTTP{c, ctx, logger}
}

//Publish sends a message to the specified channel
func (io *IoHTTP) Publish(ch string) {
	err := io.Client.Publish(io.context, ch, []byte(`{"input": "test"}`))
	if err != nil {
		log.Fatalf("Error calling publish: %v", err)
	}
	log.Printf("Publish into channel %s successful", ch)
}
