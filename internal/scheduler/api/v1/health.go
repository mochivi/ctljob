package v1

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, "ok")
}

func Ready(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []byte("ready"))
}
