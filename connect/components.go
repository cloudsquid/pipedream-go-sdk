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
	"strings"

	"github.com/cloudsquid/pipedream-go-sdk/internal"
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
	Observations  []any    `json:"observations,omitempty"`
	Context       any      `json:"context,omitempty"` // TODO
	Options       []Value  `json:"options,omitempty"`
	Errors        []string `json:"errors,omitempty"`
	StringOptions any      `json:"string_options,omitempty"`
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

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var propOptions PropOptions
	if err := json.Unmarshal(bodyBytes, &propOptions); err != nil {
		return nil, fmt.Errorf("unmarshalling body into propOptions: %w: %w",
			errors.New(string(bodyBytes)), err)
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

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body: %w", err)
		}

		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

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
	limit int,
) (*ListComponentResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(componentType))})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "app", appName)
	internal.AddQueryParams(queryParams, "q", searchTerm)

	if limit > 0 {
		internal.AddQueryParams(queryParams, "limit", strconv.Itoa(limit))
	}

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
