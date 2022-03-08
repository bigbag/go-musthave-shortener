package storage

type storageService struct {
	storageRepository StorageRepository
}

func NewStorageService(r StorageRepository) StorageService {
	return &storageService{
		storageRepository: r,
	}
}

func (s *storageService) Get(key string) (string, error) {
	return s.storageRepository.Get(key)
}

func (s *storageService) Save(key string, value string) (string, error) {
	return s.storageRepository.Save(key, value)
}
