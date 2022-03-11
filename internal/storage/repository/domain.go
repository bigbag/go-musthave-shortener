package repository

type Record struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StorageRepository interface {
	Get(key string) (string, error)
	Save(key string, value string) (string, error)
	Close() error
}
