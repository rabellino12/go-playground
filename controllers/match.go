package controller

import (
	"encoding/json"
	"log"

	"github.com/centrifugal/centrifuge-go"
	redis "github.com/go-redis/redis/v7"
	game "github.com/rabellino12/go-playground/db/collections"
	"github.com/rabellino12/go-playground/ioclient"
	matchIO "github.com/rabellino12/go-playground/ioclient/match"
	"github.com/rabellino12/go-playground/iohttp"
	"github.com/rabellino12/go-playground/loop"
	scenes "github.com/rabellino12/go-playground/scenes/match"
)

// Match is the struct for handling lobby actions on loop
type Match struct {
	IO         *iohttp.Client
	Logger     *log.Logger
	Redis      *redis.Client
	ID         string
	WorldScene *scenes.WorldScene
	Moves      []matchIO.Move
}

// MakeMatch starts a new match instance with its own loop, intended to be used on its own goroutine
func MakeMatch(io *iohttp.Client, logger *log.Logger, redis *redis.Client, gameObj game.Game) {
	worldScene := scenes.MakeMatch(gameObj)
	matchHandler := &Match{
		IO:         io,
		Logger:     logger,
		Redis:      redis,
		WorldScene: worldScene,
	}
	c := ioclient.NewConnection()
	defer c.Close()
	go matchIO.ListenMatch(c, matchHandler)
	go loop.Initialize(matchHandler, 60)
	select {}
}

// GetID returns the current match id
func (m *Match) GetID() string {
	return m.ID
}

// RunLoop method acts as init for match loop handler
func (m *Match) RunLoop() {
	channel := "match:" + m.ID
	messagesQuery := m.Redis.LRange(channel, 0, -1)
	messages, err := messagesQuery.Result()
	if err != nil {
		m.Logger.Println("error getting game history: ", err.Error())
		return
	}
	m.Redis.Del(channel)
	var move matchIO.Move
	for _, message := range messages {
		err := json.Unmarshal([]byte(message), &move)
		if err != nil {
			m.Logger.Println("error getting game history: ", err.Error())
		} else {
			m.WorldScene.AddMove(move)
		}
	}

}

// OnPublish handles centrifuge subscription publish
func (m *Match) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {

}
