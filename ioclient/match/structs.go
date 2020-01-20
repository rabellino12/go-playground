package match

import game "github.com/rabellino12/go-playground/db/collections"

// JoinEvent is the initial lobby join event when a game is created
type JoinEvent struct {
	Event   string        `json:"event"`
	Game    string        `json:"game"`
	Players []game.Player `json:"players"`
}
