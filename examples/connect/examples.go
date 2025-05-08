package main

import (
	"context"
	"github.com/cloudsquid/pipedream-go-sdk"
	"log"
)

type StdLogger struct{}

func (l *StdLogger) Debug(msg string, keyvals ...any) { log.Println("[DEBUG]", msg, keyvals) }
func (l *StdLogger) Info(msg string, keyvals ...any)  { log.Println("[INFO]", msg, keyvals) }
func (l *StdLogger) Warn(msg string, keyvals ...any)  { log.Println("[WARN]", msg, keyvals) }
func (l *StdLogger) Error(msg string, keyvals ...any) { log.Println("[ERROR]", msg, keyvals) }

func main() {
	sdk := pipedream.NewPipedreamClient(
		&StdLogger{},
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
		log.Printf("\nAccount: %s (%s)", acc.Name, acc.ID)
	}
}
