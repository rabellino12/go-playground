package home

import (
	"log"
	"net/http"
	"time"
)

const message = "Hello World"

// Handlers is a Struct that contains handler methods and shared server data
type Handlers struct {
	logger *log.Logger
}

// Home is a handler funciton for home route "/"
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

// Logger is the Home logging middleware
func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.logger.Println("Processing home request")
		defer h.logger.Printf("Request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

// SetupRoutes creates all home related routes
func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Logger(h.Home))
}

// NewHandlers returns a home page handlers struct
func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}
