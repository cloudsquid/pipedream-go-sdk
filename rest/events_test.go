package rest

import (
	"context"
	"github.com/cloudsquid/pipedream-go-sdk/client"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type eventsTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *eventsTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *eventsTestSuite) TestGetSourceEvents_Success() {
	require := suite.Require()
	sourceID := "dc_test"
	limit := 21

	expectedPath := "/sources/dc_test/event_summaries"

	expectedResponse := `{
        "page_info": {
            "start_cursor": "1745858986889-0",
            "total_count": 13,
            "end_cursor": "1745858986853-0",
            "count": 3,
            "excluded_count": 0
        },
        "data": [
            {
                "id": "1745858986870-0",
                "indexed_at_ms": 1745858986870,
                "event": {
                    "rowNumber": 11
                },
                "metadata": {
                    "emitter_id": "dc_test",
                    "emit_id": "emit_123",
                    "name": "",
                    "summary": "New row #11",
                    "id": "0111",
                    "ts": 1745858986784
                }
            }
        ]
    }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodGet, r.Method)
		require.Equal(strconv.Itoa(limit), r.URL.Query().Get("limit"))

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetSourceEvents(
		context.Background(),
		sourceID,
		limit,
		false,
	)
	require.NoError(err)
	require.Equal(float64(11), resp.Data[0].Event["rowNumber"])
}

func (suite *eventsTestSuite) TestDeleteSourceEvents_Success() {
	require := suite.Require()
	sourceID := "dc_test"
	startID := "1745858986889-0"

	expectedPath := "/sources/dc_test/events"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodDelete, r.Method)
		require.Equal(startID, r.URL.Query().Get("start_id"))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	base := client.NewClient(&mockLogger{}, "dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	err := suite.pipedreamClient.DeleteSourceEvents(
		context.Background(),
		sourceID,
		startID,
		"",
	)
	require.NoError(err)
}

func TestEvents(t *testing.T) {
	suite.Run(t, new(eventsTestSuite))
}
