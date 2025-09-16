package getItem

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
	cr := true
	res, err := database.Client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: "status"},
		},
		TableName:      aws.String(tableName),
		ConsistentRead: &cr,
	})
	if err != nil {
		log.Printf("Error retrieving status from DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"failed to retrieve item"}`}, err
	}
	var out string
	_ = attributevalue.Unmarshal(res.Item["value"], &out)
	body, _ := json.Marshal(models.Status{Value: out})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}
