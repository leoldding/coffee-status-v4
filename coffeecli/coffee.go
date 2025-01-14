package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "setup",
				Usage: "write database name to env file",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// get path
					exePath, err := os.Executable()
					if err != nil {
						log.Fatal(err)
					}
					dirPath := filepath.Dir(exePath)

					// check that dynamodb table name exists
					input := cmd.Args().First()
					if input == "" {
						log.Fatal(errors.New("database name cannot be empty"))
					}

					dirPath = filepath.Join(dirPath, "../envs")
					if _, err := os.Stat(dirPath); os.IsNotExist(err) {
						// directory doesn't exist, create it
						err := os.Mkdir(dirPath, 0755)
						if err != nil {
							fmt.Println("Error creating directory:", err)
						}
					}
					filePath := filepath.Join(dirPath, "coffee.env")

					// write to env file
					err = os.WriteFile(filePath, []byte("DYNAMODB_TABLENAME="+input), 0666)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "set",
				Aliases: []string{"s"},
				Usage:   "set status value",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// get path
					exePath, err := os.Executable()
					if err != nil {
						log.Fatal(err)
					}
					dirPath := filepath.Dir(exePath)
					filePath := filepath.Join(dirPath, "../envs/coffee.env")
					// load table name
					err = godotenv.Load(filePath)
					if err != nil {
						log.Fatal(err)
					}
					tableName := os.Getenv("DYNAMODB_TABLENAME")

					// create dynamodb client
					cfg, err := config.LoadDefaultConfig(context.TODO())
					if err != nil {
						log.Fatal(err)
					}
					client := dynamodb.NewFromConfig(cfg)

					// check that new status exists
					input := cmd.Args().First()

					// check if new status is valid
					validStatuses := []string{"yes", "otw", "no"}
					if !slices.Contains(validStatuses, input) {
						log.Fatal(errors.New("invalid input"))
					}
					// put new status into table
					_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
						Item: map[string]types.AttributeValue{
							"key":   &types.AttributeValueMemberS{Value: "status"},
							"value": &types.AttributeValueMemberS{Value: input},
						},
						TableName: aws.String(tableName),
					})

					if err != nil {
						log.Fatal(err)
					}
					log.Println("Status set to", input)
					return nil
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "get current status",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// get path
					exePath, err := os.Executable()
					if err != nil {
						log.Fatal(err)
					}
					dirPath := filepath.Dir(exePath)
					filePath := filepath.Join(dirPath, "../envs/coffee.env")
					// load table name
					err = godotenv.Load(filePath)
					if err != nil {
						log.Fatal(err)
					}
					tableName := os.Getenv("DYNAMODB_TABLENAME")

					// create dynamodb client
					cfg, err := config.LoadDefaultConfig(context.TODO())
					if err != nil {
						log.Fatal(err)
					}
					client := dynamodb.NewFromConfig(cfg)
					cr := true // use to set ConsistentRead to true
					// get current status from dynamodb
					res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
						Key: map[string]types.AttributeValue{
							"key": &types.AttributeValueMemberS{Value: "status"},
						},
						TableName:      aws.String(tableName),
						ConsistentRead: &cr,
					})
					if err != nil {
						log.Fatal(err)
					}

					// unmarshal value into string
					var out string
					err = attributevalue.Unmarshal(res.Item["value"], &out)
					if err != nil {
						log.Fatal(err)
					}

					log.Println("Current Status: " + out)

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
