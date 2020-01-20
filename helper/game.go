package helper

import "strconv"

// GetPlayerInitialPosition calculates the player initial position on the map based on index
func GetPlayerInitialPosition(i int) string {
	x := 100 + 30*(i+1)
	y := 450
	return strconv.FormatInt(int64(x), 10) + "," + strconv.FormatInt(int64(y), 10)
}
