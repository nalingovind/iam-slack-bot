package main

import (
	"log"
	"net/http"

	"github.com/nalingovind/iam-slack-bot/internal/config"
	"github.com/nalingovind/iam-slack-bot/internal/dynamo"
	handlers "github.com/nalingovind/iam-slack-bot/internal/handlers"
	"github.com/nalingovind/iam-slack-bot/internal/lambda"
	"github.com/nalingovind/iam-slack-bot/internal/slack"
)

func main() {
	// Load configuration and clients
	config.Load()
	slack.InitClient()
	dynamo.InitClient()
	lambda.InitClient()

	// Register HTTP handlers
	http.HandleFunc("/workflows/aws-access", handlers.HandleWorkflow)
	http.HandleFunc("/slack/interactions", handlers.HandleInteraction)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
