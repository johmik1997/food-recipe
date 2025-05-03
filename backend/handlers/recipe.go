package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
)

// Define the structure matching your GraphQL response
type recipe struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// GetRecipesHandler returns a Gin handler that fetches recipes using GraphQL
func GetRecipesHandler(client *graphql.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query struct {
			Recipes []recipe `graphql:"recipes"`
		}

		err := client.Query(context.Background(), &query, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, query.Recipes)
	}
}
