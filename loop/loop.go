package loop

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/rabellino12/go-playground/iohttp"
)

const emptyLobbies string = "emptyLobbies"
const loadingLobbies string = "loadingLobbies"
const fullLobbies string = "fullLobbies"

// Handler is the loop methods handler
type Handler struct {
	io     *iohttp.Client
	quit   chan struct{}
	ticker *time.Ticker
	logger *log.Logger
	redis  *redis.Client
}

// Initialize starts the loop and creates the quit chanel
func Initialize(io *iohttp.Client, logger *log.Logger, redis *redis.Client) {
	ticker := time.NewTicker(time.Second / 20)
	quit := make(chan struct{})
	h := &Handler{io, quit, ticker, logger, redis}
	go h.loop()
	// close(quit)
}

func (h *Handler) loop() {
	for {
		select {
		case <-h.ticker.C:
			h.lobby()
		case <-h.quit:
			fmt.Println("ticker stopped")
			return
		}
	}
}

func (h *Handler) lobby() {
	users, err := h.io.Presence("$lobby:index")
	if err != nil {
		h.logger.Println("error getting lobby users: ", err.Error())
	}
	if len(users) < 2 {
		h.io.Publish("$lobby:index", []byte(`{"status": "wait"}`))
	}
	if len(users) == 2 {
		ctx := context.Background()
		pipe := h.io.Client.Pipe()
		now := time.Now()
		nanos := now.UnixNano()
		gameStamp := nanos / 1000000
		fmt.Println("New Game: ", gameStamp)
		for _, user := range users {
			if user.User != "112" {
				h.redis.Append(fmt.Sprintf("game:%s", strconv.FormatInt(gameStamp, 2)), user.User)
				pipe.AddPublish(fmt.Sprintf("lobby#%s", user.User), []byte(fmt.Sprintf(`{"status": "join", "game": "%s"}`, strconv.FormatInt(gameStamp, 2))))
				pipe.AddUnsubscribe("$lobby:index", user.User)
			}
		}
		h.redis.Set(loadingLobbies, gameStamp, 0)
		h.io.Client.SendPipe(ctx, pipe)
	}
}
