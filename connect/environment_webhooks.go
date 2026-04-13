package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/cloudsquid/pipedream-go-sdk/internal"
)

type EnvironmentWebhook struct {
	ID            string `json:"id,omitempty"`
	URL           string `json:"url,omitempty"`
	SigningKey    string `json:"signing_key,omitempty"`
	SigningKeySet bool   `json:"signing_key_set,omitempty"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
}

type GetEnvironmentWebhookResponse struct {
	Data *EnvironmentWebhook `json:"data,omitempty"`
}

type SetEnvironmentWebhookResponse struct {
	Data *EnvironmentWebhook `json:"data,omitempty"`
}

type SetEnvironmentWebhookRequest struct {
	URL string `json:"url"`
}

type RegenerateEnvironmentWebhookSigningKeyResponse struct {
	Data *EnvironmentWebhook `json:"data,omitempty"`
}

// GetEnvironmentWebhook retrieves the project environment webhook.
func (c *Client) GetEnvironmentWebhook(
	ctx context.Context,
) (*GetEnvironmentWebhookResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "webhook"),
	})

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating get environment webhook request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get environment webhook request: %w", err)
	}

	var result GetEnvironmentWebhookResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get environment webhook response: %w", err)
	}

	return &result, nil
}

// SetEnvironmentWebhook creates or updates the project environment webhook.
func (c *Client) SetEnvironmentWebhook(
	ctx context.Context,
	webhookURL string,
) (*SetEnvironmentWebhookResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "webhook"),
	})

	body, err := json.Marshal(SetEnvironmentWebhookRequest{URL: webhookURL})
	if err != nil {
		return nil, fmt.Errorf("marshalling set environment webhook request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, baseURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating set environment webhook request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing set environment webhook request: %w", err)
	}

	var result SetEnvironmentWebhookResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling set environment webhook response: %w", err)
	}

	return &result, nil
}

// DeleteEnvironmentWebhook deletes the project environment webhook.
func (c *Client) DeleteEnvironmentWebhook(
	ctx context.Context,
) error {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "webhook"),
	})

	req, err := http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating delete environment webhook request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return fmt.Errorf("executing delete environment webhook request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("expected status %d, got %d", http.StatusNoContent, response.StatusCode)
}

// RegenerateEnvironmentWebhookSigningKey regenerates the signing key for the project environment webhook.
func (c *Client) RegenerateEnvironmentWebhookSigningKey(
	ctx context.Context,
) (*RegenerateEnvironmentWebhookSigningKeyResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "webhook", "regenerate_signing_key"),
	})

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating regenerate webhook signing key request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing regenerate webhook signing key request: %w", err)
	}

	var result RegenerateEnvironmentWebhookSigningKeyResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling regenerate webhook signing key response: %w", err)
	}

	return &result, nil
}
