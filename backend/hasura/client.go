package hasura

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foodrecipe/models"
	"net/http"
)

const hasuraEndpoint = "http://graphql-engine:8080/v1/graphql"

const registerMutation = `

mutation Register($input: RegisterInput!) {
  register(input: $input) {
    token
    user { id name email }
  }
}`

type Client struct {
	Endpoint string
}

func NewClient() *Client {
	return &Client{
		Endpoint: hasuraEndpoint,
	}
}

func (c *Client) ExecuteRegister(variables map[string]interface{}, headers http.Header) (*models.HasuraResponse, error) {
	requestBody := map[string]interface{}{
		"query":     registerMutation,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Copy headers from original request
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Hasura: %v", err)
	}
	defer resp.Body.Close()

	var result models.HasuraResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
// Add this method to your existing Client struct
func (c *Client) ExecuteQuery(query string, variables map[string]interface{}) (*struct {
	Data struct {
		Users []struct {
			ID       int
			Name     string
			Email    string
			Password string
		}
	}
	Errors []struct {
		Message string
	}
}, error) {
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Hasura: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Users []struct {
				ID       int
				Name     string
				Email    string
				Password string
			}
		}
		Errors []struct {
			Message string
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("hasura error: %v", result.Errors[0].Message)
	}

	return &result, nil
}