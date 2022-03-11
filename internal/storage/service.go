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

func (s *StorageService) Get(key string) (string, error) {
	return s.storageRepository.Get(key)
}

func (s *StorageService) Save(key string, value string) (string, error) {
	return s.storageRepository.Save(key, value)
}

func (s *StorageService) Shutdown() error {
	return s.storageRepository.Close()
}
