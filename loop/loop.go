package loop

import (
	"fmt"
	"time"
)

const emptyLobbies string = "emptyLobbies"
const loadingLobbies string = "loadingLobbies"
const fullLobbies string = "fullLobbies"

// Handler is the loop methods handler
type Handler struct {
	quit       chan struct{}
	ticker     *time.Ticker
	controller Controller
}

// Initialize starts the loop and creates the quit chanel
func Initialize(controller Controller) {
	ticker := time.NewTicker(time.Second / 20)
	quit := make(chan struct{})
	h := &Handler{quit, ticker, controller}
	go h.loop()
	// close(quit)
}

func (h *Handler) loop() {
	for {
		select {
		case <-h.ticker.C:
			h.controller.RunLoop()
		case <-h.quit:
			fmt.Println("ticker stopped")
			return
		}
	}
}
