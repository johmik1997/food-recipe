package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foodrecipe/models"
	"io"
	"net/http"
	"os"
)

type DatabaseInterface interface {
	GetRecipe(recipeID string) (*models.Recipe, error)
	CreatePayment(payment models.Payment) (string, error)
	UpdatePaymentStatus(paymentID string, status string) error
	UpdatePaymentStatusByTxRef(txRef string, status string, chapaRefID string) error
	GrantRecipeAccess(txRef string) error
}

type HasuraDB struct {
	GraphQLURL string
	AdminSecret string
}

func (db *HasuraDB) ExecuteGraphQL(query string, variables map[string]interface{}) ([]byte, error) {
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	if db.GraphQLURL == "" {
		db.GraphQLURL = os.Getenv("HASURA_GRAPHQL_URL")
	}
	if db.AdminSecret == "" {
		db.AdminSecret = os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")
	}

	req, err := http.NewRequest("POST", db.GraphQLURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", db.AdminSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("graphql error: %s", string(body))
	}

	return body, nil
}

func (db *HasuraDB) GetRecipe(recipeID string) (*models.Recipe, error) {
	query := `
	query GetRecipe($id: uuid!) {
		recipes_by_pk(id: $id) {
			id
			price
		}
	}`

	vars := map[string]interface{}{"id": recipeID}
	resp, err := db.ExecuteGraphQL(query, vars)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Recipe *models.Recipe `json:"recipes_by_pk"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %v", result.Errors)
	}

	if result.Data.Recipe == nil {
		return nil, fmt.Errorf("recipe not found")
	}

	return result.Data.Recipe, nil
}

func (db *HasuraDB) CreatePayment(payment models.Payment) (string, error) {
	mutation := `
	mutation CreatePayment(
		$user_id: uuid!,
		$recipe_id: uuid!,
		$amount: numeric!,
		$currency: String!,
		$tx_ref: String!,
		$customer_email: String!,
		$customer_phone: String,
		$first_name: String,
		$last_name: String,
		$callback_url: String!,
		$return_url: String!
	) {
		insert_recipe_payments_one(object: {
			user_id: $user_id,
			recipe_id: $recipe_id,
			amount: $amount,
			currency: $currency,
			tx_ref: $tx_ref,
			customer_email: $customer_email,
			customer_phone: $customer_phone,
			first_name: $first_name,
			last_name: $last_name,
			callback_url: $callback_url,
			return_url: $return_url
		}) {
			id
		}
	}`

	vars := map[string]interface{}{
		"user_id":        payment.UserID,
		"recipe_id":      payment.RecipeID,
		"amount":         payment.Amount,
		"currency":       payment.Currency,
		"tx_ref":         payment.TxRef,
		"customer_email": payment.CustomerEmail,
		"customer_phone": payment.CustomerPhone,
		"first_name":     payment.FirstName,
		"last_name":      payment.LastName,
		"callback_url":   payment.CallbackURL,
		"return_url":     payment.ReturnURL,
	}

	resp, err := db.ExecuteGraphQL(mutation, vars)
	if err != nil {
		return "", err
	}

	var result struct {
		Data struct {
			Payment struct {
				ID string `json:"id"`
			} `json:"insert_payments_one"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Data.Payment.ID, nil
}

func (db *HasuraDB) UpdatePaymentStatus(paymentID string, status string) error {
	return db.UpdatePaymentStatusByFields(map[string]interface{}{
		"id":     paymentID,
		"status": status,
	})
}

func (db *HasuraDB) UpdatePaymentStatusByTxRef(txRef string, status string, chapaRefID string) error {
	updateFields := map[string]interface{}{
		"tx_ref": txRef,
		"status": status,
	}
	if chapaRefID != "" {
		updateFields["chapa_ref_id"] = chapaRefID
	}
	return db.UpdatePaymentStatusByFields(updateFields)
}

func (db *HasuraDB) UpdatePaymentStatusByFields(fields map[string]interface{}) error {
	mutation := `
	mutation UpdatePayment($where: payments_bool_exp!, $_set: payments_set_input!) {
		update_recipe_payments(where: $where, _set: $_set) {
			affected_rows
		}
	}`

	where := make(map[string]interface{})
	for k, v := range fields {
		if k != "status" && k != "chapa_ref_id" {
			where[k] = map[string]interface{}{"_eq": v}
		}
	}

	set := make(map[string]interface{})
	if status, ok := fields["status"]; ok {
		set["status"] = status
	}
	if chapaRefID, ok := fields["chapa_ref_id"]; ok {
		set["chapa_ref_id"] = chapaRefID
	}

	vars := map[string]interface{}{
		"where": where,
		"_set":  set,
	}

	_, err := db.ExecuteGraphQL(mutation, vars)
	return err
}

func (db *HasuraDB) GrantRecipeAccess(txRef string) error {
	// First get payment details
	query := `
	query GetPayment($tx_ref: String!) {
		recipe_payments(where: {tx_ref: {_eq: $tx_ref}}) {
			user_id
			recipe_id
			id
		}
	}`

	vars := map[string]interface{}{"tx_ref": txRef}
	resp, err := db.ExecuteGraphQL(query, vars)
	if err != nil {
		return err
	}

	var result struct {
		Data struct {
			Payments []struct {
				UserID   string `json:"user_id"`
				RecipeID string `json:"recipe_id"`
				ID       string `json:"id"`
			} `json:"payments"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result.Data.Payments) == 0 {
		return fmt.Errorf("payment not found")
	}

	payment := result.Data.Payments[0]

	// Grant access
	mutation := `
	mutation GrantAccess(
		$user_id: uuid!,
		$recipe_id: uuid!,
		$payment_id: uuid!
	) {
		insert_recipe_access_one(object: {
			user_id: $user_id,
			recipe_id: $recipe_id,
			payment_id: $payment_id
		}) {
			user_id
		}
		update_recipe_payments_by_pk(
			pk_columns: {id: $payment_id},
			_set: {access_granted: true}
		) {
			id
		}
	}`

	vars = map[string]interface{}{
		"user_id":    payment.UserID,
		"recipe_id":  payment.RecipeID,
		"payment_id": payment.ID,
	}

	_, err = db.ExecuteGraphQL(mutation, vars)
	return err
}