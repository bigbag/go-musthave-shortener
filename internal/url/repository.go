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

func (r *urlRepository) makeShortID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func (r *urlRepository) GetURL(shortID string) (*URL, error) {
	record, err := r.urlStorage.GetByKey(shortID)
	if err != nil {
		return nil, err
	}
	return &URL{ShortID: record.Key, FullURL: record.Value}, nil
}

func (r *urlRepository) CreateURL(fullURL string, userID string) (string, error) {
	shortID := r.makeShortID()
	record, err := r.urlStorage.Save(
		&repository.Record{
			Key:    shortID,
			Value:  fullURL,
			UserID: userID,
		},
	)

	switch err.(type) {
	case *storage.NotUniqueError:
		return record.Key, &NotUniqueURLError{}
	case nil:
		return shortID, nil
	default:
		return "", err
	}

}

func (r *urlRepository) CreateBatchOfURL(
	items BatchRequest,
	userID string,
) ([]*URL, error) {
	var (
		shortID string
		record  *repository.Record
		url     *URL
	)

	records := make([]*repository.Record, 0, 100)
	result := make([]*URL, 0, 100)
	for _, item := range items {
		shortID = r.makeShortID()
		record = &repository.Record{
			Key:           shortID,
			Value:         item.FullURL,
			UserID:        userID,
			CorrelationID: item.CorrelationID,
		}
		records = append(records, record)

		url = &URL{
			ShortID:       shortID,
			FullURL:       item.FullURL,
			CorrelationID: item.CorrelationID,
		}
		result = append(result, url)
	}

	err := r.urlStorage.SaveBatchOfURL(records)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *urlRepository) FindAllByUserID(userID string) ([]*URL, error) {
	records, err := r.urlStorage.GetAllByUserID(userID)
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

func (r *urlRepository) Status() error {
	return r.urlStorage.Status()
}

func (r *urlRepository) Close() error {
	return r.urlStorage.Shutdown()
}
