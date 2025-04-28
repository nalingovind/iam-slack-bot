package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nalingovind/iam-slack-bot/internal/config"
)

// LookupProjectOwner returns the Slack ID of the owner for a project
type ProjectRecord struct {
	OwnerSlackID string `dynamodbav:"OwnerSlackID"`
}

func LookupProjectOwner(ctx context.Context, project string) (string, error) {
	out, err := Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &config.Cfg.ProjectsTable,
		Key: map[string]types.AttributeValue{
			"Project": &types.AttributeValueMemberS{Value: project},
		},
	})
	if err != nil {
		return "", err
	}
	if out.Item == nil {
		return "", fmt.Errorf("project not found")
	}
	owner := out.Item["OwnerSlackID"].(*types.AttributeValueMemberS).Value
	return owner, nil
}
