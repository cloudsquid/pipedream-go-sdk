package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
	"net/http"
	"net/url"
	"path"
	"time"
)

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

// retrieve a short-lived token for that user
func (c *Client) AcquireUserToken(
	ctx context.Context,
	externalUserID string,
	webhookURI string, // optional, left empty won't be configured
) (*UserTokenResponse, error) {
	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, c.ProjectID(), "tokens")})

	endpoint := baseURL.String()

	request := UserTokenRequest{
		ExternalUserID: externalUserID,
		AllowedOrigins: c.AllowedOrigins(),
		WebhookURI:     webhookURI,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil,
			fmt.Errorf("couldn't marshalling user token request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil,
			fmt.Errorf("creating new request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil,
			fmt.Errorf("executing request to get user token: %w", err)
	}

	var userToken *UserTokenResponse
	if err := internal.UnmarshalResponse(response, &userToken); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response: %w", err)
	}

	return userToken, nil
}
