package pipedream

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type workspacesTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *workspacesTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.pipedreamClient = &Client{
		projectID:   "project-abc",
		environment: "development",
		token: &Token{
			AccessToken: "dummy-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
			CreatedAt:   int(time.Now().Unix()),
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		},
		logger: &mockLogger{},
		apiKey: "dummy-key",
	}
}

func (suite *workspacesTestSuite) TestGetWorkspaces_Success() {
	require := suite.Require()
	orgID := "o_Qa8I1Z"
	expectedResponse := `{
		"data": {
			"id": "o_Qa8I1Z",
			"orgname": "asdf",
			"name": "asdf",
			"email": "dev@pipedream.com",
			"daily_credits_quota": 100,
			"daily_credits_used": 0
		}
	}`

	expectedPath := "/workspaces/o_Qa8I1Z"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	resp, err := suite.pipedreamClient.GetWorkspace(
		context.Background(),
		orgID,
	)
	require.NoError(err)
	require.NotNil(resp)
	require.Equal("o_Qa8I1Z", resp.Data.ID)
}

func (suite *workspacesTestSuite) TestGetWorkspaceConnectedAccounts_Success() {
	require := suite.Require()
	orgID := "o_Qa8I1Z"
	expectedResponse := `{
	  "page_info": {
		"total_count": 1,
		"count": 1,
		"start_cursor": "YXBuXzJrVmhMUg",
		"end_cursor": "YXBuXzJrVmhMUg"
	  },
	  "data": [
		{
		  "id": "apn_2kVhLR",
		  "name": "Google Sheets #1"
		}
	  ]
	}`

	expectedPath := "/workspaces/o_Qa8I1Z/accounts"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	resp, err := suite.pipedreamClient.GetWorkspaceConnectedAccounts(
		context.Background(),
		orgID,
		"",
	)
	require.NoError(err)
	require.NotNil(resp)
	require.Equal("apn_2kVhLR", resp.Data[0].ID)
}

func (suite *workspacesTestSuite) TestGetWorkspaceSubscriptions_Success() {
	require := suite.Require()
	orgID := "o_Qa8I1Z"
	expectedResponse := `{
	  "data": [
		{
		  "id": "sub_abc123",
		  "emitter_id": "dc_abc123",
		  "listener_id": "p_abc123",
		  "event_name": ""
		},
		{
		  "id": "sub_def456",
		  "emitter_id": "dc_def456",
		  "listener_id": "p_def456",
		  "event_name": ""
		}
	  ]
	}`

	expectedPath := "/workspaces/o_Qa8I1Z/subscriptions"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	resp, err := suite.pipedreamClient.GetWorkspaceSubscriptions(
		context.Background(),
		orgID,
	)
	require.NoError(err)
	require.NotNil(resp)
	require.Equal("sub_def456", resp.Data[1].ID)
}

func (suite *workspacesTestSuite) TestGetWorkspaceSources_Success() {
	require := suite.Require()
	orgID := "o_Qa8I1Z"
	expectedResponse := `{
	  "page_info": {
		"total_count": 19,
		"count": 10,
		"start_cursor": "ZGNfSzB1QWVl",
		"end_cursor": "ZGNfeUx1alJx"
	  },
	  "data": [
		{
		  "id": "dc_abc123",
		  "component_id": "sc_def456",
		  "configured_props": {
			"http": {
			  "endpoint_url": "https://myendpoint.m.pipedream.net"
			}
		  },
		  "active": true,
		  "created_at": 1587679599,
		  "updated_at": 1587764467,
		  "name": "test",
		  "name_slug": "test"
		}
	  ]
	}`

	expectedPath := "/workspaces/o_Qa8I1Z/sources"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	resp, err := suite.pipedreamClient.GetWorkspaceSources(
		context.Background(),
		orgID,
	)
	require.NoError(err)
	require.NotNil(resp)
	require.Equal("dc_abc123", resp.Data[0].ID)
}

func TestWorkspaces(t *testing.T) {
	suite.Run(t, new(workspacesTestSuite))
}
