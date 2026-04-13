package connect

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/cloudsquid/pipedream-go-sdk/internal"
)

type AppCategory struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ListAppsOptions struct {
	Q             string
	After         *string
	Before        *string
	Limit         *int
	SortKey       *string
	SortDirection *string
	CategoryIDs   []string
	HasComponents *bool
	HasActions    *bool
	HasTriggers   *bool
}

type ListAppsResponse struct {
	PageInfo PageInfo `json:"page_info,omitzero"`
	Data     []*App   `json:"data,omitempty"`
}

type GetAppResponse struct {
	Data *App `json:"data,omitempty"`
}

type ListAppCategoriesResponse []AppCategory

// ListApps lists available apps. Not scoped to a project.
func (c *Client) ListApps(
	ctx context.Context,
	opts *ListAppsOptions,
) (*ListAppsResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "apps"),
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
		if opts.SortKey != nil {
			internal.AddQueryParams(queryParams, "sort_key", *opts.SortKey)
		}
		if opts.SortDirection != nil {
			internal.AddQueryParams(queryParams, "sort_direction", *opts.SortDirection)
		}
		if len(opts.CategoryIDs) > 0 {
			internal.AddQueryParams(queryParams, "category_ids", strings.Join(opts.CategoryIDs, ","))
		}
		internal.AddQueryParamBool(queryParams, "has_components", opts.HasComponents)
		internal.AddQueryParamBool(queryParams, "has_actions", opts.HasActions)
		internal.AddQueryParamBool(queryParams, "has_triggers", opts.HasTriggers)
	}

	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating list apps request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing list apps request: %w", err)
	}

	var result ListAppsResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling list apps response: %w", err)
	}

	return &result, nil
}

// GetApp retrieves a single app by its name slug or ID.
func (c *Client) GetApp(
	ctx context.Context,
	appID string,
) (*GetAppResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "apps", appID),
	})

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating get app request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing get app request: %w", err)
	}

	var result GetAppResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling get app response: %w", err)
	}

	return &result, nil
}

// ListAppCategories retrieves the list of app categories.
func (c *Client) ListAppCategories(
	ctx context.Context,
) (ListAppCategoriesResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, "app_categories"),
	})

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating list app categories request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing list app categories request: %w", err)
	}

	var result ListAppCategoriesResponse
	if err := internal.UnmarshalResponse(response, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling list app categories response: %w", err)
	}

	return result, nil
}
