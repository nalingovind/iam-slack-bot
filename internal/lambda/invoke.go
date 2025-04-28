package lambda

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// InvokeProvisioningLambda triggers the specified Lambda function
func InvokeProvisioningLambda(ctx context.Context, functionName string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = Client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: &functionName,
		Payload:      body,
	})
	return err
}
