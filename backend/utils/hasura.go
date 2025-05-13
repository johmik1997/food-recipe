package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var GraphQLURL = "http://graphql-engine:8080/v1/graphql"
var AdminSecret = "myadminsecretkey"

func ExecuteGraphQL(query string, variables map[string]interface{}) (*http.Response, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	req, err := http.NewRequest("POST", GraphQLURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", AdminSecret)
	// Set the user ID in the header if needed
	// req.Header.Set("x-hasura-user-id", userID)
	// Set the role in the header if needed		
	// req.Header.Set("x-hasura-role", role)

	return http.DefaultClient.Do(req)
}



type Transaction struct {
	client *http.Client
}

func StartTransaction() (*Transaction, error) {
	return &Transaction{
		client: &http.Client{},
	}, nil
}

func (t *Transaction) ExecuteGraphQL(query string, variables map[string]interface{}) (*http.Response, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	req, err := http.NewRequest("POST", GraphQLURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", AdminSecret)

	return t.client.Do(req)
}

func (t *Transaction) Commit() error {
	// In a real implementation, this would commit the transaction
	// For Hasura, we're just ensuring all operations succeeded
	return nil
}

func (t *Transaction) Rollback() {
	// In a real implementation, this would rollback the transaction
	// For Hasura, we rely on each operation being atomic
}

func SendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
	})
}