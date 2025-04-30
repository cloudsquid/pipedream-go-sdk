package pipedream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type ListAppsResponse struct {
	PageInfo PageInfo `json:"page_info,omitzero"`
	Data     []*App   `json:"data,omitzero"`
}

type GetAppResponse struct {
	Data *App `json:"data,omitzero"`
}

// Retrieve a list of all apps available on Pipedream
func (p *Client) ListApps(ctx context.Context, q string, hasComponents, hasActions, hasTriggers bool) (*ListAppsResponse, error) {
	p.logger.Info("Listing apps")

	baseURL := p.baseURL.ResolveReference(&url.URL{
		Path: path.Join(p.baseURL.Path, "apps")})

	queryParams := url.Values{}
	addQueryParams(queryParams, "q", q)
	if hasComponents {
		addQueryParams(queryParams, "has_components", "1")
	}
	if hasActions {
		addQueryParams(queryParams, "has_actions", "1")
	}
	if hasTriggers {
		addQueryParams(queryParams, "has_triggers", "1")
	}

	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get request for endpoint %s: %w", endpoint, err)
	}

	response, err := p.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	p.logger.Info("Response status code", "code", response.StatusCode)

	var appList ListAppsResponse
	err = json.NewDecoder(response.Body).Decode(&appList)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &appList, nil
}

// GetApp Retrieve metadata for a specific app
func (p *Client) GetApp(ctx context.Context, appID string) (*GetAppResponse, error) {
	p.logger.Info("get an apps")

	baseURL := p.baseURL.ResolveReference(&url.URL{
		Path: path.Join(p.baseURL.Path, "apps", appID)})

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get request for endpoint %s: %w", endpoint, err)
	}

	response, err := p.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

	var app GetAppResponse
	err = json.NewDecoder(response.Body).Decode(&app)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &app, nil
}
