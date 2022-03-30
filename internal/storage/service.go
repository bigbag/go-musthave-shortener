package storage

import (
	"context"

	"github.com/bigbag/go-musthave-shortener/internal/config"
	"github.com/bigbag/go-musthave-shortener/internal/storage/repository"
)

type NotUniqueError struct{}

func (e *NotUniqueError) Error() string {
	return "not unique value"
}

type StorageService struct {
	cfg               *config.Storage
	storageRepository repository.StorageRepository
}

func NewStorageService(ctx context.Context, cfg *config.Storage) (StorageService, error) {
	var (
		r   repository.StorageRepository
		err error
	)

	if cfg.DatabaseDSN != "" {
		r, err = repository.NewPGRepository(ctx, cfg.DatabaseDSN, cfg.ConnectionTimeout)
	} else {
		if cfg.FileStoragePath != "" {
			r, err = repository.NewFileRepository(cfg.FileStoragePath)
		} else {
			r, err = repository.NewMemoryRepository()
		}
	}

	service := StorageService{storageRepository: r, cfg: cfg}
	if err != nil {
		return service, err
	}

	return service, nil
}

func (s *StorageService) GetByKey(key string) (*repository.Record, error) {
	return s.storageRepository.GetByKey(key)
}

func (s *StorageService) GetAllByUserID(userID string) ([]*repository.Record, error) {
	return s.storageRepository.GetAllByUserID(userID)
}

func (s *StorageService) Save(record *repository.Record) (*repository.Record, error) {
	oldRecord, err := s.storageRepository.GetByValue(record.Value)
	if err != nil {
		return oldRecord, err
	}

	if oldRecord != nil {
		return oldRecord, &NotUniqueError{}
	}

	err = s.storageRepository.Save(record)
	return record, err
}

func (s *StorageService) SaveBatchOfURL(records []*repository.Record) error {
	return s.storageRepository.SaveBatchOfURL(records)
}

func (s *StorageService) Status() error {
	return s.storageRepository.Status()
}

func (s *StorageService) Shutdown() error {
	return s.storageRepository.Close()
}
