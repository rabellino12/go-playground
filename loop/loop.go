package loop

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/centrifugal/gocent"
	redis "github.com/go-redis/redis/v7"
	game "github.com/rabellino12/go-playground/db/collections"
	"github.com/rabellino12/go-playground/helper"
	"github.com/rabellino12/go-playground/ioclient/match"
	"github.com/rabellino12/go-playground/iohttp"
	"go.mongodb.org/mongo-driver/mongo"
)

const emptyLobbies string = "emptyLobbies"
const loadingLobbies string = "loadingLobbies"
const fullLobbies string = "fullLobbies"

// Handler is the loop methods handler
type Handler struct {
	io          *iohttp.Client
	quit        chan struct{}
	ticker      *time.Ticker
	logger      *log.Logger
	redis       *redis.Client
	mongoClient *mongo.Client
	gameHandler *game.Handler
}

// Initialize starts the loop and creates the quit chanel
func Initialize(io *iohttp.Client, logger *log.Logger, redis *redis.Client, mongoClient *mongo.Client) {
	ticker := time.NewTicker(time.Second / 20)
	quit := make(chan struct{})
	gameHandler := game.NewHandler(mongoClient)
	h := &Handler{io, quit, ticker, logger, redis, mongoClient, gameHandler}
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
		var newGame game.Body
		playersList := []game.Player{}
		for i, player := range players {
			initialPosition := helper.GetPlayerInitialPosition(i)
			playersList = append(playersList, game.Player{Index: i, Position: initialPosition, ID: player.User})
		}
		newGame = game.Body{Players: playersList}
		resGame, err := h.gameHandler.Insert(&newGame)
		if err != nil {
			h.logger.Println("error creating new game: ", err.Error())
			return
		}
		h.redis.LPush("games", resGame.ID.Hex())
		for _, p := range playersList {
			h.logger.Println("game player: ", p)
			joinJS, joinErr := json.Marshal(match.JoinEvent{
				Event:   "join",
				Game:    resGame.ID.Hex(),
				Players: playersList,
			})
			if joinErr != nil {
				h.logger.Println("error parsing join event: ", joinErr.Error())
				return
			}
			err = pipe.AddPublish(fmt.Sprintf("lobby#%s", p.ID), joinJS)
			if err != nil {
				h.logger.Println("error adding personal lobby notification: ", err.Error())
				return
			}
			err = pipe.AddUnsubscribe("$lobby:index", p.ID)
			if err != nil {
				h.logger.Println("error unsubscribing user from lobby:index: ", err.Error())
				return
			}
		}
		lobbyReply, err2 := h.io.Client.SendPipe(ctx, pipe)
		if err2 != nil {
			h.logger.Println("error sending pipe on lobby loop ", err2.Error())
			return
		}
		for index, rep := range lobbyReply {
			if rep.Error != nil {
				h.logger.Println("error sending pipe "+strconv.FormatInt(int64(index), 10)+" on lobby loop ", rep.Error)
			} else {
				h.logger.Println("sent pipe "+strconv.FormatInt(int64(index), 10)+" on lobby loop ", rep.Result)
			}
		}
	}
}
