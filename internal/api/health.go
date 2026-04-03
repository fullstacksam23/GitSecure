package api

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Status OK"))
}
