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

type triggerTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *triggerTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *triggerTestSuite) TestDeployTrigger_Success() {
	require := suite.Require()
	componentKey := "gitlab-new-issue"
	externalUserID := "jay"
	configuredProp := ConfiguredProps{
		"gitlab": map[string]string{
			"authProvisionId": "apn_kVh9AoD",
		},
		"projectId": "1BfWjFF2dTW",
	}
	webhookURL := "https://events.example.com/gitlab-new-issue"
	expectedPath := "/project-abc/triggers/deploy"

	expectedResponse := `{
	  "data": {
		"id": "dc_dAuGmW7",
		"owner_id": "exu_oedidz",
		"component_id": "sc_3vijzQr",
		"configurable_props": [
		  {
			"name": "gitlab",
			"type": "app",
			"app": "gitlab"
		  }
		],
		"configured_props": {
		  "gitlab": {
			"authProvisionId": "apn_kVh9AoD"
		  },
		  "db": {
			"type": "$.service.db"
		  },
		  "http": {
			"endpoint_url": "https://xxxxxxxxxx.m.pipedream.net"
		  },
		  "projectId": 45672541
		},
		"active": true,
		"created_at": 1734028283,
		"updated_at": 1734028283
	  }
	}`

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
			body, err := io.ReadAll(r.Body)
			require.NoError(err)

			var reqBody DeployTriggerRequest
			err = json.Unmarshal(body, &reqBody)
			require.NoError(err)

			require.Equal(componentKey, reqBody.ComponentKey)
			require.Equal(externalUserID, reqBody.ExternalUserID)
			require.Equal(webhookURL, reqBody.WebhookURL)
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.DeployTrigger(
		context.Background(),
		componentKey,
		externalUserID,
		configuredProp,
		webhookURL,
		"",
		"",
	)

	require.NoError(err)
	require.NotNil(resp)
}

func (suite *triggerTestSuite) TestListDeployedTriggers_Success() {
	require := suite.Require()
	externalUserID := "jay"

	expectedPath := "/project-abc/deployed-triggers"

	expectedResponse := `{
	  "page_info": {
		"total_count": 1,
		"count": 1,
		"start_cursor": "ZGNfZ3p1bUsyZQ",
		"end_cursor": "ZGNfdjN1QllYZw"
	  },
	  "data": [
		{
		  "id": "dc_gzumK2e",
		  "owner_id": "exu_2LniLm",
		  "component_id": "sc_r1ixBpL",
		  "configurable_props": [
			{
			  "name": "googleDrive",
			  "type": "app",
			  "app": "google_drive"
			}
		  ],
		  "configured_props": {
			"googleDrive": {
			  "authProvisionId": "apn_V1hMeLM"
			}
		  },
		  "active": true,
		  "created_at": 1733512889,
		  "updated_at": 1733512889
		}
	  ]
	}`

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
			require.Equal(http.MethodGet, r.Method)
			require.Equal(expectedPath, r.URL.Path)
			require.Equal("jay", r.URL.Query().Get("external_user_id"))

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.ListDeployedTriggers(
		context.Background(),
		externalUserID,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Len(resp.Data, 1)
	require.Equal(resp.Data[0].ComponentID, "sc_r1ixBpL")
}

func (suite *triggerTestSuite) TestGetDeployedTrigger_Success() {
	require := suite.Require()
	externalUserID := "jay"
	deployedComponentID := "component_id"

	expectedPath := "/project-abc/deployed-triggers/component_id"

	expectedResponse := `{
	"data": {
		"id": "dc_gzumK2e",
		"owner_id": "exu_2LniLm",
		"component_id": "sc_r1ixBpL",
		"configurable_props": [
			{
			"name": "googleDrive",
			"type": "app",
			"app": "google_drive"
			}
		],
		"configured_props": {
			"googleDrive": {
			"authProvisionId": "apn_V1hMeLM"
			}
		},
		"active": true,
		"created_at": 1733512889,
		"updated_at": 1733512889,
		"name": "Danny Connect - exu_2LniLm",
		"name_slug": "danny-connect---exu-2-lni-lm-3"
	}
}`

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
			require.Equal(http.MethodGet, r.Method)
			require.Equal(expectedPath, r.URL.Path)
			require.Equal(externalUserID, r.URL.Query().Get("external_user_id"))

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetDeployedTrigger(
		context.Background(),
		deployedComponentID,
		externalUserID,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.ID, "dc_gzumK2e")
}

func (suite *triggerTestSuite) TestDeleteDeployedTrigger_Success() {
	require := suite.Require()
	externalUserID := "jay"
	deployedComponentID := "component_id"

	expectedPath := "/project-abc/deployed-triggers/component_id"

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
			require.Equal(http.MethodDelete, r.Method)
			require.Equal(expectedPath, r.URL.Path)
			require.Equal(externalUserID, r.URL.Query().Get("external_user_id"))

			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	err := suite.pipedreamClient.DeleteDeployedTrigger(
		context.Background(),
		deployedComponentID,
		externalUserID,
	)

	require.NoError(err)
}

func (suite *triggerTestSuite) TestRetrieveTriggerEvents_Success() {
	require := suite.Require()
	externalUserID := "jay"
	deployedComponentID := "component_id"

	expectedPath := "/project-abc/deployed-triggers/component_id/events"

	expectedResponse := `{
	  "data": [
		{
		  "e": {
			"method": "PUT",
			"path": "/",
			"query": [],
			"client_ip": "127.0.0.1",
			"url": "http://6c367a3dcffce4d01a7b691e906f8982.m.d.pipedream.net/",
			"headers": {
			  "host": "6c367a3dcffce4d01a7b691e906f8982.m.d.pipedream.net",
			  "connection": "close",
			  "user-agent": "curl/8.7.1",
			  "accept": "*/*"
			}
		  },
		  "k": "emit",
		  "ts": 1737155977519,
		  "id": "1737155977519-0"
		}
	  ]
	}`

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
			require.Equal(http.MethodGet, r.Method)
			require.Equal(expectedPath, r.URL.Path)
			require.Equal("jay", r.URL.Query().Get("external_user_id"))

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.RetrieveTriggerEvents(
		context.Background(),
		deployedComponentID,
		externalUserID,
		0,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Len(resp.Data, 1)
	require.Equal(resp.Data[0].E.Method, "PUT")
}

func (suite *triggerTestSuite) TestListTriggerWebhooks_Success() {
	require := suite.Require()
	externalUserID := "jay"
	deployedComponentID := "component_id"

	expectedPath := "/project-abc/deployed-triggers/component_id/webhooks"

	expectedResponse := `{
	  "webhook_urls": [
		"https://events.example.com/gitlab-new-issue"
	  ]
	}`

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
			require.Equal(http.MethodGet, r.Method)
			require.Equal(expectedPath, r.URL.Path)
			require.Equal("jay", r.URL.Query().Get("external_user_id"))

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.ListTriggerWebhooks(
		context.Background(),
		deployedComponentID,
		externalUserID,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.WebhookURLs[0], "https://events.example.com/gitlab-new-issue")
}

func (suite *triggerTestSuite) TestUpdateTriggerWebhooks_Success() {
	require := suite.Require()
	externalUserID := "jay"
	deployedComponentID := "component_id"
	webhookURLs := []string{"https://events.example.com/gitlab-new-issue"}

	expectedPath := "/project-abc/deployed-triggers/component_id/webhooks"

	expectedRequest := UpdateTriggerWebhooksRequest{
		ExternalUserID: "jay",
		WebhookURLs:    webhookURLs,
	}
	expectedResponse := `{
		"webhook_urls": ["https://events.example.com/gitlab-new-issue"]
	}`

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
			require.Equal(http.MethodPut, r.Method)
			require.Equal(expectedPath, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(err)

			var reqBody UpdateTriggerWebhooksRequest
			err = json.Unmarshal(body, &reqBody)
			require.NoError(err)

			require.Equal(expectedRequest.ExternalUserID, reqBody.ExternalUserID)
			require.Equal(expectedRequest.WebhookURLs, reqBody.WebhookURLs)

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.UpdateTriggerWebhooks(
		context.Background(),
		deployedComponentID,
		externalUserID,
		webhookURLs,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.WebhookURLs[0], "https://events.example.com/gitlab-new-issue")
}

func (suite *triggerTestSuite) TestRetrieveTriggerWorkflows_Success() {
	require := suite.Require()
	externalUserID := "jay"
	deployedComponentID := "component_id"

	expectedPath := "/project-abc/deployed-triggers/component_id/workflows"

	expectedResponse := `{
	  "workflow_ids": [
		"123"
	  ]
	}`

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
			require.Equal(http.MethodGet, r.Method)
			require.Equal(expectedPath, r.URL.Path)
			require.Equal("jay", r.URL.Query().Get("external_user_id"))

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.RetrieveTriggerWorkflows(
		context.Background(),
		deployedComponentID,
		externalUserID,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.WorkflowIDs[0], "123")
}

func (suite *triggerTestSuite) TestUpdateTriggerWorkflows_Success() {
	require := suite.Require()
	externalUserID := "jay"
	deployedComponentID := "component_id"
	workflowIDs := []string{"123"}

	expectedPath := "/project-abc/deployed-triggers/component_id/workflows"

	expectedRequest := UpdateTriggerWorkflowsRequest{
		ExternalUserID: "jay",
		WorkflowIDs:    workflowIDs,
	}
	expectedResponse := `{
		"workflow_ids": ["123"]
	}`

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
			require.Equal(http.MethodPut, r.Method)
			require.Equal(expectedPath, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(err)

			var reqBody UpdateTriggerWorkflowsRequest
			err = json.Unmarshal(body, &reqBody)
			require.NoError(err)

			require.Equal(expectedRequest.ExternalUserID, reqBody.ExternalUserID)
			require.Equal(expectedRequest.WorkflowIDs, reqBody.WorkflowIDs)

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, expectedResponse)
		}
	}))
	defer server.Close()

	base := client.NewClient("", "project-abc", "development", "",
		"", nil, server.URL, server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.UpdateTriggerWorkflows(
		context.Background(),
		deployedComponentID,
		externalUserID,
		workflowIDs,
	)

	require.NoError(err)
	require.NotNil(resp)
	require.Equal(resp.WorkflowIDs[0], "123")
}

func TestTrigger(t *testing.T) {
	suite.Run(t, new(triggerTestSuite))
}
