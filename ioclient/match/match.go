package match

import (
	"encoding/json"
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

// Move interface
type Move struct {
	Action    string `json:"action"`
	Timestamp int    `json:"timestamp"`
	MatchID   string `json:"matchId"`
	UserID    string `json:"userId"`
}

func (h *subEventHandler) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	h.logger.Printf("Unsubscribed from channel %s", sub.Channel())
}

func (h *subEventHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	h.logger.Printf("New message received from channel %s: %s", sub.Channel(), string(e.Data))
	var move *Move
	err := json.Unmarshal(e.Data, &move)
	if err != nil {
		return
	}
	h.redis.LPush("match:"+move.MatchID, move)
}

// Initialize lobby io controller
func Initialize(c *centrifuge.Client, r *redis.Client, logger *log.Logger, http *gocent.Client) {
	sub, err := c.NewSubscription("match")
	if err != nil {
		logger.Println(err)
	}

	subEventHandler := &subEventHandler{logger, r, http}
	sub.OnUnsubscribe(subEventHandler)
	sub.OnPublish(subEventHandler)
	// Subscribe on private channel.
	sub.Subscribe()
}
