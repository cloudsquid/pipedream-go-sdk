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
	"strconv"
)

type GetSourceEventsResponse struct {
	PageInfo connect.PageInfo `json:"page_info"`
	Data     []SourceEvent    `json:"data"`
}

type SourceEvent struct {
	ID          string                 `json:"id"`
	IndexedAtMs int64                  `json:"indexed_at_ms"`
	Event       map[string]interface{} `json:"event"`
	Metadata    EventMetadata          `json:"metadata"`
}

type EventMetadata struct {
	EmitterID string `json:"emitter_id"`
	EmitID    string `json:"emit_id"`
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	ID        string `json:"id"`
	TS        int64  `json:"ts"`
}

// GetSourceEvents retrieves up to the last 100 events emitted by a source
func (c *Client) GetSourceEvents(
	ctx context.Context,
	sourceID string,
	limit int,
	expand bool,
) (*GetSourceEventsResponse, error) {
	c.Logger.Info("get source events", "sourceID", sourceID)

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "sources", sourceID, "event_summaries")})

	queryParams := url.Values{}

	if limit > 0 {
		internal.AddQueryParams(queryParams, "limit", strconv.Itoa(limit))
	}

	if expand {
		internal.AddQueryParams(queryParams, "expand", "event")
	}

	baseURL.RawQuery = queryParams.Encode()

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get source events request %s: %w", endpoint, err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get source events request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var respJson GetSourceEventsResponse
	if err := internal.UnmarshalResponse(response, &respJson); err != nil {
		return nil, fmt.Errorf(
			"unmarshalling get source events response: %w", err)
	}

	return &respJson, nil
}

// DeleteSourceEvents deletes events for a source starting from startID (inclusive).
func (c *Client) DeleteSourceEvents(
	ctx context.Context,
	sourceID,
	startID,
	endID string, // optional
) error {
	c.Logger.Info("deleting source events")

	if sourceID == "" || startID == "" {
		return fmt.Errorf("both sourceID and startID are required")
	}

	endpointURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "sources", sourceID, "events"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "start_id", startID)

	if endID != "" {
		internal.AddQueryParams(queryParams, "end_id", endID)
	}

	endpointURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodDelete, endpointURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating delete source events request: %w", err)
	}

	resp, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return fmt.Errorf("executing delete source events request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
