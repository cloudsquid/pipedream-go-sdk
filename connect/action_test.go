package connect

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

type mockLogger struct{}

func (l *mockLogger) Debug(msg string, keyvals ...any) {}
func (l *mockLogger) Info(msg string, keyvals ...any)  {}
func (l *mockLogger) Warn(msg string, keyvals ...any)  {}
func (l *mockLogger) Error(msg string, keyvals ...any) {}

type actionTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *actionTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *actionTestSuite) TestInvokeAction_Success() {
	require := suite.Require()
	componentKey := "gitlab-new-issue"
	externalUserID := "jverce"
	configuredProp := ConfiguredProps{
		"gitlab": map[string]string{
			"authProvisionId": "apn_kVh9AoD",
		},
		"projectId": 45672541,
		"refName":   "main",
	}
	expectedResponse := `{
	  "exports": {
		"$summary": "Retrieved 1 commit"
	  },
	  "os": [],
	  "ret": [
		{
		  "id": "387262aea5d4a6920ac76c1e202bc9fd0841fea5",
		  "short_id": "387262ae",
		  "created_at": "2023-05-03T03:03:25.000+00:00",
		  "parent_ids": []
		}
	  ]
	}`
	expectedPath := "/project-abc/actions/run"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == oathPath:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{
				"access_token": "new-access-token",
				"expires_in": 3600
			}`)
			return
		case r.Method == http.MethodPost && r.URL.Path == expectedPath:
			require.Equal(http.MethodPost, r.Method)
			require.Equal(expectedPath, r.URL.Path)
			body, err := io.ReadAll(r.Body)
			require.NoError(err)

			var reqPayload InvokeActionRequest
			err = json.Unmarshal(body, &reqPayload)

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))

	defer server.Close()

	base := client.NewClient(&mockLogger{}, "", "project-abc", "development", "",
		"", nil, server.URL, "")
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.InvokeAction(
		context.Background(),
		componentKey,
		externalUserID,
		configuredProp,
	)

	require.NoError(err)
	require.EqualValues("Retrieved 1 commit", resp["exports"].(map[string]interface{})["$summary"].(string))
}

func TestAction(t *testing.T) {
	suite.Run(t, new(actionTestSuite))
}
