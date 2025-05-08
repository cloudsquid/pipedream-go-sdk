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

type componentTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *componentTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *componentTestSuite) TestCreateComponent_Success() {
	require := suite.Require()

	expectedPath := "/components"

	expectedResponse := `{
		  "data": {
			"id": "sc_JDi8EB",
			"code": "component code here",
			"code_hash": "6",
			"name": "rss",
			"version": "0.0.1",
			"configurable_props": [
			  {
				"name": "url",
				"type": "string",
				"label": "Feed URL"
			  }
			],
			"created_at": 1588866900,
			"updated_at": 1588866900
		  }
		}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodPost, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqPayload CreateComponentRequest
		err = json.Unmarshal(body, &reqPayload)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.CreateComponent(
		context.Background(),
		"",
		"https://github.com/PipedreamHQ/pipedream/new-item-in-feed.ts",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.Data.Name, "rss")
}

func (suite *componentTestSuite) TestGetComponents_Success() {
	require := suite.Require()

	expectedPath := "/components/my-component"

	expectedResponse := `{
	  "data": {
		"id": "sc_JDi8EB",
		"code": "component code here",
		"code_hash": "685c7a680d055eaf505b08d5d814feef9fabd516d5960837d2e0838d3e1c9ed1",
		"name": "rss",
		"version": "0.0.1",
		"configurable_props": [
		  {
			"name": "url",
			"type": "string",
			"label": "Feed URL",
			"description": "Enter the URL for any public RSS feed."
		  }
		],
		"created_at": 1588866900,
		"updated_at": 1588866900
	  }
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetComponent(
		context.Background(),
		"my-component",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal("rss", resp.Data.Name)
}

func (suite *componentTestSuite) TestGetRegistryComponents_Success() {
	require := suite.Require()

	expectedPath := "/components/registry/github-new-repository"

	expectedResponse := `{
		  "data": {
			"id": "sc_JDi8EB",
			"code": "component code here",
			"code_hash": "6",
			"name": "rss",
			"version": "0.0.1",
			"configurable_props": [
			  {
				"name": "url",
				"type": "string",
				"label": "Feed URL"
			  }
			],
			"created_at": 1588866900,
			"updated_at": 1588866900
		  }
		}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetRegistryComponents(
		context.Background(),
		"github-new-repository",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.Data.Name, "rss")
}

func (suite *componentTestSuite) TestSearchRegistryComponents_Success() {
	require := suite.Require()

	expectedPath := "/components/search"

	expectedResponse := `{
			"sources": ["hubspot-new-contact"],
			"actions": ["twilio-send-sms"]
		}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(r.URL.RawQuery, "query=SendSMS")

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.SearchRegistryComponents(
		context.Background(),
		"SendSMS",
		"", 0, false,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.Actions, []string{"twilio-send-sms"})
	require.Equal(resp.Sources, []string{"hubspot-new-contact"})
}

func TestComponent(t *testing.T) {
	suite.Run(t, new(componentTestSuite))
}
