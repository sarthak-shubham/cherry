package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sarthak-shubham/cherry/store"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()

	walPath := filepath.Join(t.TempDir(), "data.jsonl")
	kvStore := store.NewStore(nil, walPath)

	return NewServer(kvStore)
}

func TestSetHandler(t *testing.T) {
	server := newTestServer(t)

	request := httptest.NewRequest(
		http.MethodPut,
		"/kv/name",
		strings.NewReader(`{"value":"Sarthak"}`),
	)
	response := httptest.NewRecorder()

	server.handleKV(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}

	var result kvResponse

	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Key != "name" {
		t.Fatalf("key = %q, want %q", result.Key, "name")
	}

	if result.Value != "Sarthak" {
		t.Fatalf("value = %q, want %q", result.Value, "Sarthak")
	}
}

func TestGetHandler(t *testing.T) {
	server := newTestServer(t)

	err := server.store.Set("name", "Sarthak")
	if err != nil {
		t.Fatalf("Set returned an unexpected error: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/kv/name",
		nil,
	)
	response := httptest.NewRecorder()

	server.handleKV(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}

	var result kvResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Key != "name" {
		t.Fatalf("key = %q, want %q", result.Key, "name")
	}

	if result.Value != "Sarthak" {
		t.Fatalf("value = %q, want %q", result.Value, "Sarthak")
	}
}

func TestGetMissingKey(t *testing.T) {
	server := newTestServer(t)

	request := httptest.NewRequest(
		http.MethodGet,
		"/kv/missing",
		nil,
	)
	response := httptest.NewRecorder()

	server.handleKV(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func TestDeleteHandler(t *testing.T) {
	server := newTestServer(t)

	err := server.store.Set("name", "Sarthak")
	if err != nil {
		t.Fatalf("Set returned an unexpected error: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodDelete,
		"/kv/name",
		nil,
	)
	response := httptest.NewRecorder()

	server.handleKV(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}

	_, err = server.store.Get("name")
	if err != store.ErrKeyNotFound {
		t.Fatalf("Get returned error %v, want %v", err, store.ErrKeyNotFound)
	}
}

func TestDeleteMissingKey(t *testing.T) {
	server := newTestServer(t)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/kv/missing",
		nil,
	)
	response := httptest.NewRecorder()

	server.handleKV(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func TestSetRejectsInvalidJSON(t *testing.T) {
	server := newTestServer(t)

	request := httptest.NewRequest(
		http.MethodPut,
		"/kv/name",
		strings.NewReader(`not json`),
	)
	response := httptest.NewRecorder()

	server.handleKV(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestHandlerRejectsUnsupportedMethod(t *testing.T) {
	server := newTestServer(t)

	request := httptest.NewRequest(
		http.MethodPost,
		"/kv/name",
		nil,
	)
	response := httptest.NewRecorder()

	server.handleKV(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandlerRejectsInvalidKey(t *testing.T) {
	server := newTestServer(t)

	request := httptest.NewRequest(
		http.MethodGet,
		"/kv/",
		nil,
	)
	response := httptest.NewRecorder()

	server.handleKV(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
