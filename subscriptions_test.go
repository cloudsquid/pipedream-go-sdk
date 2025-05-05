package pipedream

import (
	"context"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type subscriptionsTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *subscriptionsTestSuite) SetupTest() {
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

func (suite *subscriptionsTestSuite) TestSubscribeToEmitter_Success() {
	require := suite.Require()
	emitterID := "dc_abc123"
	listenerID := "p_xyz456"
	eventName := "$errors"

	expectedPath := "/subscriptions"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodPost, r.Method)
		require.Equal(emitterID, r.URL.Query().Get("emitter_id"))
		require.Equal(listenerID, r.URL.Query().Get("listener_id"))
		require.Equal(eventName, r.URL.Query().Get("event_name"))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	err = suite.pipedreamClient.SubscribeToEmitter(
		context.Background(),
		emitterID,
		listenerID,
		eventName,
	)
	require.NoError(err)
}

func (suite *subscriptionsTestSuite) TestAutoSubscribeToEvent_Success() {
	require := suite.Require()
	listenerID := "p_xyz456"
	eventName := "$errors"

	expectedPath := "/auto_subscriptions"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodPost, r.Method)
		require.Equal(eventName, r.URL.Query().Get("event_name"))
		require.Equal(listenerID, r.URL.Query().Get("listener_id"))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	err = suite.pipedreamClient.AutoSubscribeToEvent(
		context.Background(),
		eventName,
		listenerID,
	)
	require.NoError(err)
}

func (suite *subscriptionsTestSuite) TestDeleteSubscription_Success() {
	require := suite.Require()
	listenerID := "p_xyz456"
	emitterID := "dc_abc123"

	expectedPath := "/subscriptions"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(emitterID, r.URL.Query().Get("emitter_id"))
		require.Equal(listenerID, r.URL.Query().Get("listener_id"))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	err = suite.pipedreamClient.DeleteSubscription(
		context.Background(),
		emitterID,
		listenerID,
		"",
	)
	require.NoError(err)
}

func TestSubscriptions(t *testing.T) {
	suite.Run(t, new(subscriptionsTestSuite))
}
