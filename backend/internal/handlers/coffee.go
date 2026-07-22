package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/leoldding/coffee/internal/awsclients"
)

type CoffeeHandler struct {
	clients   *awsclients.Clients
	tableName string
}

type status struct {
	Value string `json:"value"`
}

func NewCoffeeHandler(clients *awsclients.Clients, tableName string) *CoffeeHandler {
	return &CoffeeHandler{
		clients:   clients,
		tableName: tableName,
	}
}

func (c *CoffeeHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := c.clients.DynamoDB.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: "status"},
		},
		TableName:      aws.String(c.tableName),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		log.Printf("dynamodb get item error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var out string
	if resp.Item == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	val, ok := resp.Item["value"]
	if !ok {
		http.Error(w, "invalid item", http.StatusInternalServerError)
		return
	}

	if err := attributevalue.Unmarshal(val, &out); err != nil {
		log.Printf("unmarshal error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(status{Value: out}); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func (c *CoffeeHandler) PutStatus(w http.ResponseWriter, r *http.Request) {
    log.Println("PutStatus handler reached")

	defer r.Body.Close()

	ctx := r.Context()
	var stat status

	err := json.NewDecoder(r.Body).Decode(&stat)
	if err != nil {
		log.Printf("failed to decode json body: %v", err)
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	_, err = c.clients.DynamoDB.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: "status"},
		},
		TableName: aws.String(c.tableName),

		UpdateExpression:    aws.String("SET #v = :val"),
		ConditionExpression: aws.String("attribute_exists(#k)"),

		ExpressionAttributeNames: map[string]string{
			"#k": "key",
			"#v": "value",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val": &types.AttributeValueMemberS{Value: stat.Value},
		},
	})
	if err != nil {
		log.Printf("dynamodb update item error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{"success": true}); err != nil {
		log.Printf("encode error: %v", err)
	}
}
