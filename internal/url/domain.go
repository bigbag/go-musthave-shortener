package url

type URL struct {
	ShortID  string
	FullURL  string
	ShortURL string
}

type ShortenRequest struct {
	FullURL string `json:"url"`
}

type UserURL struct {
	FullURL  string `json:"original_url"`
	ShortURL string `json:"short_url"`
}

type URLRepository interface {
	GetURL(shortID string) (*URL, error)
	CreateURL(fullURL string, userID string) (*URL, error)
	FindAllByUserID(userID string) ([]*URL, error)
	Status() error
	Close() error
}

type URLService interface {
	FetchURL(shortID string) (*URL, error)
	BuildURL(baseURL string, fullURL string, userID string) (*URL, error)
	FetchUserURLs(baseURL string, userID string) ([]*UserURL, error)
	Status() error
	Shutdown() error
}
