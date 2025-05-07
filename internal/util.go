package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
)

//func (c *pipedream.Client) doRequestViaApiKey(
//	ctx context.Context,
//	req *http.Request,
//) (*http.Response, error) {
//	req = req.WithContext(ctx)
//
//	req.Header.Set("Authorization", "Bearer "+c.apiKey)
//	req.Header.Set("X-PD-Environment", c.environment)
//	req.Header.Set("Content-Type", "application/json")
//
//	c.Logger.Info("Executing request",
//		"url", req.URL.String(),
//		"request", req.Header,
//		"environment", c.environment,
//	)
//
//	response, err := c.httpClient.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("request to pipedream api in environment %s failed: %w",
//			c.environment, err)
//	}
//
//	return response, nil
//}

func UnmarshalResponse(response *http.Response, result any, okStatusCodes ...int) error {
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

//func generateEndpointURL(postfix ...string) (string, error) {
//	baseURL, err := url.Parse(pipedreamApiURL)
//	if err != nil {
//		return "", fmt.Errorf("parsing url %s: %w", pipedreamApiURL, err)
//	}
//	baseURL.Path = path.Join(postfix...)
//
//	return baseURL.String(), nil
//}

func AddQueryParams(params url.Values, key, value string) {
	if value != "" {
		params.Add(key, value)
	}
}
