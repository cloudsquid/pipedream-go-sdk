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
		return nil, fmt.Errorf("executing create component request: %w", err)
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
