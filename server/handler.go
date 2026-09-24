package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

type setRequest struct {
	Value string `json:"value"`
}

type kvResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleKV(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/kv/")

	if key == "" || strings.Contains(key, "/") {
		writeError(w, http.StatusBadRequest, "invalid key")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleGet(w, key)

	case http.MethodPut:
		s.handleSet(w, r, key)

	case http.MethodDelete:
		s.handleDelete(w, key)

	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{
		Error: message,
	})
}
