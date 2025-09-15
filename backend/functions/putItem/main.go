package putItem

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/leoldding/coffee-status-v4/backend/internal/database"
	"github.com/leoldding/coffee-status-v4/backend/internal/models"
)

var tableName string

func init() {
	tableName = os.Getenv("DYNAMODB_TABLENAME")
	if tableName == "" {
		log.Fatal("DYNAMODB_TABLENAME must be set")
	}
	database.Init()
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var status models.Status

	if err := json.Unmarshal([]byte(req.Body), &status); err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	_, err := database.Client.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"key":   &types.AttributeValueMemberS{Value: "status"},
			"value": &types.AttributeValueMemberS{Value: status.Value},
		},
		TableName: aws.String(tableName),
	})

	if err != nil {
		log.Printf("Error putting status to DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
