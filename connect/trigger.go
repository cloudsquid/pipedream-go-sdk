package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/cloudsquid/pipedream-go-sdk/internal"
)

type Trigger struct {
	ID                   string             `json:"id"`
	OwnerID              string             `json:"owner_id"`
	ComponentID          string             `json:"component_id"`
	ComponentKey         string             `json:"component_key,omitempty"`
	ConfigurableProps    []ConfigurableProp `json:"configurable_props,omitempty"`
	ConfiguredProps      ConfiguredProps    `json:"configured_props,omitempty"`
	Active               bool               `json:"active,omitempty"`
	CreatedAt            int                `json:"created_at"`
	UpdatedAt            int                `json:"updated_at"`
	Name                 string             `json:"name"`
	NameSlug             string             `json:"name_slug"`
	EmailAddress         string             `json:"email_address,omitempty"`
	EmitOnDeploy         *bool              `json:"emit_on_deploy,omitempty"`
	WebhookSigningKey    string             `json:"webhook_signing_key,omitempty"`
	Type                 string             `json:"type,omitempty"`
	CallbackObservations any                `json:"callback_observations,omitempty"`
}

type GetTriggerResponse struct {
	Data Trigger `json:"data"`
}

type DeployTriggerResponse struct {
	Data Trigger `json:"data"`
}

type PageInfo struct {
	TotalCount  int    `json:"total_count,omitempty"`
	Count       int    `json:"count,omitempty"`
	StartCursor string `json:"start_cursor,omitempty"`
	EndCursor   string `json:"end_cursor,omitempty"`
}

type TriggerList struct {
	PageInfo PageInfo  `json:"page_info,omitzero"`
	Data     []Trigger `json:"data,omitempty"`
}

func (t Trigger) String() string {
	return fmt.Sprintf("\tName: %s\t String: %s", t.Name, t.ComponentID)
}

type TriggerEvent struct {
	E  Event  `json:"e,omitzero"`
	K  string `json:"k,omitempty"`
	TS int    `json:"ts,omitempty"`
	ID string `json:"id,omitempty"`
}

type Event struct {
	Method   string            `json:"method,omitempty"`
	Path     string            `json:"path,omitempty"`
	Query    []string          `json:"query,omitempty"`
	ClientIP string            `json:"client_ip,omitempty"`
	URL      string            `json:"url,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

type TriggerEventList struct {
	Data []TriggerEvent `json:"data,omitempty"`
}

type DeployTriggerRequest struct {
	ComponentKey    string          `json:"id"`
	ConfiguredProps ConfiguredProps `json:"configured_props"`
	ExternalUserID  string          `json:"external_user_id"`
	WebhookURL      string          `json:"webhook_url,omitempty"`
	WorkflowID      string          `json:"workflow_id,omitempty"`
	DynamicPropsID  string          `json:"dynamic_props_id,omitempty"`
	Version         string          `json:"version,omitempty"`
	EmitOnDeploy    *bool           `json:"emit_on_deploy,omitempty"`
}

type DeployTriggerOptions struct {
	ComponentKey    string
	ExternalUserID  string
	ConfiguredProps ConfiguredProps
	WebhookURL      string
	DynamicPropsID  string
	WorkflowID      string
	Version         string
	EmitOnDeploy    *bool
}

type ListDeployedTriggersOptions struct {
	ExternalUserID string
	EmitterType    *string
	Limit          *int
	After          *string
	Before         *string
}

type DeleteDeployedTriggerOptions struct {
	DeployedTriggerID string
	ExternalUserID    string
	IgnoreHookErrors  *bool
}


type UpdateTriggerWebhooksRequest struct {
	WebhookURLs []string `json:"webhook_urls,omitempty"`
}

type UpdateTriggerWorkflowsRequest struct {
	WorkflowIDs []string `json:"workflow_ids,omitempty"`
}

type TriggerWebhook struct {
	ID            string `json:"id,omitempty"`
	URL           string `json:"url,omitempty"`
	SigningKey    string `json:"signing_key,omitempty"`
	SigningKeySet bool   `json:"signing_key_set,omitempty"`
}

type TriggerWebhookURLs struct {
	WebhookURLs []string         `json:"webhook_urls,omitempty"`
	Webhooks    []TriggerWebhook `json:"webhooks,omitempty"`
}

type TriggerWorkflowIDs struct {
	WorkflowIDs []string `json:"workflow_ids,omitempty"`
}

func (c *Client) DeployTrigger(
	ctx context.Context,
	componentKey string,
	externalUserID string,
	configuredProps ConfiguredProps,
	webhookURL string,
	dynamicPropsID string, // OPTIONAL
	workflowID string, // OPTIONAL
) (*Trigger, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "triggers", "deploy"),
	})

	trigger := DeployTriggerRequest{
		ComponentKey:    componentKey,
		ConfiguredProps: configuredProps,
		WebhookURL:      webhookURL,
		WorkflowID:      workflowID,
		DynamicPropsID:  dynamicPropsID,
		ExternalUserID:  externalUserID,
	}

	jsonBytes, err := json.Marshal(trigger)
	if err != nil {
		return nil, fmt.Errorf("marshalling deploy trigger request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL.String(),
		bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("creating new deploy trigger request: %w", err)
	}

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing deploy trigger request: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body for deploying trigger: %w", err)
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("failed to deploy trigger to: %s", trigger.WebhookURL)
	}

	var response DeployTriggerResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("unmarhalling response for trigger response: %w: %w",
			errors.New(string(bodyBytes)),
			err)
	}

	return &response.Data, nil
}

func (c *Client) ListDeployedTriggers(
	ctx context.Context,
	externalUserID string,
) (*TriggerList, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating new request to list deployed triggers: %w", err)
	}

	listResponse, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request to list deployed triggers: %w", err)
	}

	var list TriggerList
	if err := internal.UnmarshalResponse(listResponse, &list); err != nil {
		return nil, fmt.Errorf("unmarshalling response to list: %w", err)
	}

	return &list, nil
}

func (c *Client) GetDeployedTrigger(
	ctx context.Context,
	deployedComponentID string,
	externalUserId string,
) (*Trigger, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", deployedComponentID),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserId)
	baseURL.RawQuery = queryParams.Encode()

	getRequest, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, getRequest)
	if err != nil {
		return nil, fmt.Errorf("executing request to retrieve trigger: %w", err)
	}

	var trigger GetTriggerResponse
	if err := internal.UnmarshalResponse(response, &trigger); err != nil {
		return nil, fmt.Errorf("unmarshalling response to retrieve trigger: %w", err)
	}

	return &trigger.Data, nil
}

func (c *Client) DeleteDeployedTrigger(
	ctx context.Context,
	deployedTriggerID string,
	externalUserID string,
) error {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", deployedTriggerID),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	deleteRequest, err := http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating new delete trigger request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, deleteRequest)
	if err != nil {
		return fmt.Errorf("executing delete trigger request: %w", err)
	}

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("deleting deployed trigger %s: %w",
			deployedTriggerID,
			NotFoundErr)
	default:
		return fmt.Errorf("expected status %d, got %d",
			http.StatusNoContent,
			response.StatusCode)
	}
}

func (c *Client) RetrieveTriggerEvents(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
	numberOfEvents int,
) (*TriggerEventList, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", deployedComponentID, "events"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	if numberOfEvents > 0 {
		internal.AddQueryParams(queryParams, "n", strconv.Itoa(numberOfEvents))
	}
	baseURL.RawQuery = queryParams.Encode()

	eventsReq, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating new retrieve trigger events request: %w", err)
	}

	triggerResponse, err := c.doRequestViaOauth(ctx, eventsReq)
	if err != nil {
		return nil, fmt.Errorf("executing request to retrieve trigger events: %w", err)
	}

	var triggerEventList TriggerEventList
	if err := internal.UnmarshalResponse(triggerResponse, &triggerEventList); err != nil {
		return nil, fmt.Errorf("unmarshalling response for trigger events: %w", err)
	}

	return &triggerEventList, nil
}

// ListTriggerWebhooks Retrieve the list of webhook URLs listening to a deployed trigger
func (c *Client) ListTriggerWebhooks(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
) (*TriggerWebhookURLs, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", deployedComponentID, "webhooks"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	eventsReq, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating list trigger webhooks request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, eventsReq)
	if err != nil {
		return nil, fmt.Errorf("executing request to list trigger webhooks: %w", err)
	}
	defer response.Body.Close()

	var webhookUrls TriggerWebhookURLs
	if err := internal.UnmarshalResponse(response, &webhookUrls); err != nil {
		return nil, fmt.Errorf("unmarshalling response for updating trigger webhooks: %w", err)
	}
	return &webhookUrls, nil
}

// UpdateTriggerWebhooks Updates the list of webhook URLs that will listen to a deployed trigger
func (c *Client) UpdateTriggerWebhooks(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
	webhookURLs []string,
) (*TriggerWebhookURLs, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", deployedComponentID, "webhooks"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	body := UpdateTriggerWebhooksRequest{
		WebhookURLs: webhookURLs,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling update trigger webhooks body request: %w", err)
	}

	eventsReq, err := http.NewRequest(http.MethodPut, baseURL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("creating update trigger webhooks request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, eventsReq)
	if err != nil {
		return nil, fmt.Errorf("executing request to update trigger webhooks: %w", err)
	}
	defer response.Body.Close()

	var webhookUrls TriggerWebhookURLs
	if err := internal.UnmarshalResponse(response, &webhookUrls); err != nil {
		return nil, fmt.Errorf("unmarshalling response for updating trigger webhooks: %w", err)
	}
	return &webhookUrls, nil
}

// RetrieveTriggerWorkflows Retrieve the workflows listening to a deployed trigger
func (c *Client) RetrieveTriggerWorkflows(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
) (*TriggerWorkflowIDs, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", deployedComponentID, "pipelines"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating Retrieve trigger workflows request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request to Retrieve trigger workflows: %w", err)
	}
	defer response.Body.Close()

	var workflowIds TriggerWorkflowIDs
	if err := internal.UnmarshalResponse(response, &workflowIds); err != nil {
		return nil, fmt.Errorf("unmarshalling response for Retrieving trigger workflows: %w", err)
	}
	return &workflowIds, nil
}

// UpdateTriggerWorkflows UUpdate the list of workflows that will listen to a deployed trigger
func (c *Client) UpdateTriggerWorkflows(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
	workflowIDs []string,
) (*TriggerWorkflowIDs, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", deployedComponentID, "pipelines"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	body := UpdateTriggerWorkflowsRequest{
		WorkflowIDs: workflowIDs,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling update trigger workflows body request: %w", err)
	}

	eventsReq, err := http.NewRequest(http.MethodPut, baseURL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("creating update trigger workflows request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, eventsReq)
	if err != nil {
		return nil, fmt.Errorf("executing request to update trigger workflows: %w", err)
	}
	defer response.Body.Close()

	var webhookIds TriggerWorkflowIDs
	if err := internal.UnmarshalResponse(response, &webhookIds); err != nil {
		return nil, fmt.Errorf("unmarshalling response for updating trigger workflows: %w", err)
	}
	return &webhookIds, nil
}

// DeployTriggerWithOptions deploys a trigger with full options support.
func (c *Client) DeployTriggerWithOptions(
	ctx context.Context,
	opts *DeployTriggerOptions,
) (*Trigger, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "triggers", "deploy"),
	})

	trigger := DeployTriggerRequest{
		ComponentKey:    opts.ComponentKey,
		ConfiguredProps: opts.ConfiguredProps,
		WebhookURL:      opts.WebhookURL,
		WorkflowID:      opts.WorkflowID,
		DynamicPropsID:  opts.DynamicPropsID,
		ExternalUserID:  opts.ExternalUserID,
		Version:         opts.Version,
		EmitOnDeploy:    opts.EmitOnDeploy,
	}

	jsonBytes, err := json.Marshal(trigger)
	if err != nil {
		return nil, fmt.Errorf("marshalling deploy trigger request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("creating new deploy trigger request: %w", err)
	}

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing deploy trigger request: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body for deploying trigger: %w", err)
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("failed to deploy trigger: %s", string(bodyBytes))
	}

	var response DeployTriggerResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling response for trigger response: %w: %w",
			errors.New(string(bodyBytes)), err)
	}

	return &response.Data, nil
}

// ListDeployedTriggersWithOptions lists deployed triggers with full filtering and pagination.
func (c *Client) ListDeployedTriggersWithOptions(
	ctx context.Context,
	opts *ListDeployedTriggersOptions,
) (*TriggerList, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", opts.ExternalUserID)
	if opts.EmitterType != nil {
		internal.AddQueryParams(queryParams, "emitter_type", *opts.EmitterType)
	}
	internal.AddQueryParamInt(queryParams, "limit", opts.Limit)
	if opts.After != nil {
		internal.AddQueryParams(queryParams, "after", *opts.After)
	}
	if opts.Before != nil {
		internal.AddQueryParams(queryParams, "before", *opts.Before)
	}
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating new request to list deployed triggers: %w", err)
	}

	listResponse, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request to list deployed triggers: %w", err)
	}

	var list TriggerList
	if err := internal.UnmarshalResponse(listResponse, &list); err != nil {
		return nil, fmt.Errorf("unmarshalling response to list: %w", err)
	}

	return &list, nil
}

// DeleteDeployedTriggerWithOptions deletes a deployed trigger with full options support.
func (c *Client) DeleteDeployedTriggerWithOptions(
	ctx context.Context,
	opts *DeleteDeployedTriggerOptions,
) error {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", opts.DeployedTriggerID),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", opts.ExternalUserID)
	internal.AddQueryParamBool(queryParams, "ignore_hook_errors", opts.IgnoreHookErrors)
	baseURL.RawQuery = queryParams.Encode()

	deleteRequest, err := http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating new delete trigger request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, deleteRequest)
	if err != nil {
		return fmt.Errorf("executing delete trigger request: %w", err)
	}

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("deleting deployed trigger %s: %w",
			opts.DeployedTriggerID, NotFoundErr)
	default:
		return fmt.Errorf("expected status %d, got %d",
			http.StatusNoContent, response.StatusCode)
	}
}

type UpdateDeployedTriggerRequest struct {
	Active          *bool           `json:"active,omitempty"`
	ConfiguredProps ConfiguredProps `json:"configured_props,omitempty"`
	Name            string          `json:"name,omitempty"`
	EmitOnDeploy    *bool           `json:"emit_on_deploy,omitempty"`
}

type TriggerWebhookDetail struct {
	ID            string `json:"id,omitempty"`
	URL           string `json:"url,omitempty"`
	SigningKey    string `json:"signing_key,omitempty"`
	SigningKeySet bool   `json:"signing_key_set,omitempty"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
}

type GetTriggerWebhookResponse struct {
	Data TriggerWebhookDetail `json:"data"`
}

// UpdateDeployedTrigger updates a deployed trigger's configuration.
func (c *Client) UpdateDeployedTrigger(
	ctx context.Context,
	triggerID string,
	externalUserID string,
	req *UpdateDeployedTriggerRequest,
) (*Trigger, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", triggerID),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling update deployed trigger request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPut, baseURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating update deployed trigger request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing update deployed trigger request: %w", err)
	}

	var result GetTriggerResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling update deployed trigger response: %w", err)
	}

	return &result.Data, nil
}

// GetTriggerWebhook retrieves details for a specific webhook on a deployed trigger.
func (c *Client) GetTriggerWebhook(
	ctx context.Context,
	triggerID string,
	webhookID string,
	externalUserID string,
) (*TriggerWebhookDetail, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", triggerID, "webhooks", webhookID),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating get trigger webhook request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get trigger webhook request: %w", err)
	}

	var result GetTriggerWebhookResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get trigger webhook response: %w", err)
	}

	return &result.Data, nil
}

// RegenerateTriggerWebhookSigningKey regenerates the signing key for a specific trigger webhook.
func (c *Client) RegenerateTriggerWebhookSigningKey(
	ctx context.Context,
	triggerID string,
	webhookID string,
	externalUserID string,
) (*TriggerWebhookDetail, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "deployed-triggers", triggerID, "webhooks", webhookID, "regenerate_signing_key"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating regenerate trigger webhook signing key request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing regenerate trigger webhook signing key request: %w", err)
	}

	var result GetTriggerWebhookResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling regenerate trigger webhook signing key response: %w", err)
	}

	return &result.Data, nil
}
