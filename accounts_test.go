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

type mockLogger struct{}

func (l *mockLogger) Debug(msg string, keyvals ...any) {}
func (l *mockLogger) Info(msg string, keyvals ...any)  {}
func (l *mockLogger) Warn(msg string, keyvals ...any)  {}
func (l *mockLogger) Error(msg string, keyvals ...any) {}

type accountsTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *accountsTestSuite) SetupTest() {
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

func (suite *accountsTestSuite) TestListAccounts_Success() {
	require := suite.Require()
	expectedResponse := `{
		 "page_info": {
			"total_count": 1,
			"count": 1,
			"start_cursor": "YXBuX0JtaEJKSm0",
			"end_cursor": "YXBuX1YxaE1lTE0"
		},
		"data": [
			{
				"id": "apn_XehyZPr",
				"name": "name",
				"external_id": "user-123",
				"healthy": true,
				"dead": false,
				"app": {
				  "id": "app_OkrhR1",
				  "name": "github"
				},
				"created_at": "2024-07-30T22:52:48.000Z",
				"updated_at": "2024-08-01T03:44:17.000Z"
      		}
		  ]
	}`
	expectedPath := "/project-abc/accounts"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		require.Equal("github", r.URL.Query().Get("app"))
		require.Equal("user-123", r.URL.Query().Get("external_user_id"))
		require.Equal("oauth-789", r.URL.Query().Get("oauth_app_id"))
		require.Equal("true", r.URL.Query().Get("include_credentials"))
		require.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	resp, err := suite.pipedreamClient.ListAccounts(
		context.Background(),
		"user-123",
		"github",
		"oauth-789",
		true,
	)

	require.NoError(err)
	require.Equal(1, len(resp.Data))
	require.Equal("apn_XehyZPr", resp.Data[0].ID)
}

func (suite *accountsTestSuite) TestListAccounts_RefreshAccessToken_Success() {
	require := suite.Require()
	expectedResponse := `{
		 "page_info": {
			"total_count": 1,
			"count": 1,
			"start_cursor": "YXBuX0JtaEJKSm0",
			"end_cursor": "YXBuX1YxaE1lTE0"
		},
		"data": [
			{
				"id": "apn_XehyZPr",
				"name": "name",
				"external_id": "user-123",
				"healthy": true,
				"dead": false,
				"app": {
				  "id": "app_OkrhR1",
				  "name": "github"
				},
				"created_at": "2024-07-30T22:52:48.000Z",
				"updated_at": "2024-08-01T03:44:17.000Z"
      		}
		  ]
	}`
	expectedPath := "/project-abc/accounts"
	expectedTokenPath := "/oauth/token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == expectedTokenPath:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{
				"access_token": "new-access-token",
				"expires_in": 3600
			}`)
			return
		case r.Method == http.MethodGet && r.URL.Path == expectedPath:
			require.Equal(http.MethodGet, r.Method)
			require.Equal(expectedPath, r.URL.Path)

			require.Equal("github", r.URL.Query().Get("app"))
			require.Equal("user-123", r.URL.Query().Get("external_user_id"))
			require.Equal("oauth-789", r.URL.Query().Get("oauth_app_id"))
			require.Equal("true", r.URL.Query().Get("include_credentials"))
			require.Equal("application/json", r.Header.Get("Content-Type"))

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed
	suite.pipedreamClient.baseURL = connectParsed
	suite.pipedreamClient.token.ExpiresAt = time.Now()

	resp, err := suite.pipedreamClient.ListAccounts(
		context.Background(),
		"user-123",
		"github",
		"oauth-789",
		true,
	)

	require.NoError(err)
	require.Equal(1, len(resp.Data))
	require.Equal("apn_XehyZPr", resp.Data[0].ID)
}

func (suite *accountsTestSuite) TestGetAccount_Success() {
	require := suite.Require()
	expectedResponse := `{
		"data": {
				"id": "apn_XehyZPr",
				"name": "shaghayegh",
				"external_id": "user-123",
				"healthy": true,
				"dead": false,
				"app": {
				  "id": "app_OkrhR1",
				  "name": "github"
				},
				"created_at": "2024-07-30T22:52:48.000Z",
				"updated_at": "2024-08-01T03:44:17.000Z"
		}
	}`
	expectedPath := "/project-abc/accounts/apn_XehyZPr"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		require.Equal("github", r.URL.Query().Get("app"))
		require.Equal("user-123", r.URL.Query().Get("external_user_id"))
		require.Equal("false", r.URL.Query().Get("include_credentials"))

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	resp, err := suite.pipedreamClient.GetAccount(
		context.Background(),
		"user-123",
		"github",
		false,
		"apn_XehyZPr",
	)

	require.NoError(err)
	require.Equal("shaghayegh", resp.Data.Name)
}

func (suite *accountsTestSuite) TestGetAccount_Failure() {
	require := suite.Require()
	expectedPath := "/project-abc/accounts/apn_XehyZPr"
	expectedResponse := `{"error": "record not found"}`
	expectedError := fmt.Errorf("unexpected status code %d: %s", http.StatusNotFound, expectedResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		require.Equal("github", r.URL.Query().Get("app"))
		require.Equal("user-123", r.URL.Query().Get("external_user_id"))
		require.Equal("false", r.URL.Query().Get("include_credentials"))

		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	resp, err := suite.pipedreamClient.GetAccount(
		context.Background(),
		"user-123",
		"github",
		false,
		"apn_XehyZPr",
	)

	require.EqualError(err, expectedError.Error())
	require.Nil(resp)
}

func (suite *accountsTestSuite) TestDeleteAccount_Success() {
	require := suite.Require()
	expectedPath := "/project-abc/accounts/apn_XehyZPr"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	err = suite.pipedreamClient.DeleteAccount(
		context.Background(),
		"apn_XehyZPr",
	)

	require.NoError(err)
}

func (suite *accountsTestSuite) TestDeleteAccount_Failure() {
	require := suite.Require()
	expectedPath := "/project-abc/accounts/apn_XehyZPr"
	expectedError := fmt.Errorf("expected status %d, got %d", http.StatusNoContent, http.StatusNotFound)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	original := pipedreamApiURL
	pipedreamApiURL = server.URL
	defer func() { pipedreamApiURL = original }()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed
	err = suite.pipedreamClient.DeleteAccount(
		context.Background(),
		"apn_XehyZPr",
	)

	require.EqualError(err, expectedError.Error())
}

func (suite *accountsTestSuite) TestDeleteAccounts_Success() {
	require := suite.Require()
	expectedPath := "/project-abc/apps/app_346/accounts"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	err = suite.pipedreamClient.DeleteAccounts(
		context.Background(),
		"app_346",
	)

	require.NoError(err)
}

func (suite *accountsTestSuite) TestDeleteAccounts_Failure() {
	require := suite.Require()
	expectedPath := "/project-abc/apps/app_346/accounts"
	expectedError := fmt.Errorf("expected status %d, got %d", http.StatusNoContent, http.StatusNotFound)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	err = suite.pipedreamClient.DeleteAccounts(
		context.Background(),
		"app_346",
	)

	require.EqualError(err, expectedError.Error())
}

func (suite *accountsTestSuite) TestEndUser_Success() {
	require := suite.Require()
	expectedPath := "/project-abc/users/user-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	err = suite.pipedreamClient.DeleteEndUser(
		context.Background(),
		"user-123",
	)

	require.NoError(err)
}

func (suite *accountsTestSuite) TestEndUser_Failure() {
	require := suite.Require()
	expectedPath := "/project-abc/users/user-123"
	expectedError := fmt.Errorf("expected status %d, got %d", http.StatusNoContent, http.StatusNotFound)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	err = suite.pipedreamClient.DeleteEndUser(
		context.Background(),
		"user-123",
	)

	require.EqualError(err, expectedError.Error())
}

func TestAccounts(t *testing.T) {
	suite.Run(t, new(accountsTestSuite))
}
