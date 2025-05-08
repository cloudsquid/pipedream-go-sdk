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

type usersTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *usersTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *usersTestSuite) TestGetCurrentUser_Success() {
	require := suite.Require()
	expectedResponse := `{
		"data": {
			"id": "u_abc123",
			"username": "dyburger",
			"email": "dylan@pipedream.com",
			"daily_compute_time_quota": 95400000,
			"daily_compute_time_used": 8420300,
			"daily_invocations_quota": 27344,
			"daily_invocations_used": 24903,
			"orgs": [
				{
					"name": "MyWorkspace",
					"id": "o_abc123",
					"orgname": "myworkspace",
					"email": "workspace@pipedream.com",
					"daily_credits_quota": 100,
					"daily_credits_used": 0
				}
			]
		}
	}`
	expectedPath := "/users/me"

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

	resp, err := suite.pipedreamClient.GetCurrentUser(
		context.Background())

	require.NoError(err)
	require.Equal("dyburger", resp.Data.Username)
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(usersTestSuite))
}
