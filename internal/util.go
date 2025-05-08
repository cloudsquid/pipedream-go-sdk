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

func AddQueryParams(params url.Values, key, value string) {
	if value != "" {
		params.Add(key, value)
	}
}
