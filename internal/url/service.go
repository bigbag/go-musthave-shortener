package url

import (
	"fmt"
)

type urlService struct {
	urlRepository URLRepository
}

func NewURLService(r URLRepository) URLService {
	return &urlService{
		urlRepository: r,
	}
}

func (s *urlService) BuildURL(
	baseURL string,
	fullURL string,
	userID string,
) (*URL, error) {
	url, err := s.urlRepository.CreateURL(fullURL, userID)
	url.ShortURL = fmt.Sprintf("%s/%s", baseURL, url.ShortID)
	return url, err
}

func (s *urlService) BuildBatchOfURL(
	baseURL string,
	items BatchRequest,
	userID string,
) (BatchResponse, error) {
	urls, err := s.urlRepository.CreateBatchOfURL(items, userID)
	if err != nil {
		return nil, err
	}

	var (
		shortURL  string
		batchItem *BatchResponseItem
	)

	result := make([]*BatchResponseItem, 0, 100)
	for _, url := range urls {
		shortURL = fmt.Sprintf("%s/%s", baseURL, url.ShortID)
		batchItem = &BatchResponseItem{
			ShortURL:      shortURL,
			CorrelationID: url.CorrelationID,
		}
		result = append(result, batchItem)
	}
	return result, nil
}

func (s *urlService) FetchURL(shortID string) (*URL, error) {
	return s.urlRepository.GetURL(shortID)
}

func (s *urlService) FetchUserURLs(
	baseURL string,
	userID string,
) ([]*UserURL, error) {
	urls, err := s.urlRepository.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var shortURL string

	result := make([]*UserURL, 0, 100)
	for _, url := range urls {
		shortURL = fmt.Sprintf("%s/%s", baseURL, url.ShortID)
		result = append(result, &UserURL{ShortURL: shortURL, FullURL: url.FullURL})
	}
	return result, nil
}

func (s *urlService) Status() error {
	return s.urlRepository.Status()
}

func (s *urlService) Shutdown() error {
	return s.urlRepository.Close()
}
