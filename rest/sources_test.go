package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/client"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type sourcesTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *sourcesTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *sourcesTestSuite) TestCreateSource_Success() {
	require := suite.Require()
	expectedResponse := `{
		"data": {
			"id": "dc_abc123",
			"user_id": "u_abc123",
			"component_id": "sc_abc123",
			"configured_props": {
				"url": "https://rss.m.pipedream.net",
				"timer": {
					"cron": null,
					"interval_seconds": 60
				}
			},
			"active": true,
			"created_at": 1589486978,
			"updated_at": 1589486978,
			"name": "your-name-here",
			"name_slug": "your-name-here"
		}
	}`
	expectedPath := "/sources"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodPost, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqPayload CreateSourceRequest
		err = json.Unmarshal(body, &reqPayload)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.CreateSource(
		context.Background(),
		"",
		"",
		"https://github.com/example/component.ts",
		"My Source",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal("your-name-here", resp.Data.Name)
}

func (suite *sourcesTestSuite) TestUpdateSource_Success() {
	require := suite.Require()
	expectedResponse := `{
		"data": {
			"id": "dc_abc123",
			"user_id": "u_abc123",
			"component_id": "sc_abc123",
			"configured_props": {
				"url": "https://rss.m.pipedream.net",
				"timer": {
					"cron": null,
					"interval_seconds": 60
				}
			},
			"active": true,
			"created_at": 1589486978,
			"updated_at": 1589486978,
			"name": "your-name-here",
			"name_slug": "your-name-here"
		}
	}`
	expectedPath := "/sources/dc_abc123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodPut, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqPayload UpdateSourceRequest
		err = json.Unmarshal(body, &reqPayload)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.UpdateSource(
		context.Background(),
		"dc_abc123",
		"",
		"",
		"https://github.com/example/component.ts",
		"My Source",
		true,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal("your-name-here", resp.Data.Name)
}

func (suite *sourcesTestSuite) TestDeleteSource_Success() {
	require := suite.Require()
	expectedPath := "/sources/dc_abc123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	err := suite.pipedreamClient.DeleteSource(
		context.Background(),
		"dc_abc123",
	)

	require.NoError(err)
}

func TestSources(t *testing.T) {
	suite.Run(t, new(sourcesTestSuite))
}
