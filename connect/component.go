package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type ComponentType string

const (
	Triggers   ComponentType = "triggers"
	Actions    ComponentType = "actions"
	Components ComponentType = "components"
)

type Component struct {
	Key         string        `json:"key,omitempty"`
	Name        string        `json:"name,omitempty"`
	Version     string        `json:"version,omitempty"`
	Type        ComponentType `json:"type,omitempty"`
	Description string        `json:"description,omitempty"`
}

func (c Component) String() string {
	return fmt.Sprintf("%-20s\t%-30s\t%-50s", c.Key, c.Name, c.Description)
}

type ConfigurableProp struct {
	Name           string `json:"name,omitempty"`
	Type           string `json:"type"`
	App            string `json:"app,omitempty"`
	CustomResponse bool   `json:"custom_response,omitempty"`
	Label          string `json:"label,omitempty"`
	Description    string `json:"description,omitempty"`
	RemoteOptions  *bool  `json:"remoteOptions,omitempty"`
	Options        []any  `json:"options,omitempty"` // this can be a string array or an object array of Value
	UseQuery       bool   `json:"use_query,omitempty"`
	Default        any    `json:"default,omitempty"`
	Min            int    `json:"min,omitempty"`
	Max            int    `json:"max,omitempty"`
	Disabled       bool   `json:"disabled,omitempty"`
	Secret         bool   `json:"secret,omitempty"`
	Optional       bool   `json:"optional,omitempty"`
	ReloadProps    bool   `json:"reloadProps,omitempty"`
}

func (c ConfigurableProp) String() string {
	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Name: %s\tDescription: %s\tOptions: %s",
			c.Name, c.Description, c.Options)
	}

	return string(bs)
}

type ConfiguredProps map[string]any

func (c ConfiguredProps) String() string {
	bs, _ := json.Marshal(c)
	return string(bs)
}

type ComponentDetails struct {
	Component
	ConfigurableProps []*ConfigurableProp `json:"configurable_props,omitempty"`
}

func (c ComponentDetails) String() string {
	output, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Key: %s\t Description: %s", c.Key, c.Description)
	}

	return string(output)
}

type PropOptions struct {
	Observations  []any            `json:"observations,omitempty"`
	Context       any              `json:"context,omitempty"` // TODO
	Options       []Value          `json:"options,omitempty"`
	Errors        []string         `json:"errors,omitempty"`
	Timings       []map[string]any `json:"timings,omitempty"`
	StringOptions any              `json:"string_options,omitempty"`
}

type DynamicProps struct {
	ID                string             `json:"id,omitempty"`
	ConfigurableProps []ConfigurableProp `json:"configurableProps,omitempty"`
}

type ReloadComponentPropsResponse struct {
	Observations []any        `json:"observations,omitempty"`
	Errors       []string     `json:"errors,omitempty"`
	DynamicProps DynamicProps `json:"dynamicProps"`
}

func (p PropOptions) String() string {
	bs, _ := json.MarshalIndent(p, "", "  ")
	return string(bs)
}

type Value struct {
	Label string `json:"label,omitempty"`
	Value any    `json:"value,omitempty"`
}

func (v Value) String() string {
	return fmt.Sprintf("\t%s\t%v", v.Label, v.Value)
}

type ReloadComponentPropsRequest struct {
	ExternalUserID  string          `json:"external_user_id,omitempty"`
	ConfiguredProps ConfiguredProps `json:"configured_props,omitempty"`
	ID              string          `json:"id,omitempty"`
	DynamicPropsID  string          `json:"dynamic_props_id,omitempty"`
}

// GetComponentResponse is the response for the get component endpoint
type GetComponentResponse struct {
	Data *ComponentDetails `json:"data,omitempty"`
}

// ListComponentResponse is the response for the component list endpoint
type ListComponentResponse struct {
	Data []*Component `json:"data,omitempty"`
}

type CreateComponentRequest struct {
	ComponentCode string `json:"component_code,omitempty"`
	ComponentURL  string `json:"component_url,omitempty"`
}

type NewComponent struct {
	ID                string             `json:"id"`
	Code              string             `json:"code"`
	CodeHash          string             `json:"code_hash"`
	Name              string             `json:"name"`
	Version           string             `json:"version"`
	ConfigurableProps []ConfigurableProp `json:"configurable_props,omitempty"`
	CreatedAt         int64              `json:"created_at"`
	UpdatedAt         int64              `json:"updated_at"`
}

type CreateComponentResponse struct {
	Data *NewComponent `json:"data,omitempty"`
}

type ComponentSearchResponse struct {
	Sources []string `json:"sources"`
	Actions []string `json:"actions"`
}

// https://pipedream.com/docs/connect/api/#configure-a-component
// ConfigureComponent calls the configure endpoint for a component in pipedream
// externalUserID is the id defined by a third party or us
// component Key is the componentID
// propName is the key in the componentDetails
func (c *Client) GetPropOptions(
	ctx context.Context,
	propName string,
	componentKey string,
	externalUserID string,
	configuredProps ConfiguredProps,
) (*PropOptions, error) {
	c.Logger.Info("Getting options for the prop")

	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "components", "configure")})

	endpoint := baseURL.String()

	type ConfigureRequest struct {
		ExternalUserID  string          `json:"external_user_id,omitempty"`
		ComponentKey    string          `json:"id,omitempty"`
		PropName        string          `json:"prop_name,omitempty"`
		ConfiguredProps ConfiguredProps `json:"configured_props,omitempty"`
	}

	requestBody := &ConfigureRequest{
		ExternalUserID:  externalUserID,
		ComponentKey:    componentKey,
		PropName:        propName,
		ConfiguredProps: configuredProps,
	}
	if propName == "" {
		requestBody = &ConfigureRequest{
			ExternalUserID:  externalUserID,
			ComponentKey:    componentKey,
			ConfiguredProps: configuredProps,
		}
	}

	bs, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil,
			fmt.Errorf("executing request to configure component %s for user %s: %w",
				componentKey, externalUserID, err)
	}
	defer response.Body.Close()

	var propOptions PropOptions
	if err := internal.UnmarshalResponse(response, &propOptions); err != nil {
		return nil, fmt.Errorf("unmarshalling prop options for component %s: %w",
			componentKey, err)
	}

	if propOptions.Errors != nil || len(propOptions.Errors) > 0 {
		return nil, errors.New(strings.Join(propOptions.Errors, "."))
	}

	return &propOptions, nil
}

// https://pipedream.com/docs/connect/api/#retrieve-a-component
// GetComponent retrieves a pipedream component and its configurable props
func (c *Client) GetComponent(
	ctx context.Context,
	componentKey string,
	componentType ComponentType,
) (*GetComponentResponse, error) {
	c.Logger.Info("Getting component details",
		"component", componentKey,
		"type", componentType)

	endpoint := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(componentType), componentKey)}).String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new get request for endpoint %s: %w",
			endpoint, err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	var component GetComponentResponse
	if err := internal.UnmarshalResponse(response, &component); err != nil {
		return nil, fmt.Errorf(
			"parsing response for getting component details for component %s: %w",
			componentKey, err)
	}

	return &component, nil
}

// https://pipedream.com/docs/connect/api/#list-components
// ListComponents lists the components available in pipedream
func (c *Client) ListComponents(
	ctx context.Context,
	componentType ComponentType,
	appName string,
	searchTerm string,
) (*ListComponentResponse, error) {
	c.Logger.Info("Listing components",
		"componentType", componentType,
		"appName", appName,
	)

	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(componentType))})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "app", appName)
	internal.AddQueryParams(queryParams, "q", searchTerm)

	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get request for endpoint %s: %w", endpoint, err)
	}

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	var respJson ListComponentResponse
	if err := internal.UnmarshalResponse(resp, &respJson); err != nil {
		return nil, fmt.Errorf(
			"parsing response for listing components for app %s: %w",
			appName, err)
	}

	return &respJson, nil
}

// ReloadComponentProps Reload the component’s props after configuring a dynamic prop,
// based on the current component’s configuration
// will use the component’s configuration to retrieve a new list of props depending on the value of the props that were configured so far
func (c *Client) ReloadComponentProps(
	ctx context.Context,
	componentType ComponentType,
	configuredProps ConfiguredProps,
	externalUserID string,
	ComponentKey string,
	dynamicPropsID string,
) (*ReloadComponentPropsResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(componentType), "props")})

	endpoint := baseURL.String()

	requestBody := &ReloadComponentPropsRequest{
		ExternalUserID:  externalUserID,
		ID:              ComponentKey,
		ConfiguredProps: configuredProps,
		DynamicPropsID:  dynamicPropsID,
	}

	bs, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling reload component props body request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("creating reload component props request: %w", err)
	}

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing reload component props request: %w", err)
	}
	defer resp.Body.Close()

	var respJson ReloadComponentPropsResponse
	if err := internal.UnmarshalResponse(resp, &respJson); err != nil {
		return nil, fmt.Errorf(
			"parsing response for reloading component props: %w", err)
	}

	return &respJson, nil
}

// TODO: MOVE TO REST
// CreateComponent returns the components id, code, configurable_props, and other metadata you’ll need to deploy a source from this component
/*
func (c *Client) CreateComponent(
	ctx context.Context,
	componentCode string,
	componentURL string,
) (*CreateComponentResponse, error) {
	c.Logger.Info("Create component")

	if componentCode == "" && componentURL == "" {
		return nil, fmt.Errorf("either componentCode or componentURL must be provided")
	}
	baseURL := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "components")})
	endpoint := baseURL.String()

	payload := &CreateComponentRequest{ComponentCode: componentCode, ComponentURL: componentURL}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling create component body request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating create component request %s: %w", endpoint, err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing create component request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

	var respJson CreateComponentResponse
	if err := internal.unmarshalResponse(response, &respJson); err != nil {
		return nil, fmt.Errorf("parsing reponse for create component request: %w", err)
	}

	return &respJson, nil
}

// GetComponentFromRegistry returns the same data as the endpoint for retrieving metadata on a component you own, but allows you to fetch data for any globally-published component
func (c *pipedream.Client) GetRegistryComponents(
	ctx context.Context,
	componentKey string,
) (*CreateComponentResponse, error) {
	c.logger.Info("Getting component from registry")

	endpoint := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "components", "registry", componentKey)}).String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new get component from registery request %s: %w",
			endpoint, err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

	var component CreateComponentResponse
	if err := internal.unmarshalResponse(response, &component); err != nil {
		return nil, fmt.Errorf(
			"parsing response for getting component details for component %s: %w",
			componentKey, err)
	}

	return &component, nil
}

// SearchRegistryComponents Search for components in the global registry with natural language
func (c *pipedream.Client) SearchRegistryComponents(
	ctx context.Context,
	query string,
	app string,
	similarityThreshold int,
	debug bool,
) (*ComponentSearchResponse, error) {
	c.logger.Info("searching registry component")

	baseURL := c.baseURL.ResolveReference(&url.URL{
		Path: path.Join(c.baseURL.Path, "components", "search")})

	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	queryParams := url.Values{}
	internal.addQueryParams(queryParams, "query", query)
	if app != "" {
		internal.addQueryParams(queryParams, "app", app)
	}
	if similarityThreshold > 0 {
		internal.addQueryParams(queryParams, "similarity_threshold", fmt.Sprintf("%d", similarityThreshold))
	}
	if debug {
		internal.addQueryParams(queryParams, "debug", "true")
	}
	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new search registry component request %s: %w",
			endpoint, err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing search registry component request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

	var component ComponentSearchResponse
	if err := internal.unmarshalResponse(response, &component); err != nil {
		return nil, fmt.Errorf(
			"parsing response for search registry component request: %w", err)
	}

	return &component, nil
}
*/
