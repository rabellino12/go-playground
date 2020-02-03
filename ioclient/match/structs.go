package ioclient

import game "github.com/rabellino12/go-playground/db/collections"

// JoinEvent is the initial lobby join event when a game is created
type JoinEvent struct {
	Event   string        `json:"event"`
	Game    string        `json:"game"`
	Players []game.Player `json:"players"`
}

// Move interface
type Move struct {
	Action    string        `json:"action"`
	Jumping   bool          `json:"jumping"`
	Timestamp int           `json:"timestamp"`
	MatchID   string        `json:"matchId"`
	UserID    string        `json:"userId"`
	Position  game.Position `json:"position"`
}

// PlayerUserData contains important player info and action status
type PlayerUserData struct {
	game.Player `json:",inline"`
	Action      string `json:"action"`
	Jumping     bool   `json:"jumping"`
}
