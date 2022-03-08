package storage

type StorageRepository interface {
	Get(key string) (string, error)
	Save(key string, value string) (string, error)
}

type StorageService interface {
	Get(key string) (string, error)
	Save(key string, value string) (string, error)
}
