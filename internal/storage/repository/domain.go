package repository

type Record struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	UserID string `json:"user_id"`
}

type StorageRepository interface {
	GetByKey(key string) (*Record, error)
	GetAllByUserID(userID string) ([]*Record, error)
	Save(record *Record) (*Record, error)
	Close() error
}
