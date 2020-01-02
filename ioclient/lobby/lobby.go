package lobby

import (
	"log"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/gocent"
	"github.com/go-redis/redis/v7"
)

type subEventHandler struct {
	logger *log.Logger
	redis  *redis.Client
	http   *gocent.Client
}

const emptyLobbies string = "emptyLobbies"
const loadingLobbies string = "loadingLobbies"
const fullLobbies string = "fullLobbies"

func (h *subEventHandler) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	h.logger.Printf("Successfully subscribed to channel %s", sub.Channel())
}

func (h *subEventHandler) OnJoin(sub *centrifuge.Subscription, e centrifuge.JoinEvent) {
	h.logger.Println("New lobby join")
}

func (h *subEventHandler) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	h.logger.Printf("Error subscribing to channel %s: %v", sub.Channel(), e.Error)
}

func (h *subEventHandler) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	h.logger.Printf("Unsubscribed from channel %s", sub.Channel())
}

func (h *subEventHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	h.logger.Printf("New message received from channel %s: %s", sub.Channel(), string(e.Data))
}

// Initialize lobby io controller
func Initialize(c *centrifuge.Client, r *redis.Client, logger *log.Logger, http *gocent.Client) {
	sub, err := c.NewSubscription("$lobby:index")
	if err != nil {
		logger.Println(err)
	}

	subEventHandler := &subEventHandler{logger, r, http}
	sub.OnSubscribeSuccess(subEventHandler)
	sub.OnSubscribeError(subEventHandler)
	sub.OnUnsubscribe(subEventHandler)
	sub.OnPublish(subEventHandler)
	sub.OnJoin(subEventHandler)
	// Subscribe on private channel.
	sub.Subscribe()
}
