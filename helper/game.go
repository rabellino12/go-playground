package helper

import game "github.com/rabellino12/go-playground/db/collections"

// GetPlayerInitialPosition calculates the player initial position on the map based on index
func GetPlayerInitialPosition(i int) game.Position {
	x := float64(100 + 30*(i+1))
	y := float64(450)
	return game.Position{X: x, Y: y}
}
