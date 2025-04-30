package pipedream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type ListAppsResponse struct {
	PageInfo PageInfo `json:"page_info,omitzero"`
	Data     []*App   `json:"data,omitzero"`
}

// https://pipedream.com/docs/rest-api/#list-apps
func (p *Client) ListApps(ctx context.Context, q string) (*ListAppsResponse, error) {
	p.logger.Info("Listing apps")

	baseURL := p.baseURL.ResolveReference(&url.URL{
		Path: path.Join(p.baseURL.Path, "apps")})

	queryParams := url.Values{}
	addQueryParams(queryParams, "q", q)

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
