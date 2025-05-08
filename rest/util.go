package rest

import (
	"context"
	"fmt"
	"net/http"
)

func (p *Client) doRequestViaApiKey(
	ctx context.Context,
	req *http.Request,
) (*http.Response, error) {
	req = req.WithContext(ctx)

	req.Header.Set("Authorization", "Bearer "+p.APIKey())
	req.Header.Set("X-PD-Environment", p.Environment())
	req.Header.Set("Content-Type", "application/json")

	p.Logger.Info("Executing request",
		"url", req.URL.String(),
		"request", req.Header,
		"environment", p.Environment(),
	)

	response, err := p.HTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to pipedream api in environment %s failed: %w",
			p.Environment(), err)
	}

	return response, nil
}

func (c *Client) doRequestViaOauth(
	ctx context.Context,
	req *http.Request,
) (*http.Response, error) {
	req = req.WithContext(ctx)

	err := c.AcquireAccessToken()
	if err != nil {
		return nil, fmt.Errorf("acquiring access token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token().AccessToken)
	req.Header.Set("X-PD-Environment", c.Environment())
	req.Header.Set("Content-Type", "application/json")

	c.Logger.Info("Executing request",
		"url", req.URL.String(),
		"request", req.Header,
		"environment", c.Environment(),
	)

	response, err := c.HTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to pipedream api in environment %s failed: %w",
			c.Environment(), err)
	}

	return response, nil
}
