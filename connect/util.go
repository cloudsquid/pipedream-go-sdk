package connect

import (
	"context"
	"fmt"
	"net/http"
)

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
