package recipe

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func (s *Service) CreateRecipe(ctx context.Context, input CreateRecipeInput, userID string) (*Recipe, error) {
	// Validate input
	if input.Title == "" {
		return nil, errors.New("title is required")
	}
	if len(input.Steps) == 0 {
		return nil, errors.New("at least one step is required")
	}
	if len(input.Ingredients) == 0 {
		return nil, errors.New("at least one ingredient is required")
	}
	if len(input.CategoryIDs) == 0 {
		return nil, errors.New("at least one category is required")
	}
	
	// Validate category IDs
	for _, catID := range input.CategoryIDs {
		if _, err := uuid.Parse(catID); err != nil {
			return nil, errors.New("invalid category ID format")
		}
	}

	return s.client.CreateRecipe(ctx, input, userID)
}