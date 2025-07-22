package store

import (
	"main/pkg/models"
	"sync"
)

type CodeRepository interface {
	Save(code models.StoredCode) error
	Get(userID string) (models.StoredCode, bool, error)
}

type InMemoryStore struct {
	mu    sync.RWMutex
	codes map[string]models.StoredCode
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		codes: make(map[string]models.StoredCode),
	}
}

func (s *InMemoryStore) Save(code models.StoredCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[code.UserID] = code
	return nil
}

func (s *InMemoryStore) Get(userID string) (models.StoredCode, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, found := s.codes[userID]
	return code, found, nil
}
