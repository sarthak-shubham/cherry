package store

import "sync"

type Store struct {
	data    map[string]string
	mu      sync.RWMutex
	walPath string
}

func NewStore(data map[string]string, walPath string) *Store {
	if data == nil {
		data = make(map[string]string)
	}

	return &Store{
		data:    data,
		walPath: walPath,
	}
}
