package graphql

import (
	"bytes"
	"encoding/json"
	"net/http"

	"foodrecipe/logger"
)

type Client struct {
	url    string
	client *http.Client
	log    logger.Logger
}

func NewClient(url string, log logger.Logger) *Client {
	return &Client{
		url:    url,
		client: &http.Client{},
		log:    log,
	}
}

func (c *Client) Run(query string, variables map[string]interface{}, response interface{}) error {
	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}