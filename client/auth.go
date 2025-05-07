package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	CreatedAt   int    `json:"created_at,omitempty"`
	ExpiresAt   time.Time
}

type UserTokenRequest struct {
	ExternalUserID     string   `json:"external_user_id"`
	AllowedOrigins     []string `json:"allowed_origins"`
	SuccessRedirectURI string   `json:"success_redirect_uri,omitempty"`
	ErrorRedirectURI   string   `json:"error_redirect_uri,omitempty"`
	WebhookURI         string   `json:"webhook_uri,omitempty"`
}

type UserTokenResponse struct {
	ConnectLinkURL string    `json:"connect_link_url,omitempty"`
	ExpiresAt      time.Time `json:"expires_at,omitzero"`
	Token          string    `json:"token,omitempty"`
}

func (c *Client) AcquireAccessToken() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != nil && time.Now().Before(c.token.ExpiresAt.Add(-1*time.Minute)) {
		return nil
	}

	endpoint := c.connectURL.ResolveReference(&url.URL{
		Path: path.Join(c.connectURL.Path, "oauth", "token"),
	}).String()

	type payload struct {
		GrantType    string `json:"grant_type,omitempty"`
		ClientID     string `json:"client_id,omitempty"`
		ClientSecret string `json:"client_secret,omitempty"`
	}

	bs, err := json.Marshal(&payload{
		GrantType:    "client_credentials",
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
	})
	if err != nil {
		return fmt.Errorf("couldn't marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("creating new request for endpoint %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request for new token: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected statuscode %s", response.Status)
	}

	var token Token
	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		return fmt.Errorf("decoding token: %w", err)
	}

	token.ExpiresAt = time.Now().Add(time.Second * time.Duration(token.ExpiresIn))
	c.token = &token

	return nil
}
