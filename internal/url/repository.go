package url

import (
	"strings"

	"github.com/google/uuid"

	"github.com/bigbag/go-musthave-shortener/internal/storage"
)

type urlRepository struct {
	storageService storage.StorageService
}

func NewURLRepository(storageService storage.StorageService) URLRepository {
	return &urlRepository{storageService: storageService}
}

func (r *urlRepository) GetURL(shortID string) (*URL, error) {
	fullURL, err := r.storageService.Get(shortID)
	if err != nil {
		return nil, err
	}
	return &URL{ShortID: shortID, FullURL: fullURL}, nil
}

func (r *urlRepository) CreateURL(fullURL string) (*URL, error) {
	var err error

	shortID := strings.Replace(uuid.New().String(), "-", "", -1)
	shortID, err = r.storageService.Save(shortID, fullURL)
	if err != nil {
		return nil, err
	}

	return &URL{ShortID: shortID, FullURL: fullURL}, nil
}

func (r *urlRepository) Close() error {
	return r.storageService.Shutdown()
}
