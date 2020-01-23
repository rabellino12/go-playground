package controller

import (
	"context"
	"log"

	redis "github.com/go-redis/redis/v7"
	"github.com/rabellino12/go-playground/iohttp"
)

// Match is the struct for handling lobby actions on loop
type Match struct {
	IO     *iohttp.Client
	Logger *log.Logger
	Redis  *redis.Client
}

// RunLoop method acts as init for match loop handler
func (m *Match) RunLoop() {
	gamesQuery := m.Redis.LRange("games", 0, -1)
	games, err := gamesQuery.Result()
	if err != nil {
		m.Logger.Println("error getting game: ", err.Error())
		return
	}
	for _, game := range games {
		m.handleMatch(game)
	}
}

func (m *Match) handleMatch(game string) {
	channel := "match:" + game
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
