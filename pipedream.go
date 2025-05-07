package pipedream

import (
	"github.com/cloudsquid/pipedream-go-sdk/client"
	"github.com/cloudsquid/pipedream-go-sdk/connect"
	"github.com/cloudsquid/pipedream-go-sdk/rest"
)

type SDK struct {
	connect *connect.Client
	rest    *rest.Client
}

func NewPipedreamClient(
	logger client.Logger,
	apiKey string,
	projectID string,
	environment string,
	clientID string,
	clientSecret string,
	allowedOrigins []string,
	connectURL string,
	restURL string,
) *SDK {

	pd := client.NewClient(logger, apiKey, projectID, environment, clientID, clientSecret, allowedOrigins, connectURL, restURL)

	return &SDK{
		connect: &connect.Client{Client: pd},
		rest:    &rest.Client{Client: pd},
	}
}

func (sdk *SDK) Connect() *connect.Client { return sdk.connect }

func (sdk *SDK) Rest() *rest.Client { return sdk.rest }
