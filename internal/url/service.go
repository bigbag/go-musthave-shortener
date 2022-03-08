package url

import "fmt"

type urlService struct {
	urlReposiory URLRepository
}

func NewURLService(r URLRepository) URLService {
	return &urlService{
		urlReposiory: r,
	}
}

func (s *urlService) FetchURL(shortID string) (*URL, error) {
	return s.urlReposiory.GetURL(shortID)
}

func (s *urlService) BuildURL(baseURL string, fullURL string) (*URL, error) {
	url, err := s.urlReposiory.CreateURL(fullURL)
	url.ShortURL = fmt.Sprintf("%s/%s", baseURL, url.ShortID)
	return url, err
}

func (s *urlService) Shutdown() error {
	return s.urlReposiory.Close()
}
