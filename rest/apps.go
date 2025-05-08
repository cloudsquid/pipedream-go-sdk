package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/connect"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
	"io"
	"net/http"
	"net/url"
	"path"
)

type ListAppsResponse struct {
	PageInfo connect.PageInfo `json:"page_info,omitzero"`
	Data     []*connect.App   `json:"data,omitzero"`
}

type GetAppResponse struct {
	Data *connect.App `json:"data,omitzero"`
}

// Retrieve a list of all apps available on Pipedream
func (c *Client) ListApps(ctx context.Context, q string, hasComponents, hasActions, hasTriggers bool) (*ListAppsResponse, error) {
	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "apps")})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "q", q)
	if hasComponents {
		internal.AddQueryParams(queryParams, "has_components", "1")
	}
	if hasActions {
		internal.AddQueryParams(queryParams, "has_actions", "1")
	}
	if hasTriggers {
		internal.AddQueryParams(queryParams, "has_triggers", "1")
	}

	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get request for endpoint %s: %w", endpoint, err)
	}

	// can be done via Oauth too
	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	var appList ListAppsResponse
	err = json.NewDecoder(response.Body).Decode(&appList)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &appList, nil
}

// GetApp Retrieve metadata for a specific app
func (c *Client) GetApp(ctx context.Context, appID string) (*GetAppResponse, error) {
	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "apps", appID)})

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get an app request %s: %w", endpoint, err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get an app request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

	var app GetAppResponse
	err = json.NewDecoder(response.Body).Decode(&app)
	if err != nil {
		return nil, fmt.Errorf("decoding response for get an app request: %w", err)
	}

	return &app, nil
}
