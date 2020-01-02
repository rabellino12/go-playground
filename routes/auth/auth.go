package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rabellino12/go-playground/helper"
	"github.com/rabellino12/go-playground/iohttp"
)

// Handlers is a Struct that contains handler methods and shared server data
type Handlers struct {
	logger *log.Logger
	ioh    *iohttp.Client
}

// Auth struct containing authentication data
type Auth struct {
	Token string `json:"token"`
}

// User type structure, contains username
type User struct {
	User string `json:"user"`
}

// PrivSubscription is the centrifuge auth subscription body
type PrivSubscription struct {
	Client   string   `json:"client"`
	Channels []string `json:"channels"`
}

type channel struct {
	Channel string `json:"channel"`
	Token   string `json:"token"`
}

type privSubscriptionResponse struct {
	Channels []channel `json:"channels"`
}

// Jwt handles the jwt authentication request route "/"
func (h *Handlers) Jwt(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.Println("Error decoding body", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := helper.GetJWT(user.User)
	if err != nil {
		h.logger.Println("Error encoding token", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(&Auth{Token: token})
	if err != nil {
		h.logger.Println("Error encoding response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(js))
}

// Centrifuge handles the centrifuge private channel subscription token request "/"
func (h *Handlers) Centrifuge(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	var req PrivSubscription

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Println("Error decoding body", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := helper.GetSubscriptionJWT(req.Client, req.Channels[0])
	if err != nil {
		h.logger.Println("Error encoding token", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(&privSubscriptionResponse{[]channel{channel{req.Channels[0], token}}})
	if err != nil {
		h.logger.Println("Error encoding response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// Logger is the Home logging middleware
func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.logger.Println("Processing auth request")
		defer h.logger.Printf("Request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

// SetupRoutes creates all home related routes
func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/jwt", h.Logger(h.Jwt))
	mux.HandleFunc("/auth/centrifuge", h.Logger(h.Centrifuge))
}

// NewHandlers returns a home page handlers struct
func NewHandlers(logger *log.Logger, ioh *iohttp.Client) *Handlers {
	return &Handlers{
		logger,
		ioh,
	}
}
