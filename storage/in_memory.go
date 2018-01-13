package storage

import (
	"errors"
	"sync"
)

type memRepository struct {
	sync.RWMutex
	store map[string]Result
}

// NewInMemoryRepository returns a newly initialized in-memory store
func NewInMemoryRepository() (API, error) {
	return &memRepository{
		store: make(map[string]Result),
	}, nil
}

func (r *memRepository) Get(id string) (*Result, error) {
	if id == "" {
		return nil, errors.New("ID is not allowed to be empty")
	}

	r.RLock()
	if val, ok := r.store[id]; ok {
		r.RUnlock()
		return &val, nil
	}
	r.RUnlock()
	return nil, nil
}

func (r *memRepository) Set(id string, res Result) error {
	if id == "" {
		return errors.New("ID is not allowed to be empty")
	}
	r.Lock()
	r.store[id] = res
	r.Unlock()
	return nil
}
