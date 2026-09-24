package store

import (
	"errors"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()

	walPath := filepath.Join(t.TempDir(), "data.jsonl")

	return NewStore(nil, walPath)
}

func TestSetAndGet(t *testing.T) {
	kvStore := newTestStore(t)

	err := kvStore.Set("name", "Sarthak")
	if err != nil {
		t.Fatalf("Set returned an unexpected error: %v", err)
	}

	value, err := kvStore.Get("name")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value != "Sarthak" {
		t.Fatalf("Get returned %q, want %q", value, "Sarthak")
	}
}

func TestSetOverwritesExistingValue(t *testing.T) {
	kvStore := newTestStore(t)

	err := kvStore.Set("name", "Sarthak")
	if err != nil {
		t.Fatalf("first Set returned an unexpected error: %v", err)
	}

	err = kvStore.Set("name", "Rahul")
	if err != nil {
		t.Fatalf("second Set returned an unexpected error: %v", err)
	}

	value, err := kvStore.Get("name")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value != "Rahul" {
		t.Fatalf("Get returned %q, want %q", value, "Rahul")
	}
}

func TestGetMissingKey(t *testing.T) {
	kvStore := newTestStore(t)

	_, err := kvStore.Get("missing")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("Get returned error %v, want %v", err, ErrKeyNotFound)
	}
}

func TestDelete(t *testing.T) {
	kvStore := newTestStore(t)

	err := kvStore.Set("name", "Sarthak")
	if err != nil {
		t.Fatalf("Set returned an unexpected error: %v", err)
	}

	err = kvStore.Delete("name")
	if err != nil {
		t.Fatalf("Delete returned an unexpected error: %v", err)
	}

	_, err = kvStore.Get("name")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("Get returned error %v, want %v", err, ErrKeyNotFound)
	}
}

func TestDeleteMissingKey(t *testing.T) {
	kvStore := newTestStore(t)

	err := kvStore.Delete("missing")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("Delete returned error %v, want %v", err, ErrKeyNotFound)
	}
}

func TestSetDoesNotChangeStoreWhenWALWriteFails(t *testing.T) {
	walPath := filepath.Join(t.TempDir(), "missing", "data.jsonl")
	kvStore := NewStore(nil, walPath)

	err := kvStore.Set("name", "Sarthak")
	if err == nil {
		t.Fatal("Set returned nil, want WAL write error")
	}

	_, err = kvStore.Get("name")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("Get returned error %v, want %v", err, ErrKeyNotFound)
	}
}

func TestDeleteDoesNotChangeStoreWhenWALWriteFails(t *testing.T) {
	kvStore := newTestStore(t)

	err := kvStore.Set("name", "Sarthak")
	if err != nil {
		t.Fatalf("Set returned an unexpected error: %v", err)
	}

	kvStore.walPath = filepath.Join(t.TempDir(), "missing", "data.jsonl")

	err = kvStore.Delete("name")
	if err == nil {
		t.Fatal("Delete returned nil, want WAL write error")
	}

	value, err := kvStore.Get("name")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value != "Sarthak" {
		t.Fatalf("Get returned %q, want %q", value, "Sarthak")
	}
}
