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

type actionTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *actionTestSuite) SetupTest() {
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
	}
}

func (suite *accountsTestSuite) TestInvokeAction_Success() {
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
		require.Equal(http.MethodPost, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqPayload InvokeActionRequest
		err = json.Unmarshal(body, &reqPayload)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

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
