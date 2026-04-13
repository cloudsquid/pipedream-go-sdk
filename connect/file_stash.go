package connect

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/cloudsquid/pipedream-go-sdk/internal"
)

// DownloadFile downloads a file from the file stash.
// The caller is responsible for closing the returned ReadCloser.
func (c *Client) DownloadFile(
	ctx context.Context,
	s3Key string,
) (io.ReadCloser, error) {
	baseURL := c.ConnectURL().ResolveReference(&url.URL{
		Path: path.Join(c.ConnectURL().Path, c.ProjectID(), "file_stash", "download"),
	})

	queryParams := url.Values{}
	internal.AddQueryParams(queryParams, "s3_key", s3Key)
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating download file request: %w", err)
	}

	response, err := c.doRequestViaOauth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing download file request: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(body))
	}

	return response.Body, nil
}
