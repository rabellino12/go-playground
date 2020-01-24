package loop

import (
	"fmt"
	"time"
)

// Handler is the loop methods handler
type Handler struct {
	quit       chan struct{}
	ticker     *time.Ticker
	controller Controller
}

// Initialize starts the loop and creates the quit chanel
func Initialize(controller Controller, hz time.Duration) {
	ticker := time.NewTicker(time.Second / hz)
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
