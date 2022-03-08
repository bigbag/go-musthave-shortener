package storage

import (
	"errors"
)

type memoryRepository struct {
	memoryStorage map[string]string
}

func NewMemoryStorage() StorageRepository {
	return &memoryRepository{memoryStorage: make(map[string]string)}
}

func (r *memoryRepository) Get(key string) (string, error) {
	value, ok := r.memoryStorage[key]
	if !ok {
		return "", errors.New("NOT FOUND URL")
	}
	return value, nil
}

func (r *memoryRepository) Save(key string, value string) (string, error) {
	for oldKey, oldValue := range r.memoryStorage {
		if oldValue == value {
			return oldKey, nil
		}
	}

	r.memoryStorage[key] = value
	return key, nil
}
