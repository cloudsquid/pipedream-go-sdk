package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type RevokeAccessTokenRequest struct {
	Token        string `json:"token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// RevokeAccessToken revokes an OAuth access token.
func (c *Client) RevokeAccessToken(
	ctx context.Context,
	token string,
	clientID string,
	clientSecret string,
) error {
	endpoint := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "oauth", "revoke"),
	}).String()

	payload := &RevokeAccessTokenRequest{
		Token:        token,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling revoke token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating revoke token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("executing revoke token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(raw))
	}

	return nil
}
