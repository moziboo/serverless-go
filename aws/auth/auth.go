package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.AuthorizationToken
	apiKeyExpected := os.Getenv("API_KEY")

	if token == apiKeyExpected {
		return generatePolicy("user", "Allow", "*"), nil
	} else {
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}
}

func generatePolicy(principalID string, effect string, resource string) events.APIGatewayCustomAuthorizerResponse {
	policyDocument := events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   effect,
				Resource: []string{resource},
			},
		},
	}

	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:    principalID,
		PolicyDocument: policyDocument,
	}
}

func main() {
	lambda.Start(handler)
}
