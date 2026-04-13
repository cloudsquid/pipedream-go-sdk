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

type Project struct {
	ID                         string  `json:"id,omitempty"`
	Name                       string  `json:"name,omitempty"`
	AppName                    *string `json:"app_name,omitempty"`
	SupportEmail               *string `json:"support_email,omitempty"`
	ConnectRequireKeyAuthTest  *bool   `json:"connect_require_key_auth_test,omitempty"`
}

type CreateProjectRequest struct {
	Name                      string  `json:"name"`
	AppName                   string  `json:"app_name,omitempty"`
	SupportEmail              string  `json:"support_email,omitempty"`
	ConnectRequireKeyAuthTest *bool   `json:"connect_require_key_auth_test,omitempty"`
}

type UpdateProjectRequest struct {
	Name                      *string `json:"name,omitempty"`
	AppName                   *string `json:"app_name,omitempty"`
	SupportEmail              *string `json:"support_email,omitempty"`
	ConnectRequireKeyAuthTest *bool   `json:"connect_require_key_auth_test,omitempty"`
}

type ListProjectsOptions struct {
	After  *string
	Before *string
	Limit  *int
	Q      string
}

type ListProjectsResponse struct {
	PageInfo PageInfo   `json:"page_info,omitzero"`
	Data     []*Project `json:"data,omitempty"`
}

type UpdateProjectLogoRequest struct {
	Logo string `json:"logo"`
}

type ProjectInfoApp struct {
	ID       *string `json:"id,omitempty"`
	NameSlug string  `json:"name_slug,omitempty"`
}

type ProjectInfoResponse struct {
	Apps []ProjectInfoApp `json:"apps,omitempty"`
}

// CreateProject creates a new Connect project.
func (c *Client) CreateProject(
	ctx context.Context,
	req *CreateProjectRequest,
) (*Project, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "projects"),
	})

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling create project request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating create project request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing create project request: %w", err)
	}

	var result Project
	if err := internal.UnmarshalResponse(response, &result, http.StatusCreated); err != nil {
		return nil, fmt.Errorf("unmarshalling create project response: %w", err)
	}

	return &result, nil
}

// ListProjects lists all Connect projects.
func (c *Client) ListProjects(
	ctx context.Context,
	opts *ListProjectsOptions,
) (*ListProjectsResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "projects"),
	})

	queryParams := url.Values{}
	if opts != nil {
		internal.AddQueryParams(queryParams, "q", opts.Q)
		if opts.After != nil {
			internal.AddQueryParams(queryParams, "after", *opts.After)
		}
		if opts.Before != nil {
			internal.AddQueryParams(queryParams, "before", *opts.Before)
		}
		internal.AddQueryParamInt(queryParams, "limit", opts.Limit)
	}
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating list projects request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing list projects request: %w", err)
	}

	var result ListProjectsResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling list projects response: %w", err)
	}

	return &result, nil
}

// GetProject retrieves a single project by ID.
func (c *Client) GetProject(
	ctx context.Context,
	projectID string,
) (*Project, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "projects", projectID),
	})

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating get project request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get project request: %w", err)
	}

	var result Project
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get project response: %w", err)
	}

	return &result, nil
}

// UpdateProject updates an existing project.
func (c *Client) UpdateProject(
	ctx context.Context,
	projectID string,
	req *UpdateProjectRequest,
) (*Project, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "projects", projectID),
	})

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling update project request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPatch, baseURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating update project request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing update project request: %w", err)
	}

	var result Project
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling update project response: %w", err)
	}

	return &result, nil
}

// DeleteProject deletes a project.
func (c *Client) DeleteProject(
	ctx context.Context,
	projectID string,
) error {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "projects", projectID),
	})

	req, err := http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating delete project request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return fmt.Errorf("executing delete project request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("expected status %d, got %d", http.StatusNoContent, response.StatusCode)
}

// UpdateProjectLogo uploads a new logo for a project.
// logoBase64 should be a data URI (e.g., "data:image/png;base64,AAAAAA...").
func (c *Client) UpdateProjectLogo(
	ctx context.Context,
	projectID string,
	logoBase64 string,
) error {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "projects", projectID, "logo"),
	})

	body, err := json.Marshal(UpdateProjectLogoRequest{Logo: logoBase64})
	if err != nil {
		return fmt.Errorf("marshalling update project logo request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating update project logo request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return fmt.Errorf("executing update project logo request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("expected 2xx status, got %d", response.StatusCode)
}

// GetProjectInfo retrieves project info including configured apps.
// Uses the client's configured project ID.
func (c *Client) GetProjectInfo(
	ctx context.Context,
) (*ProjectInfoResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "projects", "info"),
	})

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating get project info request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get project info request: %w", err)
	}

	var result ProjectInfoResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get project info response: %w", err)
	}

	return &result, nil
}
