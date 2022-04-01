package url

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type urlService struct {
	l             logrus.FieldLogger
	r             URLRepository
	p             *TaskPool
	deleteTimeout time.Duration
}

func NewURLService(l logrus.FieldLogger, r URLRepository, p *TaskPool) URLService {
	return &urlService{l: l, r: r, p: p}
}

func (s *urlService) BuildURL(
	baseURL string,
	fullURL string,
	userID string,
) (string, error) {
	shortID, err := s.r.CreateURL(fullURL, userID)
	return fmt.Sprintf("%s/%s", baseURL, shortID), err
}

func (s *urlService) BuildBatchOfURL(
	baseURL string,
	items BatchRequest,
	userID string,
) (BatchResponse, error) {
	urls, err := s.r.CreateBatchOfURL(items, userID)
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
	return s.r.GetURL(shortID)
}

func (s *urlService) FetchUserURLs(
	baseURL string,
	userID string,
) ([]*UserURL, error) {
	urls, err := s.r.FindAllByUserID(userID)
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

func (s *urlService) DeleteUserURLs(userID string, shortIDs []string) error {
	return s.p.Push(userID, shortIDs)
}

func (s *urlService) Status() error {
	return s.r.Status()
}

func (s *urlService) Shutdown() error {
	return s.r.Close()
}
