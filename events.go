package pipedream

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type GetSourceEventsResponse struct {
	PageInfo PageInfo      `json:"page_info"`
	Data     []SourceEvent `json:"data"`
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
func (p *Client) GetSourceEvents(
	ctx context.Context,
	sourceID string,
	limit int,
	expand bool,
) (*GetSourceEventsResponse, error) {
	p.logger.Info("get source events", "sourceID", sourceID)

	baseURL := p.baseURL.ResolveReference(&url.URL{
		Path: path.Join(p.baseURL.Path, "sources", sourceID, "event_summaries")})

	queryParams := url.Values{}

	if limit > 0 {
		addQueryParams(queryParams, "limit", strconv.Itoa(limit))
	}

	if expand {
		addQueryParams(queryParams, "expand", "event")
	}

	baseURL.RawQuery = queryParams.Encode()

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get source events request %s: %w", endpoint, err)
	}

	response, err := p.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get source events request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var respJson GetSourceEventsResponse
	if err := unmarshalResponse(response, &respJson); err != nil {
		return nil, fmt.Errorf(
			"unmarshalling get source events response: %w", err)
	}

	return &respJson, nil
}
