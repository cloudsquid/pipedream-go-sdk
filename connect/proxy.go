package connect

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
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
	Status     int         `json:"status"`
	Body       []byte      `json:"body"`
	Header     http.Header `json:"header"`
	StatusText string      `json:"status_text"`
}

// Proxy posts to: POST /v1/connect/{project_id}/proxy
// https://pipedream.com/docs/connect/api-proxy
func (c *Client) Proxy(
	ctx context.Context,
	pr ProxyRequest,
) (*ProxyResponse, error) {
	if err := pr.Validate(); err != nil {
		return nil, fmt.Errorf("proxy validation: %w", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString([]byte(pr.URL))
	proxyURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(
			c.ConnectURL().Path,
			c.ProjectID(),
			"proxy",
			encoded,
		),
	})
	q := proxyURL.Query()
	q.Set("external_user_id", pr.ExternalUserID)
	q.Set("account_id", pr.AccountID)
	proxyURL.RawQuery = q.Encode()

	body, err := c.prepareRequestBody(pr.Body)
	if err != nil {
		return nil, fmt.Errorf("prepare request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		pr.Method,
		proxyURL.String(),
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("create proxy request: %w", err)
	}

	req.Header = pr.Headers

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
		Status:     resp.StatusCode,
		Body:       b,
		Header:     resp.Header,
		StatusText: resp.Status,
	}, nil
}

func (pr *ProxyRequest) Validate() error {
	if strings.TrimSpace(pr.ExternalUserID) == "" {
		return fmt.Errorf("external_user_id is required")
	}
	if strings.TrimSpace(pr.AccountID) == "" {
		return fmt.Errorf("account_id is required")
	}
	if strings.TrimSpace(pr.URL) == "" {
		return fmt.Errorf("url is required")
	}
	if _, err := url.ParseRequestURI(pr.URL); err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	validMethods := []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"HEAD",
		"OPTIONS",
	}
	method := strings.ToUpper(pr.Method)
	for _, v := range validMethods {
		if method == v {
			return nil
		}
	}
	return fmt.Errorf("invalid method: %s", pr.Method)
}

func (c *Client) prepareRequestBody(
	bodyData json.RawMessage,
) (io.Reader, error) {
	if len(bodyData) == 0 {
		return nil, nil
	}
	if !json.Valid(bodyData) {
		return nil, fmt.Errorf("invalid JSON in request body")
	}

	return bytes.NewReader(bodyData), nil
}
