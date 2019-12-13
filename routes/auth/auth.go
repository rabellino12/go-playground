package auth

import (
	"encoding/json"
	"github.com/rabellino12/go-playground/helper"
	"github.com/rabellino12/go-playground/ioclient/iohttp"
	"log"
	"net/http"
	"time"
)

// Handlers is a Struct that contains handler methods and shared server data
type Handlers struct {
	logger *log.Logger
	ioh    *iohttp.Client
}

// Auth struct containing authentication data
type Auth struct {
	Jwt string `json:"jwt"`
}

// User type structure, contains username
type User struct {
	Username string `json:"username"`
}

// Jwt handles the jwt authentication request route "/"
func (h *Handlers) Jwt(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.Fatalln("Error decoding body", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.ioh.GetJWT(user.Username)
	if err != nil {
		h.logger.Println("Error encoding token", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(&Auth{Jwt: token})
	if err != nil {
		h.logger.Println("Error encoding response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(js))
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
}

// NewHandlers returns a home page handlers struct
func NewHandlers(logger *log.Logger, ioh *iohttp.Client) *Handlers {
	return &Handlers{
		logger,
		ioh,
	}
}
