package slack

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

// Client is the Slack API client
var Client *slack.Client

// InitClient initializes the Slack client
func InitClient() {
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		log.Fatal("SLACK_BOT_TOKEN is required")
	}
	Client = slack.New(token)
}
