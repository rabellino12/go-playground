package controller

import (
	"context"
	"log"

	redis "github.com/go-redis/redis/v7"
	game "github.com/rabellino12/go-playground/db/collections"
	"github.com/rabellino12/go-playground/iohttp"
	"github.com/rabellino12/go-playground/loop"
	scenes "github.com/rabellino12/go-playground/scenes/match"
)

// Match is the struct for handling lobby actions on loop
type Match struct {
	IO     *iohttp.Client
	Logger *log.Logger
	Redis  *redis.Client
	ID     string
}

// MakeMatch starts a new match instance with its own loop, intended to be used on its own goroutine
func MakeMatch(io *iohttp.Client, logger *log.Logger, redis *redis.Client, gameObj game.Game) {
	matchHandler := &Match{
		IO:     io,
		Logger: logger,
		Redis:  redis,
	}
	scenes.MakeMatch(gameObj)
	loop.Initialize(matchHandler, 60)
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
	pipe := m.IO.Client.Pipe()
	m.Redis.Del(channel)
	for _, message := range messages {
		pipe.AddPublish("$"+channel, []byte(message))
	}
	m.IO.Client.SendPipe(context.Background(), pipe)

}
