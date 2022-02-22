package url

type URL struct {
	ShortID  string `json:"-"`
	FullURL  string `json:"url"`
	ShortURL string `json:"-"`
}

type URLRepository interface {
	GetURL(shortID string) (*URL, error)
	CreateURL(fullURL string) (*URL, error)
}

type URLService interface {
	FetchURL(shortID string) (*URL, error)
	BuildURL(baseURL string, fullURL string) (*URL, error)
}
