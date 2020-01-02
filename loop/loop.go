package loop

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabellino12/go-playground/iohttp"
)

// Handler is the loop methods handler
type Handler struct {
	io     *iohttp.Client
	quit   chan struct{}
	ticker *time.Ticker
	logger *log.Logger
}

// Initialize starts the loop and creates the quit chanel
func Initialize(io *iohttp.Client, logger *log.Logger) {
	ticker := time.NewTicker(time.Second / 20)
	quit := make(chan struct{})
	h := &Handler{io, quit, ticker, logger}
	go h.loop()
	// close(quit)
}

func (h *Handler) loop() {
	for {
		select {
		case t := <-h.ticker.C:
			fmt.Println("Tick at", t)
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
		gameStamp := time.Now().Nanosecond() / 1000
		for _, user := range users {
			h.redis.Append(fmt.Sprintf("game:%s", gameStamp), user.User)
			pipe.AddPublish(fmt.Sprintf("lobby#%s", user.User), []byte(fmt.Sprintf(`{"status": "join", "game": "%s"}`, gameStamp)))
		}
		h.redis.Set(loadingLobbies, gameStamp, 0)
		h.io.Client.SendPipe(ctx, pipe)
	}
}
