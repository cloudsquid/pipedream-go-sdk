package rest

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

type workflowsTestSuite struct {
	suite.Suite
	ctx             context.Context
	pipedreamClient *Client
}

func (suite *workflowsTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *workflowsTestSuite) TestCreateWorkflow_Success() {
	require := suite.Require()
	orgID := "org123"
	projectID := "proj123"
	templateID := "tmpl456"
	steps := []WorkflowStep{
		{
			Namespace: "code",
			Props: map[string]interface{}{
				"stringProp": "asdf",
				"intProp":    5,
			},
		},
	}

	triggers := []WorkflowTrigger{
		{
			Props: map[string]interface{}{
				"oauth":  map[string]string{"authProvisionId": "apn_123"},
				"string": "hello",
			},
		},
	}
	settings := &WorkflowSettings{Name: "example workflow name", AutoDeploy: true}

	expectedResponse := `{
	  "data": {
	    "id": "p_48rCxZ",
	    "name": "example workflow name",
	    "active": true,
	    "steps": [
	      {
	        "id": "c_bDf10L",
	        "type": "CodeCell",
	        "namespace": "code",
	        "disabled": false,
	        "code_raw": null,
	        "codeRaw": null,
	        "codeConfigJson": null,
	        "lang": "nodejs20.x",
	        "text_raw": null,
	        "appConnections": [],
	        "flat_params_visibility_json": null,
	        "params_json": "{}",
	        "component": true,
	        "savedComponent": {
	          "id": "sc_PRYiAZ",
	          "code": "component-code",
	          "codeHash": "hash",
	          "configurableProps": [
	            { "name": "stringProp", "type": "string" },
	            { "name": "intProp", "type": "integer" }
	          ],
	          "key": null,
	          "description": null,
	          "entryPath": null,
	          "version": "",
	          "apps": []
	        },
	        "component_key": null,
	        "component_owner_id": null,
	        "configured_props_json": "{\"intProp\":5,\"stringProp\":\"asdf\"}",
	        "authProvisionIdMap": {},
	        "authProvisionIds": []
	      }
	    ],
	    "triggers": [
	      {
	        "id": "dc_rmXuv3",
	        "owner_id": "o_BYDI5y",
	        "component_id": "sc_PgliBJ",
	        "configured_props": {},
	        "active": true,
	        "created_at": 1707241571,
	        "updated_at": 1707241571,
	        "name": "Emit hello world",
	        "name_slug": "emit-hello-world-6"
	      }
	    ]
	  }
	}`

	expectedPath := "/workflows"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqPayload CreateWorkflowRequest
		err = json.Unmarshal(body, &reqPayload)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.CreateWorkflow(
		context.Background(),
		orgID,
		projectID,
		templateID,
		steps,
		triggers,
		settings,
	)
	require.NoError(err)
	require.NotNil(resp)
	require.Equal("p_48rCxZ", resp.Data.ID)
	require.Equal("example workflow name", resp.Data.Name)
}

func (suite *workflowsTestSuite) TestUpdateWorkflow_Success() {
	require := suite.Require()
	workflowID := "p_48rCxZ"
	orgID := "org123"
	active := false
	expectedResponse := `{
		  "data": {
			"inactive": true,
			"name_slug": "test-http-trigger",
			"id": "p_48rCxZ",
			"endpoint_id": "en8745sd2vo1fo5",
			"owner_id": "o_JvIwWMD",
			"owner_type": "Org",
			"name": "test http trigger",
			"description": null,
			"created_at": "2025-04-15T13:44:28.000Z",
			"updated_at": "2025-05-02T12:05:30.000Z",
			"emits": [
			  {
				"e": {
				  "orgId": "o_JvIwWMD",
				  "email": "someemial@example.com",
				  "subscriptionActive": false,
				  "useCredits": true,
				  "isDev": false,
				  "devNamespace": null
				},
				"k": "emit",
				"ts": 1745858287801,
				"id": "1745858287801-0"
			  }
			],
			"emitter_connected": null,
			"project_id": 468082,
			"route_params": {
			  "ownerName": "name",
			  "id": "p_yKCwOWDz",
			  "nameSlug": "test-http-trigger-",
			  "projectId": "proj_zNswc1XRz"
			},
			"edit": true,
			"deployments": [
			  [
				"d_v7sjwdcrm1K8",
				"2025-04-15T13:44:28.000Z"
			  ]
			]
		  }
		}`

	expectedPath := "/workflows/p_48rCxZ"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodPut, r.Method)

		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		var reqPayload UpdateWorkflowRequest
		err = json.Unmarshal(body, &reqPayload)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.UpdateWorkflow(
		context.Background(),
		workflowID,
		orgID,
		active,
	)
	require.NoError(err)
	require.NotNil(resp)
	data := (*resp)["data"].(map[string]any)
	require.Equal("p_48rCxZ", data["id"])
}

func (suite *workflowsTestSuite) TestGetWorkflowDetails_Success() {
	require := suite.Require()
	workflowID := "p_48rCxZ"
	orgID := "org123"
	expectedResponse := `{
	  "triggers": [
		{
		  "id": "hi_ABpHKz",
		  "key": "eabcdefghiklmnop",
		  "endpoint_url": "http://eabcdefghiklmnop.m.d.pipedream.net",
		  "custom_response": false
		}
	  ],
	  "steps": [
		{
		  "id": "c_abc123",
		  "namespace": "code",
		  "disabled": false,
		  "lang": "nodejs20.x",
		  "appConnections": [],
		  "component": true,
		  "savedComponent": {
			"id": "sc_abc123",
			"codeHash": "long-hash-here",
			"configurableProps": [],
			"version": ""
		  },
		  "component_key": null,
		  "component_owner_id": "o_abc123",
		  "configured_props_json": "{}"
		}
	  ]
	}`

	expectedPath := "/workflows/p_48rCxZ"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodGet, r.Method)
		require.Equal(orgID, r.URL.Query().Get("org_id"))

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetWorkflowDetails(
		context.Background(),
		workflowID,
		orgID,
	)
	require.NoError(err)
	require.NotNil(resp)
	require.Equal("hi_ABpHKz", resp.Triggers[0].ID)
}

func (suite *workflowsTestSuite) TestGetWorkflowEmits_Success() {
	require := suite.Require()
	workflowID := "p_48rCxZ"
	orgID := "org123"
	expectedResponse := `{
	  "page_info": {
		"total_count": 1,
		"start_cursor": "1606511826306-0",
		"end_cursor": "1606511826306-0",
		"count": 1
	  },
	  "data": [
		{
		  "id": "1606511826306-0",
		  "indexed_at_ms": 1606511826306,
		  "event": {
			"raw_event": {
			  "name": "Luke",
			  "title": "Jedi"
			}
		  },
		  "metadata": {
			"emit_id": "1ktF96gAMsLqdYSRWYL9KFS5QqW",
			"name": "",
			"emitter_id": "p_abc123"
		  }
		}
	  ]
	}`

	expectedPath := "/workflows/p_48rCxZ/event_summaries"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodGet, r.Method)
		require.Equal(orgID, r.URL.Query().Get("org_id"))

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetWorkflowEmits(
		context.Background(),
		workflowID,
		orgID,
		false,
		0,
	)
	require.NoError(err)
	require.NotNil(resp)
	require.Equal("1606511826306-0", resp.Data[0].ID)
}

func (suite *workflowsTestSuite) TestGetWorkflowErrors_Success() {
	require := suite.Require()
	workflowID := "p_48rCxZ"
	expectedResponse := `{
	  "page_info": {
		"total_count": 100,
		"start_cursor": "1606370816223-0",
		"end_cursor": "1606370816223-0",
		"count": 1
	  },
	  "data": [
		{
		  "id": "1606370816223-0",
		  "indexed_at_ms": 1606370816223,
		  "event": {
			"original_event": {
			  "name": "Luke",
			  "title": "Jedi"
			},
			"original_context": {
			  "id": "1kodJIW7jVnKfvB2yp1OoPrtbFk",
			  "ts": "2020-11-26T06:06:44.652Z",
			  "workflow_id": "p_abc123",
			  "deployment_id": "d_abc123",
			  "source_type": "SDK"
			},
			"error": {
			  "code": "InternalFailure",
			  "cellId": "c_abc123",
			  "ts": "2020-11-26T06:06:56.077Z",
			  "stack": "    at Request.extractError ..."
			},
			"metadata": {}
		  }
		}
	  ]
	}`

	expectedPath := "/workflows/p_48rCxZ/$errors/event_summaries"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(expectedPath, r.URL.Path)
		require.Equal(http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	base := client.NewClient("dummy-key", "project-abc", "development", "",
		"", nil, "", server.URL)
	suite.pipedreamClient = &Client{Client: base}

	resp, err := suite.pipedreamClient.GetWorkflowErrors(
		context.Background(),
		workflowID,
		false,
		0,
	)
	require.NoError(err)
	require.NotNil(resp)
	require.Equal("1606370816223-0", resp.Data[0].ID)
}

func TestWorkflows(t *testing.T) {
	suite.Run(t, new(workflowsTestSuite))
}
