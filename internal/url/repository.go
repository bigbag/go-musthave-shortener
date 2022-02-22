package url

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type storageRepository struct {
	shortStorage map[string]*URL
	fullStorage  map[string]*URL
}

func NewURLRepository() URLRepository {
	return &storageRepository{shortStorage: make(map[string]*URL), fullStorage: make(map[string]*URL)}
}

func (r *storageRepository) GetURL(shortURL string) (*URL, error) {
	url, ok := r.shortStorage[shortURL]
	if !ok {
		return nil, errors.New("NOT FOUND URL")
	}
	return url, nil
}

func (r *storageRepository) CreateURL(fullURL string) (*URL, error) {
	url, ok := r.fullStorage[fullURL]
	if ok {
		return url, nil
	}

	shortID := strings.Replace(uuid.New().String(), "-", "", -1)
	url = &URL{ShortID: shortID, FullURL: fullURL}

	r.fullStorage[fullURL] = url
	r.shortStorage[shortID] = url

	return url, nil
}
