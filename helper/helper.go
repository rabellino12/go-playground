package helper

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// CentrifugoSecret is the centrifugo server instance secret key
const CentrifugoSecret = "some-centrifugo-secret-key"

// EnableCors configures the requests
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// GetJWT receives a username and returns a jwt
func GetJWT(userName string) (string, error) {
	claims := jwt.MapClaims{"sub": userName}
	// if exp > 0 {
	// 	claims["exp"] = exp
	// }
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(CentrifugoSecret))
	if err != nil {
		return "", err
	}
	return t, err
}

// GetSubscriptionJWT receives a username and channel and returns a jwt
func GetSubscriptionJWT(client string, channel string) (string, error) {
	signingKey := []byte(CentrifugoSecret)

	// Create the Claims
	claims := jwt.MapClaims{"channel": channel, "client": client}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
