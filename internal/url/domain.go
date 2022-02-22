package url

type URL struct {
	ShortID string
	FullURL string
}

type URLRepository interface {
	GetURL(shortID string) (*URL, error)
	CreateURL(fullURL string) (*URL, error)
}

type URLService interface {
	FetchURL(shortID string) (*URL, error)
	BuildURL(fullURL string) (*URL, error)
}
