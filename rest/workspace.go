package rest

import (
	"context"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/connect"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
	"io"
	"net/http"
	"net/url"
	"path"
)

type GetWorkspaceResponse struct {
	Data Workspace `json:"data"`
}

type Workspace struct {
	ID                string `json:"id"`
	OrgName           string `json:"orgname"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	DailyCreditsQuota int    `json:"daily_credits_quota"`
	DailyCreditsUsed  int    `json:"daily_credits_used"`
}

type GetWorkspaceConnectedAccountsResponse struct {
	PageInfo connect.PageInfo   `json:"page_info"`
	Data     []ConnectedAccount `json:"data"`
}

type ConnectedAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetWorkspaceSubscriptionsResponse struct {
	Data []Subscription `json:"data"`
}

type Subscription struct {
	ID         string `json:"id"`
	EmitterID  string `json:"emitter_id"`
	ListenerID string `json:"listener_id"`
	EventID    string `json:"event_id"`
}

type GetWorkspaceSourcesResponse struct {
	PageInfo connect.PageInfo `json:"page_info"`
	Data     []Source         `json:"data"`
}

type Source struct {
	ID              string                  `json:"id"`
	ComponentID     string                  `json:"component_id"`
	ConfiguredProps connect.ConfiguredProps `json:"configured_props"`
	Active          bool                    `json:"active"`
	CreatedAt       int64                   `json:"created_at"`
	UpdatedAt       int64                   `json:"updated_at"`
	Name            string                  `json:"name"`
	NameSlug        string                  `json:"name_slug"`
}

type Data struct {
	ID              string                  `json:"id"`
	ComponentID     string                  `json:"component_id"`
	ConfiguredProps connect.ConfiguredProps `json:"configured_props"`
	Active          bool                    `json:"active"`
	CreatedAt       int64                   `json:"created_at"`
	UpdatedAt       int64                   `json:"updated_at"`
	Name            string                  `json:"name"`
	NameSlug        string                  `json:"name_slug"`
}

// GetWorkspace views your workspace’s current credit usage for the billing period in real time
func (c *Client) GetWorkspace(
	ctx context.Context,
	orgID string,
) (*GetWorkspaceResponse, error) {
	c.Logger.Debug("get workspace")

	if orgID == "" {
		return nil, fmt.Errorf("orgID is required")
	}

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "workspaces", orgID),
	})

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get workspace request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get workspace request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var result GetWorkspaceResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get workspace response:e: %w", err)
	}

	return &result, nil
}

// GetWorkspaceConnectedAccounts Retrieves all the connected accounts for a specific workspace
func (c *Client) GetWorkspaceConnectedAccounts(
	ctx context.Context,
	orgID string,
	query string, // optional
) (*GetWorkspaceConnectedAccountsResponse, error) {
	c.Logger.Debug("get workspace connected accounts")

	if orgID == "" {
		return nil, fmt.Errorf("orgID is required")
	}

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "workspaces", orgID, "accounts"),
	})

	queryParams := url.Values{}

	internal.AddQueryParams(queryParams, "query", query)

	baseURL.RawQuery = queryParams.Encode()

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get workspace accounts request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get workspace accounts request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var result GetWorkspaceConnectedAccountsResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get workspace accounts response:e: %w", err)
	}

	return &result, nil
}

// GetWorkspaceSubscriptions Retrieves all the subscriptions configured for a specific workspace
func (c *Client) GetWorkspaceSubscriptions(
	ctx context.Context,
	orgID string,
) (*GetWorkspaceSubscriptionsResponse, error) {
	c.Logger.Debug("get workspace subscriptions")

	if orgID == "" {
		return nil, fmt.Errorf("orgID is required")
	}

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "workspaces", orgID, "subscriptions"),
	})

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get workspace subscriptions request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get workspace subscriptions request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var result GetWorkspaceSubscriptionsResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get workspace subscriptions response:e: %w", err)
	}

	return &result, nil
}

// GetWorkspaceSources Retrieves all the event sources configured for a specific workspace
func (c *Client) GetWorkspaceSources(
	ctx context.Context,
	orgID string,
) (*GetWorkspaceSourcesResponse, error) {
	c.Logger.Debug("get workspace sources")

	if orgID == "" {
		return nil, fmt.Errorf("orgID is required")
	}

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "workspaces", orgID, "sources"),
	})

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get workspace sources request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get workspace sources request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var result GetWorkspaceSourcesResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get workspace sources response: %w", err)
	}

	return &result, nil
}
