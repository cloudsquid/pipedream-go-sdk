package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/cloudsquid/pipedream-go-sdk/internal"
)

type UserTokenRequest struct {
	ExternalUserID     string   `json:"external_user_id"`
	AllowedOrigins     []string `json:"allowed_origins,omitempty"`
	SuccessRedirectURI string   `json:"success_redirect_uri,omitempty"`
	ErrorRedirectURI   string   `json:"error_redirect_uri,omitempty"`
	WebhookURI         string   `json:"webhook_uri,omitempty"`
	ExpiresIn          *int     `json:"expires_in,omitempty"`
	Scope              string   `json:"scope,omitempty"`
}

type AcquireUserTokenOptions struct {
	ExternalUserID     string
	WebhookURI         string
	SuccessRedirectURI string
	ErrorRedirectURI   string
	ExpiresIn          *int
	Scope              string
}

type UserTokenResponse struct {
	ConnectLinkURL string    `json:"connect_link_url,omitempty"`
	ExpiresAt      time.Time `json:"expires_at,omitzero"`
	Token          string    `json:"token,omitempty"`
}

// Deprecated: Use AcquireUserTokenWithOptions instead.
func (c *Client) AcquireUserToken(
	ctx context.Context,
	externalUserID string,
	webhookURI string,
) (*UserTokenResponse, error) {
	return c.AcquireUserTokenWithOptions(ctx, &AcquireUserTokenOptions{
		ExternalUserID: externalUserID,
		WebhookURI:     webhookURI,
	})
}

// AcquireUserTokenWithOptions retrieves a short-lived connect token for an end user.
func (c *Client) AcquireUserTokenWithOptions(
	ctx context.Context,
	opts *AcquireUserTokenOptions,
) (*UserTokenResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "tokens"),
	})

	endpoint := baseURL.String()

	request := UserTokenRequest{
		ExternalUserID:     opts.ExternalUserID,
		AllowedOrigins:     c.AllowedOrigins(),
		WebhookURI:         opts.WebhookURI,
		SuccessRedirectURI: opts.SuccessRedirectURI,
		ErrorRedirectURI:   opts.ErrorRedirectURI,
		ExpiresIn:          opts.ExpiresIn,
		Scope:              opts.Scope,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshalling user token request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request to get user token: %w", err)
	}

	var userToken *UserTokenResponse
	if err := internal.UnmarshalResponse(response, &userToken); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response: %w", err)
	}

	return userToken, nil
}
