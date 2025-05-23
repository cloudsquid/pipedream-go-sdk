package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
	"net/http"
	"net/url"
	"path"
)

type CreateSourceRequest struct {
	ComponentID   string `json:"component_id,omitempty"`
	ComponentCode string `json:"component_code,omitempty"`
	ComponentURL  string `json:"component_url,omitempty"`
	Name          string `json:"name,omitempty"`
}

type UpdateSourceRequest struct {
	ComponentID   string `json:"component_id,omitempty"`
	ComponentCode string `json:"component_code,omitempty"`
	ComponentURL  string `json:"component_url,omitempty"`
	Name          string `json:"name,omitempty"`
	Active        bool   `json:"active,omitempty"` // default is true
}

// CreateSourceResponse is the response for both creating and updating a source
type CreateSourceResponse struct {
	Data SourceData `json:"data"`
}

type SourceData struct {
	ID              string      `json:"id"`
	UserID          string      `json:"user_id"`
	ComponentID     string      `json:"component_id"`
	ConfiguredProps SourceProps `json:"configured_props"`
	Active          bool        `json:"active"`
	CreatedAt       int64       `json:"created_at"`
	UpdatedAt       int64       `json:"updated_at"`
	Name            string      `json:"name"`
	NameSlug        string      `json:"name_slug"`
}

type SourceProps struct {
	URL   string        `json:"url"`
	Timer TimerSchedule `json:"timer"`
}

type TimerSchedule struct {
	Cron            *string `json:"cron"` // nullable
	IntervalSeconds int     `json:"interval_seconds"`
}

// CreateSource Event run code to collect events from an API, or receive events via webhooks, emitting those events for use on Pipedream
// Event sources can function as workflow triggers
func (c *Client) CreateSource(
	ctx context.Context,
	componentID,
	componentCode,
	componentURL,
	name string,
) (*CreateSourceResponse, error) {
	if componentID == "" && componentCode == "" && componentURL == "" {
		return nil, fmt.Errorf("one of component_id, component_code, or component_url is required")
	}
	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "sources")})
	endpoint := baseURL.String()

	body := &CreateSourceRequest{
		ComponentID:   componentID,
		ComponentCode: componentCode,
		ComponentURL:  componentURL,
		Name:          name,
	}

	rb, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling creating source request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(rb))
	if err != nil {
		return nil, fmt.Errorf("creating creating source for endpoint %s: %w", endpoint, err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing creating source request: %w", err)
	}
	defer response.Body.Close()

	var source CreateSourceResponse
	if err := internal.UnmarshalResponse(response, &source); err != nil {
		return nil, fmt.Errorf(
			"parsing response for creating source request: %w", err)
	}

	return &source, nil
}

// UpdateSource updates a source
func (c *Client) UpdateSource(
	ctx context.Context,
	sourceID,
	componentID,
	componentCode,
	componentURL,
	name string,
	active bool,
) (*CreateSourceResponse, error) {
	if componentID == "" && componentCode == "" && componentURL == "" {
		return nil, fmt.Errorf("one of component_id, component_code, or component_url is required")
	}
	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "sources", sourceID)})
	endpoint := baseURL.String()

	body := &UpdateSourceRequest{
		ComponentID:   componentID,
		ComponentCode: componentCode,
		ComponentURL:  componentURL,
		Name:          name,
		Active:        active,
	}

	rb, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling updating source request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(rb))
	if err != nil {
		return nil, fmt.Errorf("creating updating source for endpoint %s: %w", endpoint, err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing updating source request: %w", err)
	}
	defer response.Body.Close()

	var source CreateSourceResponse
	if err := internal.UnmarshalResponse(response, &source); err != nil {
		return nil, fmt.Errorf(
			"parsing response for updating source request: %w", err)
	}

	return &source, nil
}

// DeleteSource deletes a source
func (c *Client) DeleteSource(
	ctx context.Context,
	sourceID string,
) error {
	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "sources", sourceID)})
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating deleting source for endpoint %s: %w", endpoint, err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return fmt.Errorf("executing deleting source request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("expected status %d, got %d",
		http.StatusNoContent, response.StatusCode)
}
