package dynamo

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Client is the DynamoDB client
var Client *dynamodb.Client

// InitClient loads AWS config and initializes DynamoDB
func InitClient() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	Client = dynamodb.NewFromConfig(cfg)
}
