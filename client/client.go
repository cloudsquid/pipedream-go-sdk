package client

import (
	"net/http"
	"net/url"
	"sync"
)

type Logger interface {
	Debug(msg string, keyvals ...any)
	Info(msg string, keyvals ...any)
	Warn(msg string, keyvals ...any)
	Error(msg string, keyvals ...any)
}

type Client struct {
	Logger         Logger
	apiKey         string
	httpClient     *http.Client
	environment    string
	projectID      string
	clientID       string
	clientSecret   string
	connectURL     *url.URL
	restURL        *url.URL
	allowedOrigins []string

	token *Token
	mu    sync.Mutex
}

var (
	pipedreamApiURL     = "https://api.pipedream.com/v1/connect"
	pipedreamApiURLBase = "https://api.pipedream.com/v1/"
)

func NewClient(
	logger Logger,
	apiKey string,
	projectID string,
	environment string,
	clientID string,
	clientSecret string,
	allowedOrigins []string,
	connectURL string,
	restURL string,
) *Client {
	if connectURL == "" {
		connectURL = pipedreamApiURL
	}
	if restURL == "" {
		restURL = pipedreamApiURLBase
	}
	connectParsed, err := url.Parse(connectURL)
	if err != nil {
		logger.Error("parsing pipedream connect api url: %w", err)
	}

	restParsed, err := url.Parse(restURL)
	if err != nil {
		logger.Error("parsing pipedream base api url: %w", err)
	}

	return &Client{
		Logger:         logger,
		apiKey:         apiKey,
		projectID:      projectID,
		httpClient:     &http.Client{},
		environment:    environment,
		clientID:       clientID,
		clientSecret:   clientSecret,
		connectURL:     connectParsed,
		restURL:        restParsed,
		allowedOrigins: allowedOrigins,
	}
}

func (c *Client) APIKey() string {
	return c.apiKey
}

func (c *Client) ProjectID() string {
	return c.projectID
}

func (c *Client) Environment() string {
	return c.environment
}

func (c *Client) ClientID() string {
	return c.clientID
}

func (c *Client) ClientSecret() string {
	return c.clientSecret
}

func (c *Client) AllowedOrigins() []string {
	return c.allowedOrigins
}

func (c *Client) ConnectURL() *url.URL {
	return c.connectURL
}

func (c *Client) RestURL() *url.URL {
	return c.restURL
}

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) Token() *Token {
	return c.token
}
