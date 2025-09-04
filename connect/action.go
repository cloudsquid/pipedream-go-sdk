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
	// ID is the ComponentID
	ID              string          `json:"id,omitempty"`
	ExternalUserID  string          `json:"external_user_id"`
	ConfiguredProps ConfiguredProps `json:"configured_props,omitempty"`

	DynamicPropsID string `json:"dynamic_props_id,omitempty"`
}

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
