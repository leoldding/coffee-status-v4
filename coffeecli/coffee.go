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
	cli "github.com/urfave/cli/v2"
)

func main() {
	var setup bool
	var get bool

	app := &cli.App{
		Name:                   "coffee",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "setup",
				Aliases:     []string{"s"},
				Value:       false,
				Usage:       "write db name in env file",
				Destination: &setup,
			},
			&cli.BoolFlag{
				Name:        "get",
				Aliases:     []string{"g"},
				Value:       false,
				Usage:       "get current status",
				Destination: &get,
			},
		},
		Action: func(cCtx *cli.Context) error {
			// get path
			exePath, err := os.Executable()
			if err != nil {
				log.Fatal(err)
			}
			dirPath := filepath.Dir(exePath)

			if setup {
				// check that dynamodb table name exists
				if cCtx.Args().Len() != 1 {
					log.Fatal(errors.New("Must have one argument"))
				}
				input := cCtx.Args().Get(0)

				dirPath = filepath.Join(dirPath, "../envs")
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					// Directory doesn't exist, create it
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
			}

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

			if get {
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
			}

			// check that new status exists
			if cCtx.Args().Len() != 1 {
				log.Fatal(errors.New("Must have one argument"))
			}
			input := cCtx.Args().Get(0)

			// check if new status is valid
			validStatuses := []string{"yes", "otw", "no"}
			if !slices.Contains(validStatuses, input) {
				log.Fatal(errors.New("Invalid input"))
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
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
