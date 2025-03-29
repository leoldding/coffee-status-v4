package repository

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joho/godotenv"
)

type CoffeeRepository struct {
	database  *dynamodb.Client
	tableName string
}

func NewCliRepository() (*CoffeeRepository, error) {
	exePath, err := os.Executable()
	if err != nil {
		log.Println("Unable to get executable path:", err)
		return nil, err
	}
	dirPath := filepath.Dir(exePath)
	filePath := filepath.Join(dirPath, "../envs/coffee.env")
	// load table name
	err = godotenv.Load(filePath)
	if err != nil {
		log.Println("Unable to load .env file:", err)
		return nil, err
	}

	tableName := os.Getenv("COFFEE_DYNAMODB_TABLENAME")

	if tableName == "" {
		err := errors.New("COFFEE_DYNAMODB_TABLENAME is not set in env file")
		log.Println(err)
		return nil, err
	}
	return newRepository(tableName)
}

func NewHttpRepository() (*CoffeeRepository, error) {
	tableName := os.Getenv("COFFEE_DYNAMODB_TABLENAME")

	if tableName == "" {
		err := errors.New("COFFEE_DYNAMODB_TABLENAME is not set in environment variables")
		log.Println(err)
		return nil, err
	}

	return newRepository(tableName)
}

func newRepository(tableName string) (*CoffeeRepository, error) {
	// create dynamodb client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Error loading default config:", err)
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)

	return &CoffeeRepository{client, tableName}, nil
}

func (r *CoffeeRepository) Get() (*dynamodb.GetItemOutput, error) {
	cr := true // use to set ConsistentRead to true

	// get current status from dynamodb
	res, err := r.database.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: "status"},
		},
		TableName:      aws.String(r.tableName),
		ConsistentRead: &cr,
	})
	if err != nil {
		log.Println("Error retrieving status from DynamoDB:", err)
		return nil, err
	}

	return res, nil
}

func (r *CoffeeRepository) Update(status string) error {
	_, err := r.database.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"key":   &types.AttributeValueMemberS{Value: "status"},
			"value": &types.AttributeValueMemberS{Value: status},
		},
		TableName: aws.String(r.tableName),
	})

	if err != nil {
		log.Println("Error updating status in DynamoDB:", err)
		return err
	}

	return nil
}
