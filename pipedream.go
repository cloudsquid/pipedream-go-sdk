package pipedream

import (
	"net/http"
	"net/url"
	"sync"
)

// TODO: deduplicate this interface by creating a logger package in internal
type Logger interface {
	Debug(msg string, keyvals ...any)
	Info(msg string, keyvals ...any)
	Warn(msg string, keyvals ...any)
	Error(msg string, keyvals ...any)
}

type Client struct {
	logger         Logger
	apiKey         string
	httpClient     *http.Client
	environment    string
	projectID      string
	clientID       string
	clientSecret   string
	connectURL     *url.URL
	baseURL        *url.URL
	allowedOrigins []string

	token *Token
	mu    sync.Mutex
}

var (
	pipedreamApiURL     = "https://api.pipedream.com/v1/connect"
	pipedreamApiURLBase = "https://api.pipedream.com/v1/"
)

// NewClient creates a new client for the pipedream connect API
func NewClient(
	logger Logger,
	apiKey string,
	projectID string,
	environment string,
	clientID string,
	clientSecret string,
	allowedOrigins []string,
) *Client {
	connectParsed, err := url.Parse(pipedreamApiURL)
	if err != nil {
		logger.Error("parsing pipedream connect api url: %w", err)
	}

	baseParsed, err := url.Parse(pipedreamApiURLBase)
	if err != nil {
		logger.Error("parsing pipedream base api url: %w", err)
	}

	return &Client{
		logger:         logger,
		apiKey:         apiKey,
		projectID:      projectID,
		httpClient:     &http.Client{},
		environment:    environment,
		clientID:       clientID,
		clientSecret:   clientSecret,
		connectURL:     connectParsed,
		baseURL:        baseParsed,
		allowedOrigins: allowedOrigins,
	}
}
