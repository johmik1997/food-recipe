package recipe

import (
	"context"
	"foodrecipe/internal/utils"
)

type RecipeService struct {
	recipeClient RecipeClient
}

func NewRecipeService(recipeClient RecipeClient) *RecipeService {
	return &RecipeService{
		recipeClient: recipeClient,
	}
}

func (s *RecipeService) CreateRecipe(ctx context.Context, input CreateRecipeInput) (RecipeOutput, error) {
	// Validate required fields
	if input.Title == "" {
		return RecipeOutput{}, utils.NewValidationError("title is required")
	}
	if len(input.Steps) == 0 {
		return RecipeOutput{}, utils.NewValidationError("at least one step is required")
	}
	if len(input.Ingredients) == 0 {
		return RecipeOutput{}, utils.NewValidationError("at least one ingredient is required")
	}
	if input.UserID == nil {
		return RecipeOutput{}, utils.NewValidationError("user_id is required")
	}

	// Create recipe via GraphQL
	result, err := s.recipeClient.CreateRecipe(ctx, input)
	if err != nil {
		return RecipeOutput{}, err
	}
	return *result, nil
}