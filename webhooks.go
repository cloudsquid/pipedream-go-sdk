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

type CreateWebhookResponse struct {
	Data Webhook `json:"data"`
}

type Webhook struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	URL         string  `json:"url"`
	Active      bool    `json:"active"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

// CreateWebhook Creates a webhook pointing to a URL
// Configure a subscription to deliver events to this webhook
func (c *Client) CreateWebhook(
	ctx context.Context,
	endpoint,
	name,
	description string,
) (*CreateWebhookResponse, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("url is required")
	}

	baseURL := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "webhooks"),
	})

	queryParams := url.Values{}

	addQueryParams(queryParams, "url", endpoint)
	if name != "" {
		addQueryParams(queryParams, "name", name)
	}
	if description != "" {
		addQueryParams(queryParams, "description", description)
	}
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating webhook request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing create webhook  request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(body))
	}

	var webhook CreateWebhookResponse
	if err := json.NewDecoder(response.Body).Decode(&webhook); err != nil {
		return nil, fmt.Errorf("decoding create webhook response: %w", err)
	}

	return &webhook, nil
}

// DeleteWebhook deletes a webhook in your account.
func (c *Client) DeleteWebhook(
	ctx context.Context,
	id string,
) error {
	baseURL := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "webhooks", id),
	})

	req, err := http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating delete webhook request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return fmt.Errorf("executing delete webhook  request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil
	} else {
		return fmt.Errorf("expected status %d, got %d",
			http.StatusNoContent, response.StatusCode)
	}
}
