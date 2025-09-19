package connect

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

type ProxyRequest struct {
	ExternalUserID string          `json:"external_user_id"`
	AccountID      string          `json:"account_id"`
	Method         string          `json:"method"`
	URL            string          `json:"url"`
	Headers        http.Header     `json:"headers,omitempty"`
	Body           json.RawMessage `json:"body,omitempty"`
}

type ProxyResponse struct {
	Status int
	Body   []byte
	Header http.Header
}

// Proxy posts to: POST /v1/connect/{project_id}/proxy
// https://pipedream.com/docs/connect/api-proxy
func (c *Client) Proxy(
	ctx context.Context,
	pr ProxyRequest,
) (*ProxyResponse, error) {
	if pr.ExternalUserID == "" || pr.AccountID == "" ||
		pr.URL == "" {
		return nil, fmt.Errorf("proxy: missing required fields")
	}
	encoded := base64.RawURLEncoding.EncodeToString([]byte(pr.URL))

	base := c.ConnectURL()
	u := *base
	u.Path = path.Join(base.Path, c.ProjectID(), "proxy", encoded)

	q := u.Query()
	q.Set("external_user_id", pr.ExternalUserID)
	q.Set("account_id", pr.AccountID)
	u.RawQuery = q.Encode()

	var body io.Reader
	if len(pr.Body) > 0 {
		body = bytes.NewReader(pr.Body)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		pr.Method,
		u.String(),
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("create proxy request: %w", err)
	}

	req.Header = pr.Headers.Clone()

	resp, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("execute proxy request: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read proxy response: %w", err)
	}

	return &ProxyResponse{
		Status: resp.StatusCode,
		Body:   b,
		Header: resp.Header.Clone(),
	}, nil
}
