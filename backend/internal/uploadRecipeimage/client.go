package uploadRecipeimage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"foodrecipe/utils"
	"io"
	"log"
	"net/http"
)

// GraphQLClient interface for inserting recipe image
type GraphQLClient interface {
	InsertRecipeImage(ctx context.Context, recipeID, imageURL string, isFeatured bool) error
	UpdateRecipeFeaturedImage(ctx context.Context, recipeID, imageURL string) error
}

// HasuraClient implements GraphQLClient
type HasuraClient struct {
	Endpoint    string
	AdminSecret string
}

func NewHasuraClient() *HasuraClient {
	return &HasuraClient{
		Endpoint:    utils.GraphQLURL,
		AdminSecret: utils.AdminSecret,
	}
}

// Implement GraphQLClient interface
func (c *HasuraClient) InsertRecipeImage(ctx context.Context, recipeID, imageURL string, isFeatured bool) error {
	log.Printf("[Hasura] Inserting image for recipe: %s", recipeID)
	log.Printf("[Hasura] Image URL: %s, Featured: %t", imageURL, isFeatured)

	mutation := `
	mutation insertRecipeImage($recipe_id: uuid!, $image_url: String!, $is_featured: Boolean!) {
		insert_recipe_images_one(object: {
			recipe_id: $recipe_id,
			image_url: $image_url,
			is_featured: $is_featured
		}) {
			id
			image_url
		}
	}`

	variables := map[string]interface{}{
		"recipe_id":   recipeID,
		"image_url":   imageURL,
		"is_featured": isFeatured,
	}

log.Print("[Hasura] Updating featured image for recipe: ", recipeID)
log.Print("[Hasura] New featured image URL: ", imageURL)
log.Print("featured: ", isFeatured)
	reqBody, err := json.Marshal(map[string]interface{}{
		"query":     mutation,
		"variables": variables,
	})
	if err != nil {
		log.Printf("[Hasura] Marshal error: %v", err)
		return err
	}

	log.Printf("[Hasura] Mutation: %s", string(reqBody))

	req, err := http.NewRequestWithContext(ctx, "POST", c.Endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("[Hasura] Request creation error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", c.AdminSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[Hasura] Request execution error: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Hasura] Error response: %s", string(body))
		return fmt.Errorf("hasura error: %s", string(body))
	}

	log.Println("[Hasura] Image inserted successfully")
	return nil
}

func (c *HasuraClient) UpdateRecipeFeaturedImage(ctx context.Context, recipeID, imageURL string) error {
    mutation := `
    mutation updateRecipeFeaturedImage($id: uuid!, $image_url: String!) {
        update_recipes_by_pk(
            pk_columns: {id: $id}
            _set: {feature_image_url: $image_url}
        ) {
            id
        }
    }`
log.Print("[Hasura] Updating featured image for recipe: ", recipeID)
log.Print("[Hasura] New featured image URL: ", imageURL)
    variables := map[string]interface{}{
        "id":        recipeID,
        "image_url": imageURL,
    }
	
	reqBody, err := json.Marshal(map[string]interface{}{
		"query":     mutation,
		"variables": variables,
	})
	if err != nil {
		log.Printf("[Hasura] Marshal error: %v", err)
		return err
	}

	log.Printf("[Hasura] Mutation: %s", string(reqBody))

	req, err := http.NewRequestWithContext(ctx, "POST", c.Endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("[Hasura] Request creation error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", c.AdminSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[Hasura] Request execution error: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Hasura] Error response: %s", string(body))
		return fmt.Errorf("hasura error: %s", string(body))
	}

	log.Println("[Hasura] Image inserted successfully")
	return nil
}