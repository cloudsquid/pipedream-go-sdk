package pipedream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type Trigger struct {
	ID                string             `json:"id"`
	OwnerID           string             `json:"owner_id"`
	ComponentID       string             `json:"component_id"`
	ConfigurableProps []ConfigurableProp `json:"configurable_props,omitempty"`
	ConfiguredProps   ConfiguredProps    `json:"configured_props,omitempty"`
	Active            bool               `json:"active,omitempty"`
	CreatedAt         int                `json:"created_at"`
	UpdatedAt         int                `json:"updated_at"`
	Name              string             `json:"name"`
	NameSlug          string             `json:"name_slug"`
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

// OPTIONAL: WorkflowID, DynamicPropsID and WebhookURL
type DeployTriggerRequest struct {
	ComponentKey    string          `json:"id"`
	ConfiguredProps ConfiguredProps `json:"configured_props"`
	ExternalUserID  string          `json:"external_user_id"`
	WebhookURL      string          `json:"webhook_url,omitempty"`
	WorkflowID      string          `json:"workflow_id,omitempty"`
	DynamicPropsID  string          `json:"dynamic_props_id,omitempty"`
}

type UpdateTriggerWebhooksRequest struct {
	ExternalUserID string   `json:"external_user_id"`
	WebhookURLs    []string `json:"webhook_urls,omitempty"`
}

type UpdateTriggerWorkflowsRequest struct {
	ExternalUserID string   `json:"external_user_id"`
	WorkflowIDs    []string `json:"workflow_ids,omitempty"`
}

type TriggerWebhookURLs struct {
	WebhookURLs []string `json:"webhook_urls,omitempty"`
}

type TriggerWorkflowIDs struct {
	WorkflowIDs []string `json:"workflow_ids,omitempty"`
}

func (p *Client) DeployTrigger(
	ctx context.Context,
	componentKey string,
	externalUserID string,
	configuredProps ConfiguredProps,
	webhookURL string,
	dynamicPropsID string, // OPTIONAL
	workflowID string, // OPTIONAL
) (map[string]any, error) {
	p.logger.Info("Deploying trigger")

	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.projectID, "triggers", "deploy")})

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

	resp, err := p.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing deploy trigger request: %w", err)
	}

	var response map[string]any
	if err := unmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("unmarhalling response for trigger response: %w", err)
	}

	return response, nil
}

func (p *Client) ListDeployedTriggers(
	ctx context.Context,
	externalUserID string,
) (*TriggerList, error) {
	p.logger.Info("List Deployed triggers")

	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.projectID, "deployed-triggers")})

	queryParams := url.Values{}
	addQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating new request to list deployed triggers: %w", err)
	}

	listResponse, err := p.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request to list deployed triggers: %w", err)
	}

	var list TriggerList
	if err := unmarshalResponse(listResponse, &list); err != nil {
		return nil, fmt.Errorf("unmarshalling response to list: %w", err)
	}

	return &list, nil
}

func (p *Client) GetDeployedTrigger(
	ctx context.Context,
	deployedComponentID string,
	externalUserId string,
) (*Trigger, error) {
	p.logger.Info("Retrieving deployed trigger", "trigger", deployedComponentID)

	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.projectID, "deployed-triggers", deployedComponentID)})

	queryParams := url.Values{}
	addQueryParams(queryParams, "external_user_id", externalUserId)
	baseURL.RawQuery = queryParams.Encode()

	getRequest, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	response, err := p.doRequestViaOauth(ctx, getRequest)
	if err != nil {
		return nil, fmt.Errorf("executing request to retrieve trigger: %w", err)
	}

	var trigger Trigger
	if err := unmarshalResponse(response, &trigger); err != nil {
		return nil, fmt.Errorf("unmarshalling response to retrieve trigger: %w", err)
	}

	return &trigger, nil
}

func (p *Client) DeleteDeployedTrigger(
	ctx context.Context,
	deployedTriggerID string,
	externalUserID string,
) error {
	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.connectURL.Path, p.projectID, "deployed-triggers", deployedTriggerID)})

	queryParams := url.Values{}
	addQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	deleteRequest, err := http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating new delete trigger request: %w", err)
	}

	response, err := p.doRequestViaOauth(ctx, deleteRequest)
	if err != nil {
		return fmt.Errorf("executing delete trigger request: %w", err)
	}

	if response.StatusCode == http.StatusNoContent {
		return nil
	} else {
		return fmt.Errorf("expected status %d, got %d",
			http.StatusNoContent, response.StatusCode)
	}
}

func (p *Client) RetrieveTriggerEvents(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
	numberOfEvents int,
) (*TriggerEventList, error) {
	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.projectID, "deployed-triggers", deployedComponentID, "events")})

	queryParams := url.Values{}
	addQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	eventsReq, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating new retrieve trigger events request: %w", err)
	}

	triggerResponse, err := p.doRequestViaOauth(ctx, eventsReq)
	if err != nil {
		return nil, fmt.Errorf("executing request to retrieve trigger events: %w", err)
	}

	var triggerEventList TriggerEventList
	if err := unmarshalResponse(triggerResponse, &triggerEventList); err != nil {
		return nil, fmt.Errorf("unmarshalling response for trigger events: %w", err)
	}

	return &triggerEventList, nil
}

// ListTriggerWebhooks Retrieve the list of webhook URLs listening to a deployed trigger
func (p *Client) ListTriggerWebhooks(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
) (*TriggerWebhookURLs, error) {
	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.projectID, "deployed-triggers", deployedComponentID, "webhooks")})

	queryParams := url.Values{}
	addQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	eventsReq, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating list trigger webhooks request: %w", err)
	}

	response, err := p.doRequestViaOauth(ctx, eventsReq)
	if err != nil {
		return nil, fmt.Errorf("executing request to list trigger webhooks: %w", err)
	}
	defer response.Body.Close()

	var webhookUrls TriggerWebhookURLs
	if err := unmarshalResponse(response, &webhookUrls); err != nil {
		return nil, fmt.Errorf("unmarshalling response for updating trigger webhooks: %w", err)
	}
	return &webhookUrls, nil
}

// UpdateTriggerWebhooks Updates the list of webhook URLs that will listen to a deployed trigger
func (p *Client) UpdateTriggerWebhooks(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
	webhookURLs []string,
) (*TriggerWebhookURLs, error) {
	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.projectID, "deployed-triggers", deployedComponentID, "webhooks")})

	body := UpdateTriggerWebhooksRequest{
		ExternalUserID: externalUserID,
		WebhookURLs:    webhookURLs,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling update trigger webhooks body request: %w", err)
	}

	eventsReq, err := http.NewRequest(http.MethodPut, baseURL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("creating update trigger webhooks request: %w", err)
	}

	response, err := p.doRequestViaOauth(ctx, eventsReq)
	if err != nil {
		return nil, fmt.Errorf("executing request to update trigger webhooks: %w", err)
	}
	defer response.Body.Close()

	var webhookUrls TriggerWebhookURLs
	if err := unmarshalResponse(response, &webhookUrls); err != nil {
		return nil, fmt.Errorf("unmarshalling response for updating trigger webhooks: %w", err)
	}
	return &webhookUrls, nil
}

// RetrieveTriggerWorkflows Retrieve the workflows listening to a deployed trigger
func (p *Client) RetrieveTriggerWorkflows(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
) (*TriggerWorkflowIDs, error) {
	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.projectID, "deployed-triggers", deployedComponentID, "workflows")})

	queryParams := url.Values{}
	addQueryParams(queryParams, "external_user_id", externalUserID)
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating Retrieve trigger workflows request: %w", err)
	}

	response, err := p.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request to Retrieve trigger workflows: %w", err)
	}
	defer response.Body.Close()

	var workflowIds TriggerWorkflowIDs
	if err := unmarshalResponse(response, &workflowIds); err != nil {
		return nil, fmt.Errorf("unmarshalling response for Retrieving trigger workflows: %w", err)
	}
	return &workflowIds, nil
}

// UpdateTriggerWorkflows UUpdate the list of workflows that will listen to a deployed trigger
func (p *Client) UpdateTriggerWorkflows(
	ctx context.Context,
	deployedComponentID string,
	externalUserID string,
	workflowIDs []string,
) (*TriggerWorkflowIDs, error) {
	baseURL := p.connectURL.ResolveReference(&url.URL{
		Path: path.Join(p.connectURL.Path, p.projectID, "deployed-triggers", deployedComponentID, "workflows")})

	body := UpdateTriggerWorkflowsRequest{
		ExternalUserID: externalUserID,
		WorkflowIDs:    workflowIDs,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling update trigger workflows body request: %w", err)
	}

	eventsReq, err := http.NewRequest(http.MethodPut, baseURL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("creating update trigger workflows request: %w", err)
	}

	response, err := p.doRequestViaOauth(ctx, eventsReq)
	if err != nil {
		return nil, fmt.Errorf("executing request to update trigger workflows: %w", err)
	}
	defer response.Body.Close()

	var webhookIds TriggerWorkflowIDs
	if err := unmarshalResponse(response, &webhookIds); err != nil {
		return nil, fmt.Errorf("unmarshalling response for updating trigger workflows: %w", err)
	}
	return &webhookIds, nil
}
