package connect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudsquid/pipedream-go-sdk/client"
	"github.com/stretchr/testify/suite"
)

type authTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *authTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *authTestSuite) TestAcquireUserToken_Success() {
	require := suite.Require()

	expectedResponse := `{
		"connect_link_url": "randomURL.com",
		"expires_at": "2025-06-01T15:04:05Z",
		"token": "mock-user-token"
	}`

	expectedPath := "/project-abc/tokens"

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
			require.Equal(http.MethodPost, r.Method)
			require.Equal(expectedPath, r.URL.Path)

			var body UserTokenRequest
			raw, err := io.ReadAll(r.Body)
			require.NoError(err)

			err = json.Unmarshal(raw, &body)
			require.NoError(err)

			require.Equal("app_346", body.ExternalUserID)

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.AcquireUserToken(
		context.Background(),
		"app_346",
		"",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal("mock-user-token", resp.Token)
	require.Equal("randomURL.com", resp.ConnectLinkURL)
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(authTestSuite))
}
