package dynamo

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/nalingovind/iam-slack-bot/internal/config"
	"github.com/nalingovind/iam-slack-bot/internal/models"
)

// CreateRequest writes a new PENDING request entry and returns its ID
func CreateRequest(ctx context.Context, req models.WorkflowRequest, ownerID string) (string, error) {
	id := uuid.New().String()
	duration, _ := time.ParseDuration(req.Duration)
	expiresAt := time.Now().Add(duration).UTC().Format(time.RFC3339)

	item := map[string]types.AttributeValue{
		"RequestID":     &types.AttributeValueMemberS{Value: id},
		"UserID":        &types.AttributeValueMemberS{Value: req.UserID},
		"UserName":      &types.AttributeValueMemberS{Value: req.UserName},
		"Project":       &types.AttributeValueMemberS{Value: req.Project},
		"Role":          &types.AttributeValueMemberS{Value: req.Role},
		"Justification": &types.AttributeValueMemberS{Value: req.Justification},
		"Status":        &types.AttributeValueMemberS{Value: "PENDING"},
		"ExpiresAt":     &types.AttributeValueMemberS{Value: expiresAt},
		"CreatedAt":     &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
		"OwnerSlackID":  &types.AttributeValueMemberS{Value: ownerID},
	}
	_, err := Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &config.Cfg.RequestsTable,
		Item:      item,
	})
	return id, err
}

// UpdateRequestStatus updates the Status of a request
func UpdateRequestStatus(ctx context.Context, id, status string) error {
	_, err := Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &config.Cfg.RequestsTable,
		Key: map[string]types.AttributeValue{
			"RequestID": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          awsString("SET #s = :status"),
		ExpressionAttributeNames:  map[string]string{"#s": "Status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":status": &types.AttributeValueMemberS{Value: status}},
	})
	return err
}

func awsString(s string) *string {
	return &s
}
