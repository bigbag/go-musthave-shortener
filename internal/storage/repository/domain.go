package repository

type StorageRepository interface {
	GetByKey(key string) (*Record, error)
	GetByValue(value string) (*Record, error)
	GetAllByUserID(userID string) ([]*Record, error)
	Save(record *Record) error
	SaveBatchOfURL(records []*Record) error
	DeleteByUserID(userID string, keys []string) error
	Status() error
	Close() error
}
type Record struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	UserID        string `json:"user_id"`
	CorrelationID string `json:"correlation_id"`
	Removed       bool   `json:"-"`
}

func (r Record) IsOwnerAndExists(userID string) bool {
	return r.IsOwner(userID) && !r.Removed
}

func (r Record) IsOwner(userID string) bool {
	return r.UserID == userID
}
