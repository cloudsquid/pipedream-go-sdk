package client

import (
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Client struct {
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
		log.Fatal("parsing pipedream connect api url: %w", err)
	}

	restParsed, err := url.Parse(restURL)
	if err != nil {
		log.Fatal("parsing pipedream connect api url: %w", err)
	}

	return &Client{
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
