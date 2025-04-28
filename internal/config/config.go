package config

import (
	"log"
	"os"
)

// Cfg holds environment configurations
var Cfg struct {
	SlackSigningSecret string
	SlackBotToken      string
	AWSRegion          string
	RequestsTable      string
	ProjectsTable      string
}

// Load reads and validates environment variables
func Load() {
	Cfg.SlackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
	Cfg.SlackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	Cfg.AWSRegion = os.Getenv("AWS_REGION")
	Cfg.RequestsTable = os.Getenv("DYNAMODB_TABLE_REQUESTS")
	Cfg.ProjectsTable = os.Getenv("DYNAMODB_TABLE_PROJECTS")

	if Cfg.SlackSigningSecret == "" || Cfg.SlackBotToken == "" {
		log.Fatal("Missing Slack credentials")
	}
	if Cfg.RequestsTable == "" || Cfg.ProjectsTable == "" {
		log.Fatal("Missing DynamoDB table names")
	}
}
