package lobby

import (
	"context"
	"fmt"
	"log"
	"time"

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
	if !e.Recovered && !e.Resubscribed {
		users, err := sub.Presence()
		if err != nil {
			h.logger.Println("error getting lobby users: ", err.Error())
		}
		if len(users) < 2 {
			sub.Publish([]byte(`{"status": "wait"}`))
		}
		if len(users) == 2 {
			ctx := context.Background()
			pipe := h.http.Pipe()
			gameStamp := time.Now()
			for _, user := range users {
				h.redis.Append(fmt.Sprintf("game:%s", gameStamp), user.User)
				pipe.AddPublish(fmt.Sprintf("lobby#%s", user.User), []byte(fmt.Sprintf(`{"status": "join", "game": "%s"}`, gameStamp)))
			}
			h.redis.Set(loadingLobbies, gameStamp, 0)
			h.http.SendPipe(ctx, pipe)
		}
	}
	h.logger.Printf("Successfully subscribed to channel %s", sub.Channel())
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
		logger.Fatalln(err)
	}

	subEventHandler := &subEventHandler{logger, r, http}
	sub.OnSubscribeSuccess(subEventHandler)
	sub.OnSubscribeError(subEventHandler)
	sub.OnUnsubscribe(subEventHandler)
	sub.OnPublish(subEventHandler)

	// Subscribe on private channel.
	sub.Subscribe()
}
