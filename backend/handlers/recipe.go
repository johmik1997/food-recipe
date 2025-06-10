package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foodrecipe/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var hasuraURL = utils.GraphQLURL    // e.g., http://hasura:8080/v1/graphql
var hasuraSecret = utils.AdminSecret // your Hasura admin secret

type CreateRecipeInput struct {
	Title            string                  `json:"title"`
	Description      *string                 `json:"description"`
	PrepTime         *int                    `json:"prep_time"`
	CookTime         *int                    `json:"cook_time"`
	Servings         *int                    `json:"servings"`
	FeaturedImageURL *string                 `json:"featured_image_url"`
	UserID           *string                 `json:"user_id"`
	Price            *int               `json:"price"`
	CategoryIDs      []string                `json:"category_ids"`
	Steps            []RecipeStepInput       `json:"steps"`
	Ingredients      []RecipeIngredientInput `json:"ingredients"`
	Images           []RecipeImageInput      `json:"images"`
}

type RecipeStepInput struct {
	StepNumber  int     `json:"step_number"`
	Instruction string  `json:"instruction"`
	ImageURL    *string `json:"image_url"`
}

type RecipeIngredientInput struct {
	Name     string   `json:"name"`
	Quantity *float64 `json:"quantity"`
	Unit     *string  `json:"unit"`
}

type RecipeImageInput struct {
	ImageURL   string `json:"image_url"`
	IsFeatured *bool  `json:"is_featured"`
}

type HasuraActionPayload struct {
	Input struct {
		Object CreateRecipeInput `json:"object"`
	} `json:"input"`
}

type RecipeOutput struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Description      *string `json:"description"`
	PrepTime         *int    `json:"prep_time"`
	CookTime         *int    `json:"cook_time"`
	TotalTime        *int    `json:"total_time"`
	Servings         *int    `json:"servings"`
	FeaturedImageURL *string `json:"featured_image_url"`
	Price		   *int `json:"price"`
	CreatedAt        string  `json:"created_at"`
}

func CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received create recipe request")

	// Read and log raw body
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	log.Printf("Raw request body: %s", string(bodyBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // reset body for decoding

	// Decode JSON
	var payload HasuraActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Error decoding payload: %v", err)
		http.Error(w, `{"message": "invalid request"}`, http.StatusBadRequest)
		return
	}

	input := payload.Input.Object
	recipeID := uuid.New().String()
	totalTime := 0
	if input.PrepTime != nil {
		totalTime += *input.PrepTime
	}
	if input.CookTime != nil {
		totalTime += *input.CookTime
	}

	// Build GraphQL mutation
mutation := fmt.Sprintf(`
mutation {
  insert_recipes_one(object: {
	id: "%s",
	title: "%s",
	description: %s,
	prep_time: %s,
	cook_time: %s,
	servings: %s,
	featured_image_url: %s,
	price: %s,
	user_id: %s,
  }) {
	id
	created_at
  }
  %s
  %s
  %s
  %s
}`,
		recipeID,
		escapeString(input.Title),
		nullOrStr(input.Description),
		nullOrInt(input.PrepTime),
		nullOrInt(input.CookTime),
		nullOrInt(input.Servings),
		nullOrStr(input.FeaturedImageURL),
		nullOrInt(input.Price),
		nullOrStr(input.UserID),
		buildCategoryInsert(recipeID,input.CategoryIDs),
		buildStepsInsert(recipeID, input.Steps),
		buildIngredientsInsert(recipeID, input.Ingredients),
		buildImagesInsert(recipeID, input.Images),
	)

	log.Printf("Executing mutation:\n%s", mutation)
	result, err := doHasuraRequest(mutation)
	if err != nil {
		log.Printf("Hasura error: %v", err)
		http.Error(w, `{"message": "recipe creation failed"}`, http.StatusInternalServerError)
		return
	}

	if errs, ok := result["errors"]; ok {
		log.Printf("GraphQL errors: %v", errs)
		http.Error(w, fmt.Sprintf(`{"message": "recipe creation failed", "errors": %v}`, errs), http.StatusInternalServerError)
		return
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		log.Println("Invalid response format from Hasura")
		http.Error(w, `{"message": "recipe creation failed"}`, http.StatusInternalServerError)
		return
	}

	insertResult, ok := data["insert_recipes_one"].(map[string]interface{})
	if !ok {
		log.Println("No insert result from Hasura")
		http.Error(w, `{"message": "recipe creation failed"}`, http.StatusInternalServerError)
		return
	}

	createdAt, _ := insertResult["created_at"].(string)

	// Respond with created recipe
	response := RecipeOutput{
		ID:               recipeID,
		Title:            input.Title,
		Description:      input.Description,
		PrepTime:         input.PrepTime,
		CookTime:         input.CookTime,
		TotalTime:        &totalTime,
		Servings:         input.Servings,
		FeaturedImageURL: input.FeaturedImageURL,
		Price: input.Price,
		CreatedAt:        createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func doHasuraRequest(query string) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{"query": query}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", hasuraURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", hasuraSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}

func buildStepsInsert(recipeID string, steps []RecipeStepInput) string {
	if len(steps) == 0 {
		return ""
	}

	var values []string
	for _, s := range steps {
		values = append(values, fmt.Sprintf(
			`{recipe_id: "%s", step_number: %d, instruction: "%s", image_url: %s}`,
			recipeID, s.StepNumber, escapeString(s.Instruction), nullOrStr(s.ImageURL),
		))
	}

	return fmt.Sprintf(`insert_recipe_steps(objects: [%s]) { affected_rows }`, strings.Join(values, ", "))
}

func buildIngredientsInsert(recipeID string, ingredients []RecipeIngredientInput) string {
	if len(ingredients) == 0 {
		return ""
	}

	var values []string
	for _, i := range ingredients {
		values = append(values, fmt.Sprintf(
			`{recipe_id: "%s", name: "%s", quantity: %s, unit: %s}`,
			recipeID, escapeString(i.Name), nullOrFloat(i.Quantity), nullOrStr(i.Unit),
		))
	}

	return fmt.Sprintf(`insert_recipe_ingredients(objects: [%s]) { affected_rows }`, strings.Join(values, ", "))
}

func buildImagesInsert(recipeID string, images []RecipeImageInput) string {
	if len(images) == 0 {
		return ""
	}

	var values []string
	for _, img := range images {
		values = append(values, fmt.Sprintf(
			`{recipe_id: "%s", image_url: "%s", is_featured: %t}`,
			recipeID, escapeString(img.ImageURL), boolOrFalse(img.IsFeatured),
		))
	}

	return fmt.Sprintf(`insert_images(objects: [%s]) { affected_rows }`, strings.Join(values, ", "))
}

func buildCategoryInsert(recipeID string, categoryIDs []string) string {
	if len(categoryIDs) == 0 {
		return ""
	}

	var values []string
	for _, id := range categoryIDs {
		values = append(values, fmt.Sprintf(`{recipe_id: "%s", category_id: "%s"}`, recipeID, id))
	}

	return fmt.Sprintf(`insert_recipe_categories(objects: [%s]) { affected_rows }`, strings.Join(values, ", "))
}

// Helper functions
func nullOrStr(s *string) string {
	if s == nil {
		return "null"
	}
	return fmt.Sprintf(`"%s"`, escapeString(*s))
}

func nullOrInt(i *int) string {
	if i == nil {
		return "null"
	}
	return fmt.Sprintf("%d", *i)
}

func nullOrFloat(f *float64) string {
	if f == nil {
		return "null"
	}
	return fmt.Sprintf("%f", *f)
}

func boolOrFalse(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	return s
}