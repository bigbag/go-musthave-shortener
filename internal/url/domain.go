package url

type NotUniqueURLError struct{}

func (e *NotUniqueURLError) Error() string {
	return "not unique url"
}

type URL struct {
	ShortID       string
	FullURL       string
	ShortURL      string
	CorrelationID string
	Removed       bool
}

type JSONRequest struct {
	FullURL string `json:"url"`
}

type BatchRequestItem struct {
	FullURL       string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

type BatchRequest []BatchRequestItem

type BatchResponseItem struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

type BatchResponse []*BatchResponseItem

type UserURL struct {
	FullURL  string `json:"original_url"`
	ShortURL string `json:"short_url"`
}

type URLRepository interface {
	GetURL(shortID string) (*URL, error)
	FindAllByUserID(userID string) ([]*URL, error)
	CreateURL(fullURL string, userID string) (string, error)
	CreateBatchOfURL(items BatchRequest, userID string) ([]*URL, error)
	DeleteUserURLs(userID string, shortIDs []string) error
	Status() error
	Close() error
}

type URLService interface {
	FetchURL(shortID string) (*URL, error)
	FetchUserURLs(baseURL string, userID string) ([]*UserURL, error)
	BuildURL(baseURL string, fullURL string, userID string) (string, error)
	BuildBatchOfURL(
		baseURL string,
		items BatchRequest,
		userID string,
	) (BatchResponse, error)
	DeleteUserURLs(userID string, shortIDs []string) error
	Status() error
	Shutdown() error
}
