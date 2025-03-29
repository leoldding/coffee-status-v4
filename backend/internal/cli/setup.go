package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func WriteEnvFile(ctx context.Context, cmd *cli.Command) error {
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
	err = os.WriteFile(filePath, []byte("COFFEE_DYNAMODB_TABLENAME="+input), 0666)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
