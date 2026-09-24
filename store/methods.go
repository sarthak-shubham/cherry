package store

import (
	"errors"

	"github.com/sarthak-shubham/cherry/persistence"
)

var ErrKeyNotFound = errors.New("key not found")

func (s *Store) Set(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := persistence.WALRecord{
		Operation: "set",
		Key:       key,
		Value:     value,
	}

	err := persistence.AppendLog(s.walPath, record)
	if err != nil {
		return err
	}

	s.data[key] = value

	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	if !exists {
		return "", ErrKeyNotFound
	}

	return value, nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]
	if !exists {
		return ErrKeyNotFound
	}

	record := persistence.WALRecord{
		Operation: "delete",
		Key:       key,
	}

	err := persistence.AppendLog(s.walPath, record)
	if err != nil {
		return err
	}

	delete(s.data, key)

	return nil
}
