package loop

import (
	"fmt"
	"time"
)

// Handler is the loop methods handler
type Handler struct {
	quit       chan struct{}
	ticker     *time.Ticker
	Controller Controller
	lastLoop   time.Time
}

// Initialize starts the loop and creates the quit chanel
func Initialize(controller Controller, hz time.Duration) {
	ticker := time.NewTicker(time.Second / hz)
	quit := make(chan struct{})
	h := &Handler{quit, ticker, controller, time.Time{}}
	go h.loop()
	// close(quit)
}

func (h *Handler) loop() {
	h.lastLoop = time.Now()
	for {
		select {
		case <-h.ticker.C:
			now := time.Now()
			elapsed := now.Sub(h.lastLoop)
			h.Controller.RunLoop(elapsed.Milliseconds())
			h.lastLoop = time.Now()
		case <-h.quit:
			fmt.Println("ticker stopped")
			return
		}
	}
}
