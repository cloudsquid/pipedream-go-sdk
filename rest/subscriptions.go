package rest

import (
	"context"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
	"io"
	"net/http"
	"net/url"
	"path"
)

// SubscribeToEmitter configures a source or workflow to receive events from any number of other workflows or sources
// For example, if you want a single workflow to run on 10 different RSS sources
// you can configure the workflow to listen for events from those 10 sources
func (c *Client) SubscribeToEmitter(
	ctx context.Context,
	emitterID,
	listenerID,
	eventName string, // optional
) error {
	c.Logger.Debug("creating subscription to emitter request")

	if emitterID == "" || listenerID == "" {
		return fmt.Errorf("emitter_id and listener_id are required")
	}

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "subscriptions"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "emitter_id", emitterID)
	internal.AddQueryParams(queryParams, "listener_id", listenerID)

	if eventName != "" {
		internal.AddQueryParams(queryParams, "event_name", eventName)
	}

	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating subscribe to emitter request: %w", err)
	}

	resp, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return fmt.Errorf("executing subscribe to emitter request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// AutoSubscribeToEvent automatically subscribes a listener to events from new workflows/sources
func (c *Client) AutoSubscribeToEvent(
	ctx context.Context,
	eventName string,
	listenerID string,
) error {
	c.Logger.Debug("creating auto-subscription to event request")

	if eventName == "" || listenerID == "" {
		return fmt.Errorf("event_name and listener_id are required")
	}

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "auto_subscriptions"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "event_name", eventName)
	internal.AddQueryParams(queryParams, "listener_id", listenerID)

	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating auto-subscription to event request: %w", err)
	}

	resp, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return fmt.Errorf("executing auto-subscription to event request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteSubscription deletes an existing subscription
// this endpoint accepts the same parameters as the POST /subscriptions endpoint for creating subscriptions.
func (c *Client) DeleteSubscription(
	ctx context.Context,
	emitterID,
	listenerID,
	eventName string,
) error {
	c.Logger.Debug("creating delete subscription request")

	if emitterID == "" || listenerID == "" {
		return fmt.Errorf("emitter_id and listener_id are required")
	}

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "subscriptions"),
	})

	queryParams := url.Values{}

	internal.AddQueryParams(queryParams, "emitter_id", emitterID)
	internal.AddQueryParams(queryParams, "listener_id", listenerID)
	internal.AddQueryParams(queryParams, "event_name", eventName)

	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating delete subscription request: %w", err)
	}

	resp, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return fmt.Errorf("executing delete subscription request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
