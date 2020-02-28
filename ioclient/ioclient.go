// Package ioclient contains the centrifugo connection logic
package ioclient

import (
	"log"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/gocent"
	"github.com/go-redis/redis/v7"
	"github.com/rabellino12/go-playground/helper"
)

var centrifugoWS = helper.GoDotEnvVariable("CENTRIFUGO_WS")

func connToken(user string, exp int64) string {
	t, err := helper.GetJWT(user)
	if err != nil {
		panic(err)
	}
	return t
}

type eventHandler struct{}

func (h *eventHandler) OnPrivateSub(c *centrifuge.Client, e centrifuge.PrivateSubEvent) (string, error) {
	token, err := helper.GetSubscriptionJWT(e.ClientID, e.Channel)
	if err != nil {
		log.Println("error on private sub")
		return "", err
	}
	return token, nil
}

func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	log.Println("Connected")
}

func (h *eventHandler) OnError(c *centrifuge.Client, e centrifuge.ErrorEvent) {
	log.Println("Error", e.Message)
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	log.Println("Disconnected", e.Reason)
}

// NewConnection creates a new centrifuge pub/sub connection
func NewConnection() *centrifuge.Client {
	wsURL := centrifugoWS + "/connection/websocket"

	c := centrifuge.New(wsURL, centrifuge.DefaultConfig())
	c.SetToken(connToken("112", 0))
	handler := &eventHandler{}
	c.OnPrivateSub(handler)
	c.OnDisconnect(handler)
	c.OnConnect(handler)
	c.OnError(handler)

	err := c.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

// Connect starts the centrifugo connection
func Connect(
	g *gocent.Client,
	r *redis.Client,
	l *log.Logger,
) {
	l.Println("Start program")
	c := NewConnection()
	defer c.Close()
	// lobby.Initialize(c, r, l, g)
	l.Println("IO initialized")
	select {}
}
