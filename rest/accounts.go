package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/cloudsquid/pipedream-go-sdk/connect"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
)

type Account struct {
	ID              string              `json:"id,omitempty"`
	Name            string              `json:"name,omitempty"`
	ExternalID      string              `json:"external_id,omitempty"`
	Healthy         bool                `json:"healthy,omitempty"`
	Dead            bool                `json:"dead,omitempty"`
	App             connect.App         `json:"app,omitzero"`
	CreatedAt       time.Time           `json:"created_at,omitzero"`
	UpdatedAt       time.Time           `json:"updated_at,omitzero"`
	Credentials     connect.Credentials `json:"credentials,omitzero"`
	ExpiresAt       any                 `json:"expires_at,omitempty"`
	Error           any                 `json:"error,omitempty"`
	LastRefreshedAt time.Time           `json:"last_refreshed_at,omitzero"`
	NextRefreshAt   time.Time           `json:"next_refresh_at,omitzero"`
}

type ListAccountsResponse struct {
	PageInfo connect.PageInfo `json:"page_info"`
	Data     []Account        `json:"data"`
}

type GetAccountResponse struct {
	Data Account `json:"data"`
}

type ListAccountsOptions struct {
	IncludeCredentials bool
	App                *string
	OauthAppID         *string
}

// ListAccounts List connected accounts accessible by the authenticated user or workspace
func (c *Client) ListAccounts(
	ctx context.Context,
	opts *ListAccountsOptions,
) (*ListAccountsResponse, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts cannot be nil")
	}

	endpoint := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "accounts"),
	})

	queryParams := url.Values{}

	internal.AddQueryParams(queryParams, "include_credentials", fmt.Sprintf("%t", opts.IncludeCredentials))

	if opts.App != nil {
		internal.AddQueryParams(queryParams, "app", *opts.App)
	}
	if opts.OauthAppID != nil {
		internal.AddQueryParams(queryParams, "oauth_app_id", *opts.OauthAppID)
	}

	endpoint.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating list accounts request: %w", err)
	}

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing list accounts request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(raw))
	}

	var result ListAccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding list accounts response: %w", err)
	}

	return &result, nil
}

// GetAccount By default, this route returns metadata for a specific connected account
// Set include_credentials=true to return credentials that you can use in any app where you need the actual credentials
// (API key or OAuth access token for example)
func (c *Client) GetAccount(
	ctx context.Context,
	accountID string,
	includeCredentials bool,
) (*GetAccountResponse, error) {
	endpoint := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "accounts", accountID),
	})

	queryParams := url.Values{}

	if includeCredentials {
		internal.AddQueryParams(queryParams, "include_credentials", "true")
	}

	endpoint.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating get account request: %w", err)
	}

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get account request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(raw))
	}

	var result GetAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding get account response: %w", err)
	}

	return &result, nil
}
