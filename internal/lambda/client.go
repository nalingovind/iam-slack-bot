package lambda

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// Client is the Lambda API client
var Client *lambda.Client

// InitClient initializes the Lambda client
func InitClient() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	Client = lambda.NewFromConfig(cfg)
}
