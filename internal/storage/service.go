package storage

import (
	"github.com/bigbag/go-musthave-shortener/internal/config"
	"github.com/bigbag/go-musthave-shortener/internal/storage/repository"
)

type StorageService struct {
	storageRepository repository.StorageRepository
}

func NewStorageService(cfg *config.Config) (StorageService, error) {
	var (
		r   repository.StorageRepository
		err error
	)

	if cfg.FileStoragePath != "" {
		r, err = repository.NewFileRepository(cfg.FileStoragePath)
	} else {
		r, err = repository.NewMemoryRepository()
	}
	service := StorageService{storageRepository: r}
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
	return s.storageRepository.Save(record)
}

func (s *StorageService) Shutdown() error {
	return s.storageRepository.Close()
}
