package connect

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type Account struct {
	ID              string      `json:"id,omitempty"`
	Name            string      `json:"name,omitempty"`
	ExternalID      string      `json:"external_id,omitempty"`
	Healthy         bool        `json:"healthy,omitempty"`
	Dead            bool        `json:"dead,omitempty"`
	App             App         `json:"app,omitzero"`
	CreatedAt       time.Time   `json:"created_at,omitzero"`
	UpdatedAt       time.Time   `json:"updated_at,omitzero"`
	Credentials     Credentials `json:"credentials,omitzero"`
	ExpiresAt       any         `json:"expires_at,omitempty"`
	Error           any         `json:"error,omitempty"`
	LastRefreshedAt time.Time   `json:"last_refreshed_at,omitzero"`
	NextRefreshAt   time.Time   `json:"next_refresh_at,omitzero"`
}

type App struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	NameSlug    string `json:"name_slug,omitempty"`
	AuthType    string `json:"auth_type,omitempty"`
	Description string `json:"description,omitempty"`
}

type Credentials struct {
	OauthClientId    string `json:"oauth_client_id,omitempty"`
	OauthAccessToken string `json:"oauth_access_token,omitempty"`
	OauthUid         string `json:"oauth_uid,omitempty"`
}

type ListAccountsResponse struct {
	PageInfo PageInfo   `json:"page_info"`
	Data     []*Account `json:"data"`
}

type GetAccountResponse struct {
	Data Account `json:"data"`
}

// ListAccounts lists all accounts related to the currently set projectID
// All the parameters are optional
func (c *Client) ListAccounts(
	ctx context.Context,
	externalUserID string,
	app string,
	oauthAppId string,
	includeCredentials bool,
) (*ListAccountsResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "accounts")})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	internal.AddQueryParams(queryParams, "app", app)
	internal.AddQueryParams(queryParams, "oauth_app_id", oauthAppId)
	internal.AddQueryParams(queryParams,
		"include_credentials",
		strconv.FormatBool(includeCredentials),
	)

	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get request for endpoint %s: %w", endpoint, err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing list account request: %w", err)
	}
	defer response.Body.Close()

	var accountsList ListAccountsResponse

	err = json.NewDecoder(response.Body).Decode(&accountsList)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &accountsList, nil
}

// GetAccount Retrieve the account details for a specific account based on the account ID
func (c *Client) GetAccount(
	ctx context.Context,
	externalUserID string,
	app string,
	includeCredentials bool,
	accountId string,
) (*GetAccountResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "accounts", accountId)})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	internal.AddQueryParams(queryParams, "app", app)
	internal.AddQueryParams(
		queryParams,
		"include_credentials",
		strconv.FormatBool(includeCredentials),
	)

	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get account request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request to get account: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil,
			fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(bodyBytes))
	}

	var accountDetail GetAccountResponse
	err = json.NewDecoder(response.Body).Decode(&accountDetail)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response for request to get account: %w", err)
	}

	return &accountDetail, nil
}

// DeleteAccount Delete a specific connected account for an end user, and any deployed triggers
func (c *Client) DeleteAccount(
	ctx context.Context,
	accountId string,
) error {
	endpoint := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "accounts", accountId)}).String()

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating delete account request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return fmt.Errorf("executing request to delete an account: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("expected status %d, got %d",
		http.StatusNoContent, response.StatusCode)
}

// DeleteAccounts Delete all connected accounts for a specific app
func (c *Client) DeleteAccounts(
	ctx context.Context,
	appID string,
) error {
	endpoint := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "apps", appID, "accounts")}).String()

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating get account request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return fmt.Errorf("executing request to get account: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("expected status %d, got %d", http.StatusNoContent, response.StatusCode)
}

// DeleteEndUser Delete an end user, all their connected accounts, and any deployed triggers.
func (c *Client) DeleteEndUser(
	ctx context.Context,
	externalUserID string,
) error {
	endpoint := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "users", externalUserID)}).String()

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating delete end user request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return fmt.Errorf("executing request to delete end usert: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("expected status %d, got %d", http.StatusNoContent, response.StatusCode)
}
