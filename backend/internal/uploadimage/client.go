package uploadimage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Interface for abstraction (optional, can remove if not needed)
type GraphQLClient interface {
	UpdateUserAvatar(userID, avatarURL string) error
}

// HTTP-based Hasura client
type hasuraClient struct {
	url         string
	adminSecret string
}

// Configuration struct
type ClientConfig struct {
	URL         string
	AdminSecret string
}

// Constructor
func NewClient(cfg ClientConfig) GraphQLClient {
	return &hasuraClient{
		url:         cfg.URL,
		adminSecret: cfg.AdminSecret,
	}
}

// Raw HTTP GraphQL mutation
func (h *hasuraClient) UpdateUserAvatar(userID, avatarURL string) error {
	mutation := `mutation UpdateUserAvatar($user_id: uuid!, $avatar_url: String!) {
		update_users_by_pk(pk_columns: {id: $user_id}, _set: {avatar_image_url: $avatar_url}) {
			id
		}
	}`

	reqBody, err := json.Marshal(map[string]interface{}{
		"query": mutation,
		"variables": map[string]interface{}{
			"user_id":    userID,
			"avatar_url": avatarURL,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", h.adminSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("hasura returned %d: %s", resp.StatusCode, string(body))
	}

	// Optional: Parse and check response if needed
	var result struct {
		Data struct {
			UpdateUsersByPk struct {
				ID string `json:"id"`
			} `json:"update_users_by_pk"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Data.UpdateUsersByPk.ID == "" {
		return fmt.Errorf("no user updated")
	}

	log.Printf("Avatar updated for user: %s", result.Data.UpdateUsersByPk.ID)
	return nil
}
