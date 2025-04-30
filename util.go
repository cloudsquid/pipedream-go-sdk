package pipedream

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
	"slices"
)

func (p *Client) doRequest(
	ctx context.Context,
	req *http.Request,
) (*http.Response, error) {
	req = req.WithContext(ctx)

	err := p.acquireAccessToken()
	if err != nil {
		return nil, fmt.Errorf("acquiring access token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.token.AccessToken)
	req.Header.Set("X-PD-Environment", p.environment)
	req.Header.Set("Content-Type", "application/json")

	p.logger.Info("Executing request",
		"url", req.URL.String(),
		"request", req.Header,
		"environment", p.environment,
	)

	response, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to pipedream api in environment %s failed: %w",
			p.environment, err)
	}

	return response, nil
}

func unmarshalResponse(response *http.Response, result any, okStatusCodes ...int) error {
	if len(okStatusCodes) == 0 {
		okStatusCodes = append(okStatusCodes, http.StatusOK)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		err = fmt.Errorf("decoding response body: %w", err)
	}

	if !slices.Contains(okStatusCodes, response.StatusCode) {
		compact := new(bytes.Buffer)
		if json.Compact(compact, body) != nil {
			compact.WriteString("[non-JSON body]")
		}

		statusErr := fmt.Errorf("unexpected status code: %d, body: %s",
			response.StatusCode, compact.String())
		if err != nil {
			return errors.Join(statusErr, err)
		}
		return statusErr
	}

	return err
}

func generateEndpointURL(postfix ...string) (string, error) {
	baseURL, err := url.Parse(pipedreamApiURL)
	if err != nil {
		return "", fmt.Errorf("parsing url %s: %w", pipedreamApiURL, err)
	}
	baseURL.Path = path.Join(postfix...)

	return baseURL.String(), nil
}

func addQueryParams(params url.Values, key, value string) {
	if value != "" {
		params.Add(key, value)
	}
}
