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
	DoneMoves  []matchIO.Move
}

// MakeMatch starts a new match instance with its own loop, intended to be used on its own goroutine
func MakeMatch(io *iohttp.Client, logger *log.Logger, redis *redis.Client, gameObj game.Game) {
	worldScene := scenes.MakeMatch(gameObj)
	matchHandler := &Match{
		IO:         io,
		Logger:     logger,
		Redis:      redis,
		WorldScene: worldScene,
		ID:         gameObj.ID.Hex(),
	}
	logger.Println("started match: " + matchHandler.ID)
	loop.Initialize(matchHandler, 60)
	logger.Println("after loop, match: " + matchHandler.ID)
	c := ioclient.NewConnection()
	defer c.Close()
	logger.Println("after connection, match: " + matchHandler.ID)
	sub, er := c.NewSubscription("$match:" + gameObj.ID.Hex())
	logger.Println("after subscription to channel: ", sub.Channel())
	if er != nil {
		logger.Println(er.Error())
	}
	sub.OnPublish(matchHandler)
	sub.OnSubscribeSuccess(matchHandler)
	sub.OnSubscribeError(matchHandler)
	sub.OnUnsubscribe(matchHandler)
	sub.OnJoin(matchHandler)
	err := sub.Subscribe()
	if err != nil {
		logger.Println(err.Error())
	}
	logger.Println("after subscribe, match: " + matchHandler.ID)
	select {}
}

// GetID returns the current match id
func (m *Match) GetID() string {
	return m.ID
}

// RunLoop method acts as init for match loop handler
func (m *Match) RunLoop() {
	channel := "$snapshot:" + m.ID
	if len(m.Moves) == 0 {
		m.Logger.Println("no moves found")
		return
	}
	for i, move := range m.Moves {
		m.Logger.Println("im adding a move")
		m.WorldScene.AddMove(move)
		m.DoneMoves = append(m.DoneMoves, move)
		m.removeMove(i)
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
	m.Logger.Println("Received message on match: " + m.ID)
	var move matchIO.Move
	err := json.Unmarshal(e.Data, &move)
	if err != nil {
		m.Logger.Println("error marshaling match message: ", err.Error())
		return
	}
	m.Moves = append(m.Moves, move)
}

func (m *Match) removeMove(i int) []matchIO.Move {
	m.Moves[len(m.Moves)-1], m.Moves[i] = m.Moves[i], m.Moves[len(m.Moves)-1]
	return m.Moves[:len(m.Moves)-1]
}

// OnSubscribeSuccess method handles the subscribe event for the match channel
func (m *Match) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	m.Logger.Printf("Successfully subscribed to channel %s", sub.Channel())
}

// OnJoin method handles the join event for the match channel
func (m *Match) OnJoin(sub *centrifuge.Subscription, e centrifuge.JoinEvent) {
	m.Logger.Println("New match join")
}

// OnSubscribeError method handles the subscribe error event for the match channel
func (m *Match) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	m.Logger.Printf("Error subscribing to channel %s: %v", sub.Channel(), e.Error)
}

// OnUnsubscribe method handles the unsubscribe event for the match channel
func (m *Match) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	m.Logger.Printf("Unsubscribed from channel %s", sub.Channel())
}
