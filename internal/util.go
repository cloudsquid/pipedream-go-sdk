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
	"strconv"
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

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func AddQueryParamInt(params url.Values, key string, value *int) {
	if value != nil {
		params.Add(key, strconv.Itoa(*value))
	}
}

func AddQueryParamBool(params url.Values, key string, value *bool) {
	if value != nil {
		params.Add(key, strconv.FormatBool(*value))
	}
}
