package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestHTTPServer(t *testing.T, handler http.Handler) {
	t.Helper()

	testServer := httptest.NewServer(handler)
	t.Cleanup(testServer.Close)

	originalServerURL := serverURL
	serverURL = testServer.URL

	t.Cleanup(func() {
		serverURL = originalServerURL
	})
}

func TestSetKey(t *testing.T) {
	requestReceived := false

	newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true

		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want %s", r.Method, http.MethodPut)
		}

		if r.URL.Path != "/kv/name" {
			t.Errorf("path = %s, want /kv/name", r.URL.Path)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}

		var request setRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			t.Errorf("failed to decode request body: %v", err)
			return
		}

		if request.Value != "Sarthak" {
			t.Errorf("value = %q, want %q", request.Value, "Sarthak")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kvResponse{
			Key:   "name",
			Value: "Sarthak",
		})
	}))

	setKey("name", "Sarthak")

	if !requestReceived {
		t.Fatal("server did not receive request")
	}
}

func TestGetKey(t *testing.T) {
	requestReceived := false

	newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true

		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want %s", r.Method, http.MethodGet)
		}

		if r.URL.Path != "/kv/name" {
			t.Errorf("path = %s, want /kv/name", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kvResponse{
			Key:   "name",
			Value: "Sarthak",
		})
	}))

	getKey("name")

	if !requestReceived {
		t.Fatal("server did not receive request")
	}
}

func TestDeleteKey(t *testing.T) {
	requestReceived := false

	newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true

		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want %s", r.Method, http.MethodDelete)
		}

		if r.URL.Path != "/kv/name" {
			t.Errorf("path = %s, want /kv/name", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kvResponse{
			Key: "name",
		})
	}))

	deleteKey("name")

	if !requestReceived {
		t.Fatal("server did not receive request")
	}
}

func TestGetKeyHandlesAPIError(t *testing.T) {
	requestReceived := false

	newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(errorResponse{
			Error: "key not found",
		})
	}))

	getKey("missing")

	if !requestReceived {
		t.Fatal("server did not receive request")
	}
}
