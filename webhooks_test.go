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

type webhooksTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *webhooksTestSuite) SetupTest() {
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

func (suite *webhooksTestSuite) TestCreateWebhook_Success() {
	require := suite.Require()
	expectedResponse := `{
		"data": {
			"id": "wh_abc123",
			"user_id": "u_abc123",
			"name": "My Webhook",
			"description": "Test webhook",
			"url": "https://webhook.site/abc",
			"active": true,
			"created_at": 1611964025,
			"updated_at": 1611964025
		}
	}}`
	expectedPath := "/webhooks"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodPost, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		require.Equal("https://webhook.site/abc", r.URL.Query().Get("url"))
		require.Equal("My Webhook", r.URL.Query().Get("name"))
		require.Equal("Test webhook", r.URL.Query().Get("description"))

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	resp, err := suite.pipedreamClient.CreateWebhook(
		context.Background(),
		"https://webhook.site/abc",
		"My Webhook",
		"Test webhook",
	)

	require.NoError(err)
	require.Equal("wh_abc123", resp.Data.ID)
	require.Equal("My Webhook", *resp.Data.Name)
}

func TestWebhooks(t *testing.T) {
	suite.Run(t, new(webhooksTestSuite))
}
