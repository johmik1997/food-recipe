package recipe

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hasura/go-graphql-client"
)

type (
	CreateRecipeInput struct {
		Title            string                 `json:"title"`
		Description      *string                `json:"description"`
		PrepTime         *int                   `json:"prep_time"`
		CookTime         *int                   `json:"cook_time"`
		Servings         *int                   `json:"servings"`
		FeaturedImageURL *string                `json:"featured_image_url"`
		Price            *int                   `json:"price"`
		UserID           *string                `json:"user_id"`
		CategoryIDs      []string               `json:"category_ids"`
		Steps            []RecipeStepInput      `json:"steps"`
		Ingredients      []RecipeIngredientInput `json:"ingredients"`
		Images           []RecipeImageInput      `json:"images"`
	}

	RecipeStepInput struct {
		StepNumber  int     `json:"step_number"`
		Instruction string  `json:"instruction"`
		ImageURL    *string `json:"image_url"`
	}

	RecipeIngredientInput struct {
		Name     string   `json:"name"`
		Quantity *float64 `json:"quantity"`
		Unit     *string  `json:"unit"`
	}

	RecipeImageInput struct {
		ImageURL   string `json:"image_url"`
		IsFeatured *bool  `json:"is_featured"`
	}

	RecipeOutput struct {
		ID               string  `json:"id"`
		Title            string  `json:"title"`
		Description      *string `json:"description"`
		PrepTime         *int    `json:"prep_time"`
		CookTime         *int    `json:"cook_time"`
		TotalTime        *int    `json:"total_time"`
		Servings         *int    `json:"servings"`
		FeaturedImageURL *string `json:"featured_image_url"`
		Price            *int    `json:"price"`
		CreatedAt        string  `json:"created_at"`
	}
)

type RecipeClient struct {
	client *graphql.Client
}

func NewRecipeClient(graphqlURL, JWTSecret string) *RecipeClient {
	return &RecipeClient{
		client: graphql.NewClient(graphqlURL, nil).
			WithRequestModifier(func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+JWTSecret)
				r.Header.Set("Content-Type", "application/json")

			}),
	}
}

func (c *RecipeClient) CreateRecipe(ctx context.Context, input CreateRecipeInput) (*RecipeOutput, error) {
	var mutation struct {
		CreateRecipe RecipeOutput `graphql:"createRecipe(object: $object)"`
	}

	vars := map[string]interface{}{
		"object": input,
	}

	err := c.client.Mutate(ctx, &mutation, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	// Calculate total time
	totalTime := 0
	if input.PrepTime != nil {
		totalTime += *input.PrepTime
	}
	if input.CookTime != nil {
		totalTime += *input.CookTime
	}
	mutation.CreateRecipe.TotalTime = &totalTime

	return &mutation.CreateRecipe, nil
}