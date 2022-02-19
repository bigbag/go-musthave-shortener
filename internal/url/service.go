package url

type urlService struct {
	urlReposiory URLRepository
}

func NewURLService(r URLRepository) URLService {
	return &urlService{
		urlReposiory: r,
	}
}

func (s *urlService) FetchURL(shortURL string) (*URL, error) {
	return s.urlReposiory.GetURL(shortURL)
}

func (s *urlService) BuildURL(fullURL string) (*URL, error) {
	return s.urlReposiory.CreateURL(fullURL)
}
