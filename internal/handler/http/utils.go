package http

import (
    "encoding/json"
    "net/http"
)

type ErrorResponse struct {
    Error  string `json:"error"`
    Status int    `json:"status"`
}

func respondError(w http.ResponseWriter, status int, msg string) {
    respondJSON(w, status, ErrorResponse{
        Error:  msg,
        Status: status,
    })
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func decodeJSON(r *http.Request, v any) error {
    return json.NewDecoder(r.Body).Decode(v)
}