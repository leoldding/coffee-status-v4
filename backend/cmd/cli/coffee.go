package main

import (
	"context"
	"fmt"
	"log"
	"os"

	handlersCli "github.com/leoldding/coffee-status-v4/internal/cli"
	"github.com/leoldding/coffee-status-v4/internal/repository"
	"github.com/leoldding/coffee-status-v4/internal/services"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "setup",
				Usage: "write database name to env file",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return handlersCli.WriteEnvFile(ctx, cmd)
				},
			},
			{
				Name:    "set",
				Aliases: []string{"s"},
				Usage:   "set status value",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					repo, err := repository.NewCliRepository()
					if err != nil {
						log.Println("Error creating repo:", err)
						return fmt.Errorf("Error creating repo: %w", err)
					}
					service := services.NewService(repo)
					handler := handlersCli.NewHandler(service)

					handler.HandleUpdateStatus(cmd)

					return nil
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "get current status",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					repo, err := repository.NewCliRepository()
					if err != nil {
						log.Println("Error creating repo:", err)
						return fmt.Errorf("Error creating repo: %w", err)
					}
					service := services.NewService(repo)
					handler := handlersCli.NewHandler(service)

					handler.HandleGetStatus()

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
