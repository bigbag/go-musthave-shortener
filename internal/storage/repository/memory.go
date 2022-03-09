package repository

import (
	"errors"
	"sync"
)

type memoryRepository struct {
	mu *sync.RWMutex
	db map[string]*Record
}

func NewMemoryRepository() (StorageRepository, error) {
	repo := &memoryRepository{
		mu: &sync.RWMutex{},
		db: make(map[string]*Record),
	}
	return repo, nil
}

func (r *memoryRepository) Get(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.db[key]
	if !ok {
		return "", errors.New("NOT FOUND URL")
	}
	return record.Value, nil
}

func (r *memoryRepository) Save(key string, value string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, record := range r.db {
		if record.Value == value {
			return record.Key, nil
		}
	}

	r.db[key] = &Record{Key: key, Value: value}
	return key, nil
}

func (r *memoryRepository) Close() error {
	return nil
}
