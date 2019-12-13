package helper

import "net/http"

// EnableCors configures the requests
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
