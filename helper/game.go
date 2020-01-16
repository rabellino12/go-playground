package helper

// GetPlayerInitialPosition calculates the player initial position on the map based on index
func GetPlayerInitialPosition(i int) string {
	x := 450 + 10*(i+1)
	y := 100 + 10*(i+1)
	return string(x) + "," + string(y)
}
