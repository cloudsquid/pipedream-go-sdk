package connect

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/cloudsquid/pipedream-go-sdk/internal"
)

type UsageRecord struct {
	CreditsUsed          int   `json:"credits_used"`
	ActionRunCreditsUsed *int  `json:"action_run_credits_used,omitempty"`
	ProxyCreditsUsed     *int  `json:"proxy_credits_used,omitempty"`
	SourceEmitCreditsUsed *int `json:"source_emit_credits_used,omitempty"`
	UsageStartTS         int64 `json:"usage_start_ts"`
	UsageEndTS           int64 `json:"usage_end_ts"`
}

type ListUsageResponse struct {
	Data []UsageRecord `json:"data,omitempty"`
}

// ListUsage retrieves Connect usage records for a given time range.
// startTS and endTS are Unix timestamps in seconds.
func (c *Client) ListUsage(
	ctx context.Context,
	startTS int64,
	endTS int64,
) (*ListUsageResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "usage"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "start_ts", strconv.FormatInt(startTS, 10))
	internal.AddQueryParams(queryParams, "end_ts", strconv.FormatInt(endTS, 10))
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating list usage request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing list usage request: %w", err)
	}

	var result ListUsageResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling list usage response: %w", err)
	}

	return &result, nil
}
