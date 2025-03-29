package services

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/leoldding/coffee-status-v4/internal/repository"
)

type CoffeeService struct {
	repo *repository.CoffeeRepository
}

func NewService(repo *repository.CoffeeRepository) *CoffeeService {
	return &CoffeeService{repo}
}

func (s *CoffeeService) GetStatus() (string, error) {
	res, err := s.repo.Get()

	if err != nil {
		return "", err
	}

	// unmarshal value into string
	var out string
	err = attributevalue.Unmarshal(res.Item["value"], &out)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (s *CoffeeService) UpdateStatus(status string) error {
	err := s.repo.Update(status)
	if err != nil {
		return err
	}
	return nil
}
