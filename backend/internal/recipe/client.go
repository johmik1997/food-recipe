package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/hasura/go-graphql-client"
)

type Client struct {
	gqlClient *graphql.Client
}

func NewClient(gqlClient *graphql.Client) *Client {
	return &Client{gqlClient: gqlClient}
}

func (c *Client) CreateRecipe(ctx context.Context, input CreateRecipeInput, userID string) (*Recipe, error) {
	var mutation struct {
		InsertRecipes struct {
			Returning []struct {
				ID string `graphql:"id"`
			} `graphql:"returning"`
		} `graphql:"insert_recipes(objects: $recipe)"`
	}

	// Calculate total time
	totalTime := 0
	if input.PrepTime != nil {
		totalTime += *input.PrepTime
	}
	if input.CookTime != nil {
		totalTime += *input.CookTime
	}

	recipeInput := map[string]interface{}{
		"title":            input.Title,
		"description":      input.Description,
		"prep_time":        input.PrepTime,
		"cook_time":        input.CookTime,
		"total_time":       totalTime,
		"servings":         input.Servings,
		"price":            input.Price,
		"user_id":          userID,
		"featured_image_url": nil,
	}

	err := c.gqlClient.Mutate(ctx, &mutation, map[string]interface{}{
		"recipe": recipeInput,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	if len(mutation.InsertRecipes.Returning) == 0 {
		return nil, errors.New("no recipe created")
	}

	recipeID := mutation.InsertRecipes.Returning[0].ID

	// Create categories
	if len(input.CategoryIDs) > 0 {
		var categoryMutation struct {
			InsertRecipeCategories struct {
				AffectedRows int `graphql:"affected_rows"`
			} `graphql:"insert_recipe_categories(objects: $categories)"`
		}

		categories := make([]map[string]interface{}, len(input.CategoryIDs))
		for i, catID := range input.CategoryIDs {
			categories[i] = map[string]interface{}{
				"recipe_id":   recipeID,
				"category_id": catID,
			}
		}

		err = c.gqlClient.Mutate(ctx, &categoryMutation, map[string]interface{}{
			"categories": categories,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create recipe categories: %w", err)
		}
	}

	// Create steps
	if len(input.Steps) > 0 {
		var stepMutation struct {
			InsertRecipeSteps struct {
				AffectedRows int `graphql:"affected_rows"`
			} `graphql:"insert_recipe_steps(objects: $steps)"`
		}

		steps := make([]map[string]interface{}, len(input.Steps))
		for i, step := range input.Steps {
			steps[i] = map[string]interface{}{
				"recipe_id":   recipeID,
				"step_number": step.StepNumber,
				"instruction": step.Instruction,
				"image_url":   step.ImageURL,
			}
		}

		err = c.gqlClient.Mutate(ctx, &stepMutation, map[string]interface{}{
			"steps": steps,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create recipe steps: %w", err)
		}
	}

	// Create ingredients
	if len(input.Ingredients) > 0 {
		var ingredientMutation struct {
			InsertRecipeIngredients struct {
				AffectedRows int `graphql:"affected_rows"`
			} `graphql:"insert_recipe_ingredients(objects: $ingredients)"`
		}

		ingredients := make([]map[string]interface{}, len(input.Ingredients))
		for i, ing := range input.Ingredients {
			ingredients[i] = map[string]interface{}{
				"recipe_id": recipeID,
				"name":      ing.Name,
				"quantity":  ing.Quantity,
				"unit":      ing.Unit,
			}
		}

		err = c.gqlClient.Mutate(ctx, &ingredientMutation, map[string]interface{}{
			"ingredients": ingredients,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create recipe ingredients: %w", err)
		}
	}

	return &Recipe{
		ID:               recipeID,
		Title:            input.Title,
		Description:      input.Description,
		PrepTime:         input.PrepTime,
		CookTime:         input.CookTime,
		TotalTime:        &totalTime,
		Servings:         input.Servings,
		Price:            input.Price,
		UserID:           userID,
		FeaturedImageURL: nil,
	}, nil
}