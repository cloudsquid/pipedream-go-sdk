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

type ComponentAnnotations struct {
	DestructiveHint *bool  `json:"destructiveHint,omitempty"`
	IdempotentHint  *bool  `json:"idempotentHint,omitempty"`
	OpenWorldHint   *bool  `json:"openWorldHint,omitempty"`
	ReadOnlyHint    *bool  `json:"readOnlyHint,omitempty"`
	Title           string `json:"title,omitempty"`
}

type Component struct {
	Key           string                `json:"key,omitempty"`
	Name          string                `json:"name,omitempty"`
	Version       string                `json:"version,omitempty"`
	Type          ComponentType         `json:"type,omitempty"`
	Description   string                `json:"description,omitempty"`
	ComponentType string                `json:"component_type,omitempty"`
	Stash         string                `json:"stash,omitempty"`
	Annotations   *ComponentAnnotations `json:"annotations,omitempty"`
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
	Disabled    bool  `json:"disabled,omitempty"`
	Hidden      *bool `json:"hidden,omitempty"`
	Secret      bool  `json:"secret,omitempty"`
	Optional    bool  `json:"optional,omitempty"`
	ReloadProps bool  `json:"reloadProps,omitempty"`
	WithLabel   *bool `json:"withLabel,omitempty"`
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

func (p *PropOptions) UnmarshalJSON(data []byte) error {
	type Alias PropOptions
	aux := &struct {
		StringOptionsCamel any `json:"stringOptions,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.StringOptionsCamel != nil && p.StringOptions == nil {
		p.StringOptions = aux.StringOptionsCamel
	}

	return nil
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
	Version         string          `json:"version,omitempty"`
	Blocking        *bool           `json:"blocking,omitempty"`
}

type ConfigureComponentRequest struct {
	ExternalUserID  string          `json:"external_user_id,omitempty"`
	ComponentKey    string          `json:"id,omitempty"`
	PropName        string          `json:"prop_name,omitempty"`
	ConfiguredProps ConfiguredProps `json:"configured_props,omitempty"`
	Version         string          `json:"version,omitempty"`
	Blocking        *bool           `json:"blocking,omitempty"`
	DynamicPropsID  string          `json:"dynamic_props_id,omitempty"`
	Page            *int            `json:"page,omitempty"`
	PrevContext     any             `json:"prev_context,omitempty"`
	Query           string          `json:"query,omitempty"`
}

type ListComponentsOptions struct {
	ComponentType ComponentType
	App           string
	Q             string
	Limit         *int
	After         *string
	Before        *string
	Registry      *string
}

type GetComponentOptions struct {
	ComponentKey  string
	ComponentType ComponentType
	Version       *string
}

type ReloadComponentPropsOptions struct {
	ComponentType   ComponentType
	ConfiguredProps ConfiguredProps
	ExternalUserID  string
	ComponentKey    string
	DynamicPropsID  string
	Version         string
	Blocking        *bool
}

type GetPropOptionsParams struct {
	PropName        string
	ComponentKey    string
	ExternalUserID  string
	ConfiguredProps ConfiguredProps
	ComponentType   ComponentType
	Version         string
	Blocking        *bool
	DynamicPropsID  string
	Page            *int
	PrevContext      any
	Query           string
}

// GetComponentResponse is the response for the get component endpoint
type GetComponentResponse struct {
	Data *ComponentDetails `json:"data,omitempty"`
}

// ListComponentResponse is the response for the component list endpoint
type ListComponentResponse struct {
	PageInfo PageInfo     `json:"page_info,omitzero"`
	Data     []*Component `json:"data,omitempty"`
}

// Deprecated: Use GetPropOptionsWithParams instead.
func (c *Client) GetPropOptions(
	ctx context.Context,
	propName string,
	componentKey string,
	externalUserID string,
	configuredProps ConfiguredProps,
) (*PropOptions, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "components", "configure"),
	})

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

// Deprecated: Use GetComponentWithOptions instead.
func (c *Client) GetComponent(
	ctx context.Context,
	componentKey string,
	componentType ComponentType,
) (*GetComponentResponse, error) {
	endpoint := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(componentType), componentKey),
	}).String()

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

// Deprecated: Use ListComponentsWithOptions instead.
func (c *Client) ListComponents(
	ctx context.Context,
	componentType ComponentType,
	appName string,
	searchTerm string,
	limit int,
) (*ListComponentResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(componentType)),
	})

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

// Deprecated: Use ReloadComponentPropsWithOptions instead.
func (c *Client) ReloadComponentProps(
	ctx context.Context,
	componentType ComponentType,
	configuredProps ConfiguredProps,
	externalUserID string,
	ComponentKey string,
	dynamicPropsID string,
) (*ReloadComponentPropsResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(componentType), "props"),
	})

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

// GetPropOptionsWithParams configures a component prop and retrieves its available options.
// Supports actions/configure, triggers/configure, and components/configure paths via ComponentType.
func (c *Client) GetPropOptionsWithParams(
	ctx context.Context,
	params *GetPropOptionsParams,
) (*PropOptions, error) {
	componentType := params.ComponentType
	if componentType == "" {
		componentType = Components
	}

	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(componentType), "configure"),
	})

	endpoint := baseURL.String()

	requestBody := &ConfigureComponentRequest{
		ExternalUserID:  params.ExternalUserID,
		ComponentKey:    params.ComponentKey,
		PropName:        params.PropName,
		ConfiguredProps: params.ConfiguredProps,
		Version:         params.Version,
		Blocking:        params.Blocking,
		DynamicPropsID:  params.DynamicPropsID,
		Page:            params.Page,
		PrevContext:     params.PrevContext,
		Query:           params.Query,
	}

	bs, err := json.Marshal(requestBody)
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
				params.ComponentKey, params.ExternalUserID, err)
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

	if len(propOptions.Errors) > 0 {
		return nil, errors.New(strings.Join(propOptions.Errors, "."))
	}

	return &propOptions, nil
}

// GetComponentWithOptions retrieves a component and its configurable props with additional options.
func (c *Client) GetComponentWithOptions(
	ctx context.Context,
	opts *GetComponentOptions,
) (*GetComponentResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(opts.ComponentType), opts.ComponentKey),
	})

	queryParams := url.Values{}
	if opts.Version != nil {
		internal.AddQueryParams(queryParams, "version", *opts.Version)
	}
	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new get request for endpoint %s: %w", endpoint, err)
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
			opts.ComponentKey, err)
	}

	return &component, nil
}

// ListComponentsWithOptions lists components with full pagination and filtering support.
func (c *Client) ListComponentsWithOptions(
	ctx context.Context,
	opts *ListComponentsOptions,
) (*ListComponentResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(opts.ComponentType)),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "app", opts.App)
	internal.AddQueryParams(queryParams, "q", opts.Q)
	internal.AddQueryParamInt(queryParams, "limit", opts.Limit)
	if opts.After != nil {
		internal.AddQueryParams(queryParams, "after", *opts.After)
	}
	if opts.Before != nil {
		internal.AddQueryParams(queryParams, "before", *opts.Before)
	}
	if opts.Registry != nil {
		internal.AddQueryParams(queryParams, "registry", *opts.Registry)
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
		return nil, fmt.Errorf("parsing response for listing components: %w", err)
	}

	return &respJson, nil
}

// ReloadComponentPropsWithOptions reloads component props with full options support.
func (c *Client) ReloadComponentPropsWithOptions(
	ctx context.Context,
	opts *ReloadComponentPropsOptions,
) (*ReloadComponentPropsResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), string(opts.ComponentType), "props"),
	})

	endpoint := baseURL.String()

	requestBody := &ReloadComponentPropsRequest{
		ExternalUserID:  opts.ExternalUserID,
		ID:              opts.ComponentKey,
		ConfiguredProps: opts.ConfiguredProps,
		DynamicPropsID:  opts.DynamicPropsID,
		Version:         opts.Version,
		Blocking:        opts.Blocking,
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
		return nil, fmt.Errorf("parsing response for reloading component props: %w", err)
	}

	return &respJson, nil
}
