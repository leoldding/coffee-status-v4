package cli

import (
	"log"
	"slices"

	"github.com/leoldding/coffee-status-v4/internal/services"
	"github.com/urfave/cli/v3"
)

type CoffeeCliHandler struct {
	service *services.CoffeeService
}

func NewHandler(service *services.CoffeeService) *CoffeeCliHandler {
	return &CoffeeCliHandler{service}
}

func (h *CoffeeCliHandler) HandleGetStatus() {
	res, err := h.service.GetStatus()
	if err != nil {
		log.Println("Error getting status from service:", err)
		return
	}

	log.Println("Current Status: " + res)
}

func (h *CoffeeCliHandler) HandleUpdateStatus(cmd *cli.Command) {
	if cmd.Args().Len() != 1 {
		log.Println("Command must contain one argument")
		return
	}

	// check that new status exists
	input := cmd.Args().First()

	// check if new status is valid
	validStatuses := []string{"yes", "otw", "no"}
	if !slices.Contains(validStatuses, input) {
		log.Println("Inputted status is invalid")
		return
	}

	h.service.UpdateStatus(input)

	log.Println("Status set to", input)
}
