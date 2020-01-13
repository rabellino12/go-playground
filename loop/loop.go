package loop

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/centrifugal/gocent"
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
			h.matches()
		case <-h.quit:
			fmt.Println("ticker stopped")
			return
		}
	}
}

func (h *Handler) handleMatch(game string) {
	channel := "match:" + game
	messagesQuery := h.redis.LRange(channel, 0, -1)
	messages, err := messagesQuery.Result()
	if err != nil {
		h.logger.Println("error getting game history: ", err.Error())
		return
	}
	pipe := h.io.Client.Pipe()
	h.redis.Del(channel)
	for _, message := range messages {
		pipe.AddPublish("$"+channel, []byte(message))
	}
	h.io.Client.SendPipe(context.Background(), pipe)
}

func (h *Handler) matches() {
	gamesQuery := h.redis.LRange("games", 0, -1)
	games, err := gamesQuery.Result()
	if err != nil {
		h.logger.Println("error getting game: ", err.Error())
		return
	}
	for _, game := range games {
		h.handleMatch(game)
	}
}

func (h *Handler) lobby() {
	users, err := h.io.Presence("$lobby:index")
	var players []gocent.ClientInfo
	for _, user := range users {
		if user.User != "112" {
			players = append(players, user)
		}
	}
	if err != nil {
		h.logger.Println("error getting lobby users: ", err.Error())
	}
	if len(players) < 2 {
		h.io.Publish("$lobby:index", []byte(`{"event": "wait"}`))
	}
	if len(players) == 2 {
		ctx := context.Background()
		pipe := h.io.Client.Pipe()
		now := time.Now()
		nanos := now.UnixNano()
		gameStamp := nanos / 1000000
		fmt.Println("New Game: ", gameStamp)
		h.redis.LPush("games", strconv.FormatInt(gameStamp, 10))
		for _, player := range players {
			pipe.AddPublish(fmt.Sprintf("lobby#%s", player.User), []byte(fmt.Sprintf(`{"event": "join", "game": "%s"}`, strconv.FormatInt(gameStamp, 10))))
			pipe.AddUnsubscribe("$lobby:index", player.User)
		}
		h.redis.Set(loadingLobbies, gameStamp, 0)
		h.io.Client.SendPipe(ctx, pipe)
	}
}
