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

type componentTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *componentTestSuite) SetupTest() {
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

func (suite *componentTestSuite) TestPropOptions_Success() {
	require := suite.Require()
	propName := "projectId"
	componentKey := "gitlab-new-issue"
	externalUserID := "jverce"
	configuredProp := ConfiguredProps{
		"googleSheets": map[string]string{
			"authProvisionId": "apn_V1hMoE7",
		},
		"sheetId": "1BfWjFF2dTW",
	}

	expectedPath := "/project-abc/components/configure"

	expectedResponse := `{
      "observations": [],
	  "context": null,
	  "options": [
			{
			  "label": "jverce/foo-massive-10231-1",
			  "value": 45672541
			},
			{
			  "label": "jverce/foo-massive-10231",
			  "value": 45672514
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodPost, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqBody map[string]interface{}
		require.NoError(json.Unmarshal(body, &reqBody))

		require.Equal(componentKey, reqBody["id"])
		require.Equal(propName, reqBody["prop_name"])
		require.Equal(externalUserID, reqBody["external_user_id"])

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	resp, err := suite.pipedreamClient.GetPropOptions(
		context.Background(),
		propName,
		componentKey,
		externalUserID,
		configuredProp,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Len(resp.Options, 2)
	require.Equal(float64(45672514), resp.Options[1].Value)
}

func (suite *componentTestSuite) TestGetComponent_Success() {
	require := suite.Require()
	componentType := Components
	componentKey := "gitlab-new-issue"

	expectedPath := "/project-abc/components/gitlab-new-issue"

	expectedResponse := `{
		  "data": {
			"name": "New Issue (Instant)",
			"version": "0.1.2",
			"key": "gitlab-new-issue",
			"configurable_props": [
			  {
				"name": "gitlab",
				"type": "app",
				"app": "gitlab"
			  },
			  {
				"name": "db",
				"type": "$.service.db"
			  }
			]
		  }
		}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	resp, err := suite.pipedreamClient.GetComponent(
		context.Background(),
		componentKey,
		componentType,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Len(resp.Data.ConfigurableProps, 2)
	require.Equal(resp.Data.Name, "New Issue (Instant)")
}

func (suite *componentTestSuite) TestListComponents_Success() {
	require := suite.Require()
	componentType := Components
	appType := "gitlab"
	searchTerm := "issue"

	expectedPath := "/project-abc/components"
	expectedQuery := url.Values{
		"app": []string{"gitlab"},
		"q":   []string{"issue"},
	}

	expectedResponse := `{
	  "page_info": {
		"total_count": 1,
		"count": 1,
		"start_cursor": "c2NfM3ZpanpRcg",
		"end_cursor": "c2NfQjVpTkJBcA"
	  },
	  "data": [
		{
		  "name": "New Issue (Instant)",
		  "version": "0.1.2",
		  "key": "gitlab-new-issue"
		}
	  ]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(expectedQuery.Get("app"), r.URL.Query().Get("app"))
		require.Equal(expectedQuery.Get("q"), r.URL.Query().Get("q"))

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	resp, err := suite.pipedreamClient.ListComponents(
		context.Background(),
		componentType,
		appType,
		searchTerm,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Len(resp.Data, 1)
	require.Equal(resp.Data[0].Name, "New Issue (Instant)")
}

func (suite *componentTestSuite) TestReloadComponentProps_Success() {
	require := suite.Require()
	componentType := Components
	configuredProp := ConfiguredProps{
		"googleSheets": map[string]string{
			"authProvisionId": "apn_V1hMoE7",
		},
		"sheetId": "1BfWjFF2dTW",
	}
	externalUserID := "jay"
	componentID := "google-sheets"

	expectedResponse := `{
 		"observations": [],
		"errors": [],
		"dynamicProps": {
			"id": "dyp_6xUyVgQ",
			"configurableProps": [
				{
					"name": "googleSheets",
					"type": "app",
					"app": "google_sheets"
				}
			]
		}
	}`
	expectedPath := "/project-abc/components/props"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodPost, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqPayload ReloadComponentPropsRequest
		err = json.Unmarshal(body, &reqPayload)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	connectParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.connectURL = connectParsed

	resp, err := suite.pipedreamClient.ReloadComponentProps(
		context.Background(),
		componentType,
		configuredProp,
		externalUserID,
		componentID,
		"",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal("googleSheets", resp.DynamicProps.ConfigurableProps[0].Name)
}

func (suite *componentTestSuite) TestCreateComponent_Success() {
	require := suite.Require()

	expectedPath := "/components"

	expectedResponse := `{
	  "data": {
		"id": "sc_JDi8EB",
		"code": "component code here",
		"code_hash": "6",
		"name": "rss",
		"version": "0.0.1",
		"configurable_props": [
		  {
			"name": "url",
			"type": "string",
			"label": "Feed URL"
		  }
		],
		"created_at": 1588866900,
		"updated_at": 1588866900
	  }
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodPost, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqPayload CreateComponentRequest
		err = json.Unmarshal(body, &reqPayload)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	resp, err := suite.pipedreamClient.CreateComponent(
		context.Background(),
		"",
		"https://github.com/PipedreamHQ/pipedream/new-item-in-feed.ts",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.Data.Name, "rss")
}

func (suite *componentTestSuite) TestGetRegistryComponents_Success() {
	require := suite.Require()

	expectedPath := "/components/registry/github-new-repository"

	expectedResponse := `{
	  "data": {
		"id": "sc_JDi8EB",
		"code": "component code here",
		"code_hash": "6",
		"name": "rss",
		"version": "0.0.1",
		"configurable_props": [
		  {
			"name": "url",
			"type": "string",
			"label": "Feed URL"
		  }
		],
		"created_at": 1588866900,
		"updated_at": 1588866900
	  }
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	resp, err := suite.pipedreamClient.GetRegistryComponents(
		context.Background(),
		"github-new-repository",
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.Data.Name, "rss")
}

func (suite *componentTestSuite) TestSearchRegistryComponents_Success() {
	require := suite.Require()

	expectedPath := "/components/search"

	expectedResponse := `{
		"sources": ["hubspot-new-contact"],
		"actions": ["twilio-send-sms"]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(http.MethodGet, r.Method)
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(r.URL.RawQuery, "query=SendSMS")

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	restParsed, err := url.Parse(server.URL)
	require.NoError(err)

	suite.pipedreamClient.httpClient = server.Client()
	suite.pipedreamClient.baseURL = restParsed

	resp, err := suite.pipedreamClient.SearchRegistryComponents(
		context.Background(),
		"SendSMS",
		"", 0, false,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.Actions, []string{"twilio-send-sms"})
	require.Equal(resp.Sources, []string{"hubspot-new-contact"})
}

func TestComponent(t *testing.T) {
	suite.Run(t, new(componentTestSuite))
}
