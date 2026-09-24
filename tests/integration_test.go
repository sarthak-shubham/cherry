package tests

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/sarthak-shubham/cherry/recovery"
	"github.com/sarthak-shubham/cherry/store"
)

func TestStorePersistenceAndRecovery(t *testing.T) {
	walPath := filepath.Join(t.TempDir(), "data.jsonl")

	firstStore := store.NewStore(nil, walPath)

	err := firstStore.Set("name", "Sarthak")
	if err != nil {
		t.Fatalf("Set returned an unexpected error: %v", err)
	}

	err = firstStore.Set("city", "Delhi")
	if err != nil {
		t.Fatalf("Set returned an unexpected error: %v", err)
	}

	err = firstStore.Set("name", "Rahul")
	if err != nil {
		t.Fatalf("Set returned an unexpected error: %v", err)
	}

	err = firstStore.Delete("city")
	if err != nil {
		t.Fatalf("Delete returned an unexpected error: %v", err)
	}

	recoveredData, err := recovery.Recover(walPath)
	if err != nil {
		t.Fatalf("Recover returned an unexpected error: %v", err)
	}

	secondStore := store.NewStore(recoveredData, walPath)

	value, err := secondStore.Get("name")
	if err != nil {
		t.Fatalf("Get returned an unexpected error: %v", err)
	}

	if value != "Rahul" {
		t.Fatalf("name = %q, want %q", value, "Rahul")
	}

	_, err = secondStore.Get("city")
	if !errors.Is(err, store.ErrKeyNotFound) {
		t.Fatalf("Get returned error %v, want %v", err, store.ErrKeyNotFound)
	}

	err = secondStore.Set("language", "Go")
	if err != nil {
		t.Fatalf("Set after recovery returned an unexpected error: %v", err)
	}

	recoveredData, err = recovery.Recover(walPath)
	if err != nil {
		t.Fatalf("second Recover returned an unexpected error: %v", err)
	}

	thirdStore := store.NewStore(recoveredData, walPath)

	value, err = thirdStore.Get("language")
	if err != nil {
		t.Fatalf("Get after second recovery returned an unexpected error: %v", err)
	}

	if value != "Go" {
		t.Fatalf("language = %q, want %q", value, "Go")
	}
}
