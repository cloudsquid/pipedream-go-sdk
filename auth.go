package pipedream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
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

// https://pipedream.com/docs/connect/api/#create-a-new-token
func (p *Client) AcquireUserToken(
	ctx context.Context,
	externalUserID string,
	webhookURI string, // optional, left empty won't be configured
) (*UserTokenResponse, error) {
	log.Infof("Acquiring user token and registering webhook uri %s", webhookURI)

	baseURL, err := url.Parse(pipedreamApiURL)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as url: %w", pipedreamApiURL, err)
	}
	baseURL.Path = path.Join(baseURL.Path, p.projectID, "tokens")

	request := UserTokenRequest{
		ExternalUserID: externalUserID,
		AllowedOrigins: p.allowedOrigins,
		WebhookURI:     webhookURI,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil,
			fmt.Errorf("couldn't marshalling user token request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil,
			fmt.Errorf("creating new request: %w", err)
	}

	response, err := p.doRequest(ctx, req)
	if err != nil {
		return nil,
			fmt.Errorf("executing request to get user token: %w", err)
	}

	var userToken *UserTokenResponse
	if err := unmarshalResponse(response, &userToken); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response: %w", err)
	}

	return userToken, nil
}

func (p *Client) acquireAccessToken() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != nil && time.Now().Before(p.token.ExpiresAt.Add(-1*time.Minute)) {
		return nil
	}

	endpoint := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join("oauth", "token"),
	}).String()

	type payload struct {
		GrantType    string `json:"grant_type,omitempty"`
		ClientID     string `json:"client_id,omitempty"`
		ClientSecret string `json:"client_secret,omitempty"`
	}

	bs, err := json.Marshal(&payload{
		GrantType:    "client_credentials",
		ClientID:     p.clientID,
		ClientSecret: p.clientSecret,
	})
	if err != nil {
		return fmt.Errorf("couldn't marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("creating new request for endpoint %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := p.httpClient.Do(req)
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
	p.token = &token

	p.logger.Info(token.AccessToken)

	return nil
}
