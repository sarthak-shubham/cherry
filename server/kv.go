package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sarthak-shubham/cherry/store"
)

func (s *Server) handleGet(w http.ResponseWriter, key string) {
	value, err := s.store.Get(key)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, kvResponse{
		Key:   key,
		Value: value,
	})
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request, key string) {
	var request setRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	var extra any

	err = decoder.Decode(&extra)
	if !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return
	}

	err = s.store.Set(key, request.Value)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, kvResponse{
		Key:   key,
		Value: request.Value,
	})
}

func (s *Server) handleDelete(w http.ResponseWriter, key string) {
	err := s.store.Delete(key)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, kvResponse{
		Key: key,
	})
}
