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
	cfg *config.Storage
	r   repository.StorageRepository
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

	service := StorageService{r: r, cfg: cfg}
	if err != nil {
		return service, err
	}

	return service, nil
}

func (s *StorageService) GetByKey(key string) (*repository.Record, error) {
	return s.r.GetByKey(key)
}

func (s *StorageService) GetAllByUserID(userID string) ([]*repository.Record, error) {
	return s.r.GetAllByUserID(userID)
}

func (s *StorageService) Save(record *repository.Record) (*repository.Record, error) {
	oldRecord, err := s.r.GetByValue(record.Value)
	if err != nil {
		return oldRecord, err
	}

	if oldRecord != nil {
		return oldRecord, &NotUniqueError{}
	}

	err = s.r.Save(record)
	return record, err
}

func (s *StorageService) SaveBatchOfRecord(
	records []*repository.Record,
) ([]*repository.Record, error) {
	var oldRecord *repository.Record

	newRecords := make([]*repository.Record, 0, 100)
	result := make([]*repository.Record, 0, 100)
	for _, record := range records {
		oldRecord, _ = s.r.GetByValue(record.Value)
		if oldRecord == nil {
			newRecords = append(newRecords, record)
			result = append(result, record)
		} else {
			result = append(result, oldRecord)
		}
	}
	err := s.r.SaveBatchOfURL(newRecords)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *StorageService) DeleteByUserID(userID string, shortIDs []string) error {
	return s.r.DeleteByUserID(userID, shortIDs)
}

func (s *StorageService) Status() error {
	return s.r.Status()
}

func (s *StorageService) Shutdown() error {
	return s.r.Close()
}
