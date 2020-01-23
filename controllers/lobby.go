package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/centrifugal/gocent"
	redis "github.com/go-redis/redis/v7"
	game "github.com/rabellino12/go-playground/db/collections"
	"github.com/rabellino12/go-playground/helper"
	"github.com/rabellino12/go-playground/ioclient/match"
	"github.com/rabellino12/go-playground/iohttp"
)

// Lobby is the struct for handling lobby actions on loop
type Lobby struct {
	IO          *iohttp.Client
	Logger      *log.Logger
	Redis       *redis.Client
	GameHandler *game.Handler
}

// RunLoop method acts as init for lobby loop handler
func (l *Lobby) RunLoop() {
	users, err := l.IO.Presence("$lobby:index")
	var players []gocent.ClientInfo
	for _, user := range users {
		if user.User != "112" {
			players = append(players, user)
		}
	}
	if err != nil {
		l.Logger.Println("error getting lobby users: ", err.Error())
	}
	if len(players) < 2 {
		l.IO.Publish("$lobby:index", []byte(`{"event": "wait"}`))
	}
	if len(players) == 2 {
		ctx := context.Background()
		pipe := l.IO.Client.Pipe()
		var newGame game.Body
		playersList := []game.Player{}
		for i, player := range players {
			initialPosition := helper.GetPlayerInitialPosition(i)
			playersList = append(playersList, game.Player{Index: i, Position: initialPosition, ID: player.User})
		}
		newGame = game.Body{Players: playersList}
		resGame, err := l.GameHandler.Insert(&newGame)
		if err != nil {
			l.Logger.Println("error creating new game: ", err.Error())
			return
		}
		l.Redis.LPush("games", resGame.ID.Hex())
		for _, p := range playersList {
			l.Logger.Println("game player: ", p)
			joinJS, joinErr := json.Marshal(match.JoinEvent{
				Event:   "join",
				Game:    resGame.ID.Hex(),
				Players: playersList,
			})
			if joinErr != nil {
				l.Logger.Println("error parsing join event: ", joinErr.Error())
				return
			}
			err = pipe.AddPublish(fmt.Sprintf("lobby#%s", p.ID), joinJS)
			if err != nil {
				l.Logger.Println("error adding personal lobby notification: ", err.Error())
				return
			}
			err = pipe.AddUnsubscribe("$lobby:index", p.ID)
			if err != nil {
				l.Logger.Println("error unsubscribing user from lobby:index: ", err.Error())
				return
			}
		}
		lobbyReply, err2 := l.IO.Client.SendPipe(ctx, pipe)
		if err2 != nil {
			l.Logger.Println("error sending pipe on lobby loop ", err2.Error())
			return
		}
		for index, rep := range lobbyReply {
			if rep.Error != nil {
				l.Logger.Println("error sending pipe "+strconv.FormatInt(int64(index), 10)+" on lobby loop ", rep.Error)
			} else {
				l.Logger.Println("sent pipe "+strconv.FormatInt(int64(index), 10)+" on lobby loop ", rep.Result)
			}
		}
	}
}
