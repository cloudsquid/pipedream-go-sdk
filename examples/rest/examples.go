package main

import (
	"context"
	"log"

	"github.com/cloudsquid/pipedream-go-sdk"
)

func main() {
	sdk := pipedream.NewPipedreamClient(
		"your-api-key",
		"your-project-id",
		"development",        // Environment: "production" or "development"
		"your-client-d",      // OAuth Client ID
		"your-client-secret", // OAuth Client Secret
		[]string{},           // Allowed Origins
		"",                   // Connect API URL (optional, defaults to public)
		"")                   // Rest API URL (optional, defaults to public)

	components, err := sdk.Rest().GetRegistryComponents(
		context.Background(),
		"github-new-repository",
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("recieved global registry components: %v", components.Data)

	events, err := sdk.Rest().GetSourceEvents(
		context.Background(),
		"p_2gCYljl",
		10,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("recieved source events: %v", events.Data)
}
