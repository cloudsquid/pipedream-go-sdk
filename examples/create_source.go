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
		"production",         // Environment: "production" or "development"
		"your-client-id",     // OAuth Client ID
		"your-client-secret", // OAuth Client Secret
		[]string{},           // Allowed Origins
		"",                   // Connect API URL (optional, defaults to public)
		"")                   // Rest API URL (optional, defaults to public)

	source, err := sdk.Rest().CreateSource(
		context.Background(),
		"sc_abc123",
		"", "",
		"example-source",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created source: %s", source.Data.ID)
}
