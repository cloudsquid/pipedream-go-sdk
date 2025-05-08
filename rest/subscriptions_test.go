package rest

import (
	"context"
	"github.com/cloudsquid/pipedream-go-sdk/client"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type subscriptionsTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *subscriptionsTestSuite) SetupTest() {
	suite.ctx = context.Background()
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

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	err := suite.pipedreamClient.SubscribeToEmitter(
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

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	err := suite.pipedreamClient.AutoSubscribeToEvent(
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

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	err := suite.pipedreamClient.DeleteSubscription(
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
