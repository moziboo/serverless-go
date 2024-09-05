package utils

import (
	"github.com/aws/aws-lambda-go/events"
)

func HandleCORS(req events.APIGatewayProxyRequest) map[string]string {
	headers := make(map[string]string)

	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Credentials"] = "true"
	headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization, Api-Key"
	headers["Access-Control-Allow-Methods"] = "GET, OPTIONS"
	headers["Access-Control-Allow-Origin"] = "*"

	return headers
}
