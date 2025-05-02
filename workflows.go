package pipedream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type CreateWorkflowRequest struct {
	OrgID      string            `json:"org_id"`
	ProjectID  string            `json:"project_id"`
	TemplateID string            `json:"template_id"`
	Steps      []WorkflowStep    `json:"steps,omitempty"`
	Triggers   []WorkflowTrigger `json:"triggers,omitempty"`
	Settings   *WorkflowSettings `json:"settings,omitempty"`
}

type WorkflowStep struct {
	Namespace string                 `json:"namespace"`
	Props     map[string]interface{} `json:"props"`
}

type WorkflowTrigger struct {
	Props map[string]interface{} `json:"props"`
}

type WorkflowSettings struct {
	Name       string `json:"name,omitempty"`
	AutoDeploy bool   `json:"auto_deploy,omitempty"`
}

type CreateWorkflowResponse struct {
	Data Workflow `json:"data"`
}

type Workflow struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Active   bool               `json:"active"`
	Steps    []WorkflowStepInfo `json:"steps"`
	Triggers []TriggerInfo      `json:"triggers"`
}

type WorkflowStepInfo struct {
	ID                       string            `json:"id"`
	Type                     string            `json:"type"`
	Namespace                string            `json:"namespace"`
	Disabled                 bool              `json:"disabled"`
	CodeRaw                  *string           `json:"code_raw"`
	CodeRawAlt               *string           `json:"codeRaw"`
	CodeConfigJSON           *string           `json:"codeConfigJson"`
	Lang                     string            `json:"lang"`
	TextRaw                  *string           `json:"text_raw"`
	AppConnections           []string          `json:"appConnections"`
	FlatParamsVisibilityJSON *string           `json:"flat_params_visibility_json"`
	ParamsJSON               string            `json:"params_json"`
	Component                bool              `json:"component"`
	SavedComponent           *SavedComponent   `json:"savedComponent,omitempty"`
	ComponentKey             *string           `json:"component_key"`
	ComponentOwnerID         *string           `json:"component_owner_id"`
	ConfiguredPropsJSON      string            `json:"configured_props_json"`
	AuthProvisionIdMap       map[string]string `json:"authProvisionIdMap"`
	AuthProvisionIds         []string          `json:"authProvisionIds"`
}

type SavedComponent struct {
	ID                string              `json:"id"`
	Code              string              `json:"code"`
	CodeHash          string              `json:"codeHash"`
	ConfigurableProps []ComponentProp     `json:"configurableProps"`
	Key               *string             `json:"key"`
	Description       *string             `json:"description"`
	EntryPath         *string             `json:"entryPath"`
	Version           string              `json:"version"`
	Apps              []map[string]string `json:"apps"`
}

type ComponentProp struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TriggerInfo struct {
	ID              string                 `json:"id"`
	OwnerID         string                 `json:"owner_id"`
	ComponentID     string                 `json:"component_id"`
	ConfiguredProps map[string]interface{} `json:"configured_props"`
	Active          bool                   `json:"active"`
	CreatedAt       int64                  `json:"created_at"`
	UpdatedAt       int64                  `json:"updated_at"`
	Name            string                 `json:"name"`
	NameSlug        string                 `json:"name_slug"`
	Key             string                 `json:"key,omitempty"`
	CustomResponse  bool                   `json:"custom_response,omitempty"`
	EndpointURL     string                 `json:"endpoint_url,omitempty"`
}

type UpdateWorkflowRequest struct {
	Active bool   `json:"active"`
	OrgID  string `json:"org_id"`
}

type GetWorkflowDetailsResponse struct {
	Triggers []TriggerInfo      `json:"triggers"`
	Steps    []WorkflowStepInfo `json:"steps"`
}

type GetWorkflowEmitsResponse struct {
	PageInfo PageInfo       `json:"page_info"`
	Data     []EventSummary `json:"data"`
}

type EventSummary struct {
	ID        string   `json:"id"`
	IndexedAt int64    `json:"indexed_at_ms"`
	Event     RawEvent `json:"event"`
	Metadata  Metadata `json:"metadata"`
}

type RawEvent struct {
	RawEvent map[string]interface{} `json:"raw_event"`
}

type Metadata struct {
	EmitID    string `json:"emit_id"`
	Name      string `json:"name"`
	EmitterID string `json:"emitter_id"`
}

type GetWorkflowErrorsResponse struct {
	PageInfo PageInfo        `json:"page_info"`
	Data     []WorkflowError `json:"data"`
}

type WorkflowError struct {
	ID              string                 `json:"id"`
	IndexedAtMS     int64                  `json:"indexed_at_ms"`
	Event           map[string]interface{} `json:"event"`
	OriginalContext ErrorOriginalContext   `json:"original_context"`
	Error           WorkflowExecutionError `json:"error"`
	Metadata        Metadata               `json:"metadata"`
}

type ErrorOriginalContext struct {
	ID              string `json:"id"`
	Timestamp       string `json:"ts"`
	WorkflowID      string `json:"workflow_id"`
	DeploymentID    string `json:"deployment_id"`
	SourceType      string `json:"source_type"`
	Verified        bool   `json:"verified"`
	OwnerID         string `json:"owner_id"`
	PlatformVersion string `json:"platform_version"`
}

type WorkflowExecutionError struct {
	Code   string `json:"code"`
	CellID string `json:"cellId"`
	TS     string `json:"ts"`
	Stack  string `json:"stack"`
}

// TODO: implement invoke workflow
// CreateWorkflow Creates a new workflow within an organization’s project
// This endpoint allows defining workflow steps, triggers, and settings, based on a supplied template
func (c *Client) CreateWorkflow(
	ctx context.Context,
	orgID,
	projectID,
	templateID string,
	steps []WorkflowStep,
	triggers []WorkflowTrigger,
	settings *WorkflowSettings,
) (*CreateWorkflowResponse, error) {
	c.logger.Debug("creating workflow")

	endpoint := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "workflows"),
	}).String()

	payload := &CreateWorkflowRequest{
		OrgID:      orgID,
		ProjectID:  projectID,
		TemplateID: templateID,
		Steps:      steps,
		Triggers:   triggers,
		Settings:   settings,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling create workflow request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating create workflow request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing create workflow request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var respJson CreateWorkflowResponse
	if err := json.NewDecoder(response.Body).Decode(&respJson); err != nil {
		return nil, fmt.Errorf("decoding create workflow response: %w", err)
	}

	return &respJson, nil
}

// UpdateWorkflow Updates the workflow’s activation status
// Does not modify the workflow’s steps, triggers, or connected accounts
func (c *Client) UpdateWorkflow(
	ctx context.Context,
	id,
	orgID string,
	active bool,
) (*map[string]any, error) {
	c.logger.Debug("update workflow")

	endpoint := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "workflows", id),
	}).String()

	payload := &UpdateWorkflowRequest{
		OrgID:  orgID,
		Active: active,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling update workflow request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating update workflow request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing update workflow request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var result map[string]any
	if err := unmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling update workflow response:e: %w", err)
	}

	return &result, nil
}

// GetWorkflowDetails Retrieves the details of a specific workflow within an organization’s project
func (c *Client) GetWorkflowDetails(
	ctx context.Context,
	id,
	orgID string,
) (*GetWorkflowDetailsResponse, error) {
	c.logger.Debug("get workflow details")

	if orgID == "" {
		return nil, fmt.Errorf("orgID is required")
	}

	baseURL := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "workflows", id),
	})

	queryParams := url.Values{}
	addQueryParams(queryParams, "org_id", orgID)
	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get workflow details request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get workflow details request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var result GetWorkflowDetailsResponse
	if err := unmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get workflow details response:e: %w", err)
	}

	return &result, nil
}

// GetWorkflowEmits Retrieves up to the last 100 events emitted from a workflow using $send.emit().
func (c *Client) GetWorkflowEmits(
	ctx context.Context,
	id,
	orgID string,
	expandEvent bool,
	limit int, // if 0 no limit is applied
) (*GetWorkflowEmitsResponse, error) {
	c.logger.Debug("get workflow emits")

	if orgID == "" {
		return nil, fmt.Errorf("orgID is required")
	}
	baseURL := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "workflows", id, "event_summaries"),
	})

	queryParams := url.Values{}
	addQueryParams(queryParams, "org_id", orgID)

	if expandEvent {
		addQueryParams(queryParams, "expand", "event")
	}

	if limit > 0 {
		addQueryParams(queryParams, "limit", strconv.Itoa(limit))
	}

	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get workflow emits request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get workflow emits request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(raw))
	}

	var result GetWorkflowEmitsResponse
	if err := unmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get workflow details response:e: %w", err)
	}

	return &result, nil
}

// GetWorkflowErrors Retrieve up to the last 100 events for a workflow that threw an error
// The details of the error, along with the original event data, will be included
func (c *Client) GetWorkflowErrors(
	ctx context.Context,
	id string,
	expandEvent bool,
	limit int,
) (*GetWorkflowErrorsResponse, error) {
	c.logger.Debug("getting workflow errors")

	baseURL := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "workflows", id, "$errors", "event_summaries"),
	})

	queryParams := url.Values{}
	if expandEvent {
		addQueryParams(queryParams, "expand", "event")
	}
	if limit > 0 {
		addQueryParams(queryParams, "limit", strconv.Itoa(limit))
	}

	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get workflow emits request: %w", err)
	}

	resp, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request for get workflow errors reaquest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(raw))
	}

	var result GetWorkflowErrorsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding get workflow errors response: %w", err)
	}

	return &result, nil
}
