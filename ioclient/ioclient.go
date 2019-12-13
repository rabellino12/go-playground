// Package ioclient contains the centrifugo connection logic
package ioclient

import (
	"log"
	"time"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/gocent"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/rabellino12/go-playground/ioclient/lobby"
)

// CentrifugoSecret is the centrifugo server instance secret key
const CentrifugoSecret = "some-centrifugo-secret-key"

func connToken(user string, exp int64) string {
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(CentrifugoSecret))
	if err != nil {
		panic(err)
	}
	return t
}

func subscribeToken(channel string, client string, exp int64) string {
	claims := jwt.MapClaims{"channel": channel, "client": client}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
	if err != nil {
		panic(err)
	}
	return t
}

type eventHandler struct{}

func (h *eventHandler) OnPrivateSub(c *centrifuge.Client, e centrifuge.PrivateSubEvent) (string, error) {
	token := subscribeToken(e.Channel, e.ClientID, time.Now().Unix()+10)
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

func newConnection() *centrifuge.Client {
	wsURL := "ws://centrifugo:8081/connection/websocket"

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
	log.Println("Start program")
	c := newConnection()
	lobby.Initialize(c, r, l, g)
	defer c.Close()
}
