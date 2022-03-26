package repository

import (
	"errors"
	"sync"
)

type fileRepository struct {
	mu       *sync.RWMutex
	db       map[string]*Record
	producer *producer
}

func NewFileRepository(fileStoragePath string) (StorageRepository, error) {
	producer, err := NewProducer(fileStoragePath)
	if err != nil {
		return nil, err
	}

	consumer, err := NewConsumer(fileStoragePath)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	db, err := consumer.ReadAll()
	if err != nil {
		return nil, err
	}

	repo := &fileRepository{
		mu:       &sync.RWMutex{},
		db:       db,
		producer: producer,
	}

	return repo, nil
}

func (r *fileRepository) GetByKey(key string) (*Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.db[key]
	if !ok {
		return record, errors.New("NOT FOUND URL")
	}
	return record, nil
}

func (r *fileRepository) GetAllByUserID(userID string) ([]*Record, error) {
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

func (r *fileRepository) Save(record *Record) (*Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, oldRecord := range r.db {
		if oldRecord.Value == record.Value {
			return oldRecord, nil
		}
	}

	r.db[record.Key] = record
	if err := r.producer.Write(record); err != nil {
		return record, err
	}

	return record, nil
}

func (r *fileRepository) Close() error {
	return r.producer.Close()
}
