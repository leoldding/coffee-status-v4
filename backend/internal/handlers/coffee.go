package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/leoldding/coffee/internal/awsclients"
)

type CoffeeHandler struct {
	clients *awsclients.Clients
}

type status struct {
	Value string `json:"value"`
}

func NewCoffeeHandler(clients *awsclients.Clients) *CoffeeHandler {
	return &CoffeeHandler{clients}
}

func (c *CoffeeHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	tableName, err := c.clients.SSM.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/leoding/dynamodb/coffee"),
	})
	if err != nil {
		log.Println("ssm error:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	cr := true
	resp, err := c.clients.DynamoDB.GetItem(context.Background(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: "status"},
		},
		TableName:      aws.String(*tableName.Parameter.Value),
		ConsistentRead: &cr,
	})
	if err != nil {
		log.Println("dynamodb get item error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var out string
	_ = attributevalue.Unmarshal(resp.Item["value"], &out)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status{Value: out})
	return
}

func (c *CoffeeHandler) PutStatus(w http.ResponseWriter, r *http.Request) {
	var stat status

	err := json.NewDecoder(r.Body).Decode(&stat)
	if err != nil {
		log.Println("failed to decode json body:", err)
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	tableName, err := c.clients.SSM.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/leoding/dynamodb/coffee"),
	})
	if err != nil {
		log.Println("ssm error:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	_, err = c.clients.DynamoDB.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"key":   &types.AttributeValueMemberS{Value: "status"},
			"value": &types.AttributeValueMemberS{Value: stat.Value},
		},
		TableName: aws.String(*tableName.Parameter.Value),
	})

	if err != nil {
		log.Printf("dynamodb put item error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})

	return
}
