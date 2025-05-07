package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
	"io"
	"net/http"
	"net/url"
	"path"
)

type CreateComponentRequest struct {
	ComponentCode string `json:"component_code,omitempty"`
	ComponentURL  string `json:"component_url,omitempty"`
}
type ConfigurableProp struct {
	Name           string `json:"name,omitempty"`
	Type           string `json:"type"`
	App            string `json:"app,omitempty"`
	CustomResponse bool   `json:"custom_response,omitempty"`
	Label          string `json:"label,omitempty"`
	Description    string `json:"description,omitempty"`
	RemoteOptions  *bool  `json:"remoteOptions,omitempty"`
	Options        []any  `json:"options,omitempty"`
	UseQuery       bool   `json:"use_query,omitempty"`
	Default        any    `json:"default,omitempty"`
	Min            int    `json:"min,omitempty"`
	Max            int    `json:"max,omitempty"`
	Disabled       bool   `json:"disabled,omitempty"`
	Secret         bool   `json:"secret,omitempty"`
	Optional       bool   `json:"optional,omitempty"`
	ReloadProps    bool   `json:"reloadProps,omitempty"`
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

// CreateComponent returns the components id, code, configurable_props, and other metadata you’ll need to deploy a source from this component
func (c *Client) CreateComponent(
	ctx context.Context,
	componentCode string,
	componentURL string,
) (*CreateComponentResponse, error) {
	c.Logger.Info("Create component")

	if componentCode == "" && componentURL == "" {
		return nil, fmt.Errorf("either componentCode or componentURL must be provided")
	}

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "components")})
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

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing create component request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

	var respJson CreateComponentResponse
	if err := internal.UnmarshalResponse(response, &respJson); err != nil {
		return nil, fmt.Errorf("parsing reponse for create component request: %w", err)
	}

	return &respJson, nil
}

// GetComponentFromRegistry returns the same data as the endpoint for retrieving metadata on a component you own, but allows you to fetch data for any globally-published component
func (c *Client) GetRegistryComponents(
	ctx context.Context,
	componentKey string,
) (*CreateComponentResponse, error) {
	c.Logger.Info("Getting component from registry")

	endpoint := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "components", "registry", componentKey)}).String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new get component from registery request %s: %w",
			endpoint, err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

	var component CreateComponentResponse
	if err := internal.UnmarshalResponse(response, &component); err != nil {
		return nil, fmt.Errorf(
			"parsing response for getting component details for component %s: %w",
			componentKey, err)
	}

	return &component, nil
}

// SearchRegistryComponents Search for components in the global registry with natural language
func (c *Client) SearchRegistryComponents(
	ctx context.Context,
	query string,
	app string,
	similarityThreshold int,
	debug bool,
) (*ComponentSearchResponse, error) {
	c.Logger.Info("searching registry component")

	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "components", "search")})

	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "query", query)
	internal.AddQueryParams(queryParams, "app", app)

	if similarityThreshold > 0 {
		internal.AddQueryParams(queryParams, "similarity_threshold", fmt.Sprintf("%d", similarityThreshold))
	}

	if debug {
		internal.AddQueryParams(queryParams, "debug", "true")
	}

	baseURL.RawQuery = queryParams.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new search registry component request %s: %w",
			endpoint, err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing search registry component request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d:%s", response.StatusCode, string(bodyBytes))
	}

	var component ComponentSearchResponse
	if err := internal.UnmarshalResponse(response, &component); err != nil {
		return nil, fmt.Errorf(
			"parsing response for search registry component request: %w", err)
	}

	return &component, nil
}
