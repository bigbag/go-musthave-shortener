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

func (r *fileRepository) Get(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.db[key]
	if !ok {
		return "", errors.New("NOT FOUND URL")
	}
	return record.Value, nil
}

func (r *fileRepository) Save(key string, value string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, record := range r.db {
		if record.Value == value {
			return record.Key, nil
		}
	}

	record := &Record{Key: key, Value: value}
	r.db[key] = record
	if err := r.producer.Write(record); err != nil {
		return "", err
	}

	return key, nil
}

func (r *fileRepository) Close() error {
	return r.producer.Close()
}
