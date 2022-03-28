package repository

type Record struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	UserID        string `json:"user_id"`
	CorrelationID string `json:"correlation_id"`
}

type StorageRepository interface {
	GetByKey(key string) (*Record, error)
	GetByValue(value string) (*Record, error)
	GetAllByUserID(userID string) ([]*Record, error)
	Save(record *Record) (*Record, error)
	SaveBatchOfURL(records []*Record) error
	Status() error
	Close() error
}
