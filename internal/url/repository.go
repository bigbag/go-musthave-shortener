package url

import (
	"github.com/google/uuid"
	"strings"

	"github.com/bigbag/go-musthave-shortener/internal/storage"
	"github.com/bigbag/go-musthave-shortener/internal/storage/repository"
)

type urlRepository struct {
	s storage.StorageService
}

func NewURLRepository(s storage.StorageService) URLRepository {
	return &urlRepository{s: s}
}

func (r *urlRepository) makeShortID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func (r *urlRepository) GetURL(shortID string) (*URL, error) {
	record, err := r.s.GetByKey(shortID)
	if err != nil {
		return nil, err
	}
	return &URL{
		ShortID: record.Key,
		FullURL: record.Value,
		Removed: record.Removed,
	}, nil
}

func (r *urlRepository) CreateURL(fullURL string, userID string) (string, error) {
	record, err := r.s.Save(
		&repository.Record{
			Key:     r.makeShortID(),
			Value:   fullURL,
			UserID:  userID,
			Removed: false,
		},
	)

	switch err.(type) {
	case *storage.NotUniqueError:
		return record.Key, &NotUniqueURLError{}
	case nil:
		return record.Key, nil
	default:
		return "", err
	}

}

func (r *urlRepository) CreateBatchOfURL(
	items BatchRequest,
	userID string,
) ([]*URL, error) {
	var (
		record *repository.Record
		url    *URL
	)

	recordsForSave := make([]*repository.Record, 0, 100)
	for _, item := range items {
		record = &repository.Record{
			Key:           r.makeShortID(),
			Value:         item.FullURL,
			UserID:        userID,
			Removed:       false,
			CorrelationID: item.CorrelationID,
		}
		recordsForSave = append(recordsForSave, record)

	}

	records, err := r.s.SaveBatchOfRecord(recordsForSave)
	if err != nil {
		return nil, err
	}

	result := make([]*URL, 0, 100)
	for _, record := range records {
		url = &URL{
			ShortID:       record.Key,
			FullURL:       record.Value,
			CorrelationID: record.CorrelationID,
			Removed:       record.Removed,
		}
		result = append(result, url)
	}

	return result, nil
}

func (r *urlRepository) FindAllByUserID(userID string) ([]*URL, error) {
	records, err := r.s.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var url *URL

	result := make([]*URL, 0, 100)
	for _, record := range records {
		url = &URL{ShortID: record.Key, FullURL: record.Value}
		result = append(result, url)
	}
	return result, nil
}

func (r *urlRepository) DeleteUserURLs(userID string, shortIDs []string) error {
	return r.s.DeleteByUserID(userID, shortIDs)
}

func (r *urlRepository) Status() error {
	return r.s.Status()
}

func (r *urlRepository) Close() error {
	return r.s.Shutdown()
}
