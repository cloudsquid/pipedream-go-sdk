package pipedream

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type sourcesTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *sourcesTestSuite) SetupTest() {
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

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

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

func TestSources(t *testing.T) {
	suite.Run(t, new(sourcesTestSuite))
}
