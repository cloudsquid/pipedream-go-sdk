package main

import (
	"context"
	"log"

	"github.com/cloudsquid/pipedream-go-sdk"
	"github.com/cloudsquid/pipedream-go-sdk/connect"
	"github.com/cloudsquid/pipedream-go-sdk/internal"
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

	accounts, err := sdk.Connect().ListAccounts(context.Background(), &connect.ListAccountsOptions{
		ExternalUserID:     "org_1234",
		App:                internal.StringPtr("slack"),
		IncludeCredentials: false,
	})
	if err != nil {
		log.Fatalf("error listing accounts: %v", err)
	}

	for _, acc := range accounts.Data {
		log.Printf("Account: %s (%s)\n", acc.Name, acc.ID)
	}
}
