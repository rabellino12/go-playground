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
	channel := "$match:" + m.ID
	for _, move := range m.Moves {
		m.WorldScene.AddMove(move)
	}
	snapshot := m.WorldScene.GetSnapshot()
	js, err := json.Marshal(snapshot)
	if err != nil {
		m.Logger.Println("error marshaling snapshot to json: ", err.Error())
		return
	}
	err = m.IO.Publish(channel, js)
	if err != nil {
		m.Logger.Println("error publishing snapshot: ", err.Error())
		return
	}

}

// OnPublish handles centrifuge subscription publish
func (m *Match) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	var move matchIO.Move
	err := json.Unmarshal(e.Data, &move)
	if err != nil {
		m.Logger.Println("error marshaling match message: ", err.Error())
		return
	}
	m.Moves = append(m.Moves, move)
}
