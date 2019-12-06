package server

import (
	"net/http"
	"time"
)

// NewServer Creats a new server
func NewServer(mux *http.ServeMux, addr string) *http.Server {
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv
}
