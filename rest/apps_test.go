package rest

import (
	"context"
	"fmt"
	"github.com/cloudsquid/pipedream-go-sdk/client"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type appsTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *appsTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *appsTestSuite) TestListApps_Success() {
	require := suite.Require()
	expectedPath := "/apps"
	expectedResponse := `{
	  "page_info": {
		"total_count": 1,
		"count": 1,
		"start_cursor": "c2xhY2s",
		"end_cursor": "c2xhY2tfYm90"
	  },
	  "data": [
		{
		  "id": "app_OkrhR1",
		  "name_slug": "slack",
		  "name": "Slack",
		  "auth_type": "oauth",
		  "description": "Slack is a channel-based messaging platform. With Slack, people can work together more effectively, connect all their software tools and services, and find the information they need to do their best work — all within a secure, enterprise-grade environment."
		}
	  ]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		require.Equal("git", r.URL.Query().Get("q"))

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.ListApps(
		context.Background(),
		"git",
		false,
		false,
		false,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.Data[0].Name, "Slack")
}

func (suite *appsTestSuite) TestGetApp_Success() {
	require := suite.Require()
	expectedPath := "/apps/app_OkrhR1"
	expectedResponse := `{
	  "data": 
		{
		  "id": "app_OkrhR1",
		  "name_slug": "slack",
		  "name": "Slack",
		  "auth_type": "oauth",
		  "description": "Slack is a channel-based messaging platform. With Slack, people can work together more effectively, connect all their software tools and services, and find the information they need to do their best work — all within a secure, enterprise-grade environment."
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetApp(
		context.Background(),
		"app_OkrhR1",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.Data.Name, "Slack")
}

func TestApps(t *testing.T) {
	suite.Run(t, new(appsTestSuite))
}
