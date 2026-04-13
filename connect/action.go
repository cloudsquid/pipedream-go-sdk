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

type InvokeActionRequest struct {
	ID              string          `json:"id,omitempty"`
	ExternalUserID  string          `json:"external_user_id"`
	ConfiguredProps ConfiguredProps `json:"configured_props,omitempty"`
	DynamicPropsID  string          `json:"dynamic_props_id,omitempty"`
	Version         string          `json:"version,omitempty"`
	StashID         any             `json:"stash_id,omitempty"`
}

type InvokeActionResponse struct {
	Exports map[string]any `json:"exports,omitempty"`
	OS      []any          `json:"os,omitempty"`
	Ret     any            `json:"ret,omitempty"`
	StashID string         `json:"stash_id,omitempty"`
}

type InvokeActionOptions struct {
	ComponentKey    string
	ExternalUserID  string
	ConfiguredProps ConfiguredProps
	DynamicPropsID  string
	Version         string
	StashID         any // string or bool
}

// Deprecated: Use InvokeActionWithOptions instead.
func (c *Client) InvokeAction(
	ctx context.Context,
	componentKey string,
	externalUserID string,
	props ConfiguredProps,
	dynamicPropsId string,
) (map[string]any, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "actions", "run")})

	invokeActionReq := InvokeActionRequest{
		ID:              componentKey,
		ConfiguredProps: props,
		ExternalUserID:  externalUserID,
	}

	if dynamicPropsId != "" {
		invokeActionReq.DynamicPropsID = dynamicPropsId
	}

	jsonBytes, err := json.MarshalIndent(invokeActionReq, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshalling invoke action request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL.String(),
		bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing invoke action request: %w", err)
	}

	var response map[string]any
	if err := internal.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling invoke action response: %w", err)
	}

	return response, nil
}

// InvokeActionWithOptions runs a component action with full options support.
func (c *Client) InvokeActionWithOptions(
	ctx context.Context,
	opts *InvokeActionOptions,
) (*InvokeActionResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "actions", "run"),
	})

	invokeActionReq := InvokeActionRequest{
		ID:              opts.ComponentKey,
		ConfiguredProps: opts.ConfiguredProps,
		ExternalUserID:  opts.ExternalUserID,
		DynamicPropsID:  opts.DynamicPropsID,
		Version:         opts.Version,
		StashID:         opts.StashID,
	}

	jsonBytes, err := json.Marshal(invokeActionReq)
	if err != nil {
		return nil, fmt.Errorf("marshalling invoke action request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL.String(),
		bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing invoke action request: %w", err)
	}

	var response InvokeActionResponse
	if err := internal.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling invoke action response: %w", err)
	}

	return &response, nil
}
