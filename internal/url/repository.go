package url

import (
	"github.com/google/uuid"
	"strings"

	"github.com/bigbag/go-musthave-shortener/internal/storage"
	"github.com/bigbag/go-musthave-shortener/internal/storage/repository"
)

type urlRepository struct {
	urlStorage storage.StorageService
}

func NewURLRepository(urlStorage storage.StorageService) URLRepository {
	return &urlRepository{urlStorage: urlStorage}
}

func (r *urlRepository) GetURL(shortID string) (*URL, error) {
	record, err := r.urlStorage.GetByKey(shortID)
	if err != nil {
		return nil, err
	}
	return &URL{ShortID: record.Key, FullURL: record.Value}, nil
}

func (r *urlRepository) CreateURL(fullURL string, userID string) (*URL, error) {
	var err error

	shortID := strings.Replace(uuid.New().String(), "-", "", -1)
	record, err := r.urlStorage.Save(
		&repository.Record{Key: shortID, Value: fullURL, UserID: userID},
	)
	if err != nil {
		return nil, err
	}

	return &URL{ShortID: record.Key, FullURL: record.Value}, nil
}

func (r *urlRepository) FindAllByUserID(userID string) ([]*URL, error) {
	var url *URL

	records, err := r.urlStorage.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]*URL, 0, 100)
	for _, record := range records {
		url = &URL{ShortID: record.Key, FullURL: record.Value}
		result = append(result, url)
	}
	return result, nil
}

func (r *urlRepository) Close() error {
	return r.urlStorage.Shutdown()
}
