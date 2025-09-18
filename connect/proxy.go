package connect

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

type ProxyRequest struct {
	ExternalUserID string            `json:"external_user_id"`
	AccountID      string            `json:"account_id"`
	Method         string            `json:"method"`
	URL            string            `json:"url"`
	Headers        map[string]string `json:"headers,omitempty"`
	Body           json.RawMessage   `json:"body,omitempty"`
}

type ProxyResponse struct {
	Status int
	Body   []byte
	Header http.Header
}

// Proxy posts to: POST /v1/connect/{project_id}/proxy
func (c *Client) Proxy(
	ctx context.Context,
	pr ProxyRequest,
) (*ProxyResponse, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "proxy"),
	})

	payload, err := json.Marshal(pr)
	if err != nil {
		return nil, fmt.Errorf("marshal proxy request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL.String(),
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("create proxy request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-pd-environment", c.Environment())

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("execute proxy request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read proxy response: %w", err)
	}

	return &ProxyResponse{
		Status: resp.StatusCode,
		Body:   body,
		Header: resp.Header.Clone(),
	}, nil
}
