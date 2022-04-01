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
		return record, errors.New("not found url")
	}
	return record, nil
}

func (r *memoryRepository) GetByValue(value string) (*Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, record := range r.db {
		if record.Value == value {
			return record, nil
		}
	}
	return nil, nil
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

func (r *memoryRepository) Save(record *Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[record.Key] = record
	return nil
}

func (r *memoryRepository) SaveBatchOfURL(records []*Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, record := range records {
		r.db[record.Key] = record
	}
	return nil
}

func (r *memoryRepository) DeleteByUserID(userID string, keys []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, key := range keys {
		record, ok := r.db[key]
		if !ok {
			continue
		}
		if !record.IsOwner(userID) {
			continue
		}
		record.Removed = true
		r.db[key] = record
	}
	return nil
}

func (r *memoryRepository) Status() error {
	return nil
}

func (r *memoryRepository) Close() error {
	return nil
}
