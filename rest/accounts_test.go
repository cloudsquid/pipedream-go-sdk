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

type accountsTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *accountsTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

const oathPath = "/oauth/token"

func (suite *accountsTestSuite) TestListAccounts_Success() {
	require := suite.Require()
	expectedResponse := `{
	  "data": [
		{
		  "id": "apn_abc123",
		  "created_at": "2022-07-27T20:37:52.000Z",
		  "updated_at": "2024-02-11T04:18:46.000Z",
		  "name": "Google Sheets — pipedream.com", 
		  "app": {
			"id": "app_abc123",
			"name": "Google Sheets"
		  },
		  "healthy": true 
		}
	  ]
	}`
	expectedPath := "/accounts"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == oathPath:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{
				"access_token": "new-access-token",
				"expires_in": 3600
			}`)
			return
		case r.URL.Path == expectedPath:

			require.Equal(http.MethodGet, r.Method)
			require.Equal(expectedPath, r.URL.Path)

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.ListAccounts(
		context.Background(),
		"",
		"",
		false,
	)

	require.NoError(err)
	require.Equal(1, len(resp.Data))
	require.Equal("apn_abc123", resp.Data[0].ID)
}

func (suite *accountsTestSuite) TestGetAccount_Success() {
	require := suite.Require()
	expectedResponse := `{
	  "data": {
		"id": "apn_abc123",
		"created_at": "2022-07-27T20:37:52.000Z",
		"updated_at": "2024-02-11T04:18:46.000Z",
		"expires_at": "2024-02-11T05:18:46.000Z",
		"name": "Google Sheets — pipedream.com",
		"app": {
		  "id": "app_abc123",
		  "name": "Google Sheets"
		},
		"credentials": {
		}
	  }
	}`
	expectedPath := "/accounts/user-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == oathPath:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{
				"access_token": "new-access-token",
				"expires_in": 3600
			}`)
			return
		case r.URL.Path == expectedPath:
			require.Equal(http.MethodGet, r.Method)
			require.Equal(expectedPath, r.URL.Path)
			require.Equal("true", r.URL.Query().Get("include_credentials"))

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetAccount(
		context.Background(),
		"user-123",
		true,
	)

	require.NoError(err)
	require.Equal("apn_abc123", resp.Data.ID)
}

func TestAccounts(t *testing.T) {
	suite.Run(t, new(accountsTestSuite))
}
