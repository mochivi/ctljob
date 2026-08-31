package v1

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, detail string) {
	_ = writeJSON(w, status, errorResponse{
		Error:  http.StatusText(status),
		Detail: detail,
	})
}

func writeInternalError(w http.ResponseWriter) {
	_ = writeJSON(w, http.StatusInternalServerError, errorResponse{
		Error:  http.StatusText(http.StatusInternalServerError),
		Detail: "something went wrong",
	})
}
