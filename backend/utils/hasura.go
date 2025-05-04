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

	return http.DefaultClient.Do(req)
}