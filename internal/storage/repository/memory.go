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

func (r *memoryRepository) GetByKey(key string) (*Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.db[key]
	if !ok {
		return record, errors.New("NOT FOUND URL")
	}
	return record, nil
}

func (r *memoryRepository) GetAllByUserID(userID string) ([]*Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Record, 0, 100)
	for _, record := range r.db {
		if record.UserID == userID {
			result = append(result, record)
		}
	}
	return result, nil
}

func (r *memoryRepository) Save(record *Record) (*Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, oldRecord := range r.db {
		if oldRecord.Value == record.Value {
			return oldRecord, nil
		}
	}

	r.db[record.Key] = record
	return record, nil
}

func (r *memoryRepository) Close() error {
	return nil
}
