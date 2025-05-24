// package graphqlclient/client.go
package uploadimage

import (
	"net/http"

	"github.com/hasura/go-graphql-client"
)

type ClientConfig struct {
	URL         string
	AdminSecret string
	JWTToken    string
}

type Client struct {
	*graphql.Client
}

func NewClient(cfg ClientConfig) *Client {
	return &Client{
		Client: graphql.NewClient(cfg.URL, nil).
			WithRequestModifier(func(r *http.Request) {
				r.Header.Set("Content-Type", "application/json")
				if cfg.AdminSecret != "" {
					r.Header.Set("x-hasura-admin-secret", cfg.AdminSecret)
				}
				if cfg.JWTToken != "" {
					r.Header.Set("Authorization", "Bearer "+cfg.JWTToken)
				}
			}),
	}
}