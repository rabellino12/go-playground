package match

import game "github.com/rabellino12/go-playground/db/collections"

// JoinEvent is the initial lobby join event when a game is created
type JoinEvent struct {
	Event   string        `json:"event"`
	Game    string        `json:"game"`
	Players []game.Player `json:"players"`
}

// MoveEvent is the move IO event structure
type MoveEvent struct {
	Action    string `json:"action"`
	MatchID   string `json:"matchId"`
	Timestamp string `json:"timestamp"`
	UserID    string `json:"userId"`
}
