package main

import (
	"context"
	"github.com/cloudsquid/pipedream-go-sdk"
	"log"
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
		"")

	accounts, err := sdk.Connect().ListAccounts(context.Background(), "org_1234", "slack", "", false)
	if err != nil {
		log.Fatalf("error listing accounts: %v", err)
	}

	for _, acc := range accounts.Data {
		log.Printf("Account: %s (%s)\n", acc.Name, acc.ID)
	}
}
