package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PaymentRequest struct {
	RecipeID    uuid.UUID  `json:"recipe_id"`
	Amount      float64    `json:"amount"`
	Currency    string  `json:"currency"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber string  `json:"phone_number"`
	SuccessURL  string  `json:"success_url"`
	ErrorURL    string  `json:"error_url"`
}

type HasuraRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type HasuraResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func getRecipePriceFromHasura(recipeID string) (float64, string, error) {
	query := `
		query GetRecipePrice($id: uuid!) {
			recipes_by_pk(id: $id) {
				id
				title
				price
			}
		}
	`

	variables := map[string]interface{}{
		"id": recipeID,
	}

	requestBody, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	req, err := http.NewRequest("POST", "http://graphql-engine:8080/v1/graphql", strings.NewReader(string(requestBody)))
	if err != nil {
		return 0, "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", "myadminsecretkey") // Or use appropriate auth

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	var result HasuraResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, "", err
	}

	if len(result.Errors) > 0 {
		return 0, "", fmt.Errorf("hasura error: %s", result.Errors[0].Message)
	}

	recipeData, ok := result.Data["recipes_by_pk"].(map[string]interface{})
	if !ok {
		return 0, "", fmt.Errorf("recipe not found")
	}

	price, ok := recipeData["price"].(float64)
	if !ok {
		return 0, "", fmt.Errorf("invalid price format")
	}

	title, _ := recipeData["title"].(string)

	return price, title, nil
}

func generateTxRef() string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(1000000)
	return fmt.Sprintf("chewatatest-%d", randomNum)
}

func handlePayment(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

    if r.Method == http.MethodOptions {
        handleCORS(w, r)
        return
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse the payment request
    var paymentReq PaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate required fields
    if paymentReq.RecipeID == uuid.Nil || paymentReq.Email == "" || 
       paymentReq.FirstName == "" || paymentReq.LastName == "" || paymentReq.Amount <= 0 {
        http.Error(w, "Missing required fields or invalid amount", http.StatusBadRequest)
        return
    }

    // Get recipe price from Hasura
    recipePrice, recipeTitle, err := getRecipePriceFromHasura(paymentReq.RecipeID.String())
    if err != nil {
        http.Error(w, "Failed to verify recipe: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Verify the payment amount matches the recipe price
    // Verify the payment amount matches the recipe price
if paymentReq.Amount != recipePrice {
    http.Error(w, fmt.Sprintf("Payment amount (%.2f) does not match recipe price (%.2f)", 
        paymentReq.Amount, recipePrice), http.StatusBadRequest)
    return
}
    // Generate a unique transaction reference
    txRef := generateTxRef()

    // Prepare the payment payload for Chapa
    payload := strings.NewReader(fmt.Sprintf(`{
        "amount": "%.2f",
        "currency": "%s",
        "email": "%s",
        "first_name": "%s",
        "last_name": "%s",
        "phone_number": "%s",
        "tx_ref": "%s",
        "callback_url": "https://your-webhook-url/callback",
        "return_url": "%s",
        "customization": {
            "title": "Purchase: %s",
            "description": "Recipe purchase"
        },
        "meta": {
            "recipe_id": "%s",
            "user_email": "%s"
        }
    }`,
        paymentReq.Amount,
        paymentReq.Currency,
        paymentReq.Email,
        paymentReq.FirstName,
        paymentReq.LastName,
        paymentReq.PhoneNumber,
        txRef,
        paymentReq.SuccessURL,
        recipeTitle,
        paymentReq.RecipeID.String(),
        paymentReq.Email,
    ))

    // Initialize payment with Chapa
    req, err := http.NewRequest("POST", "https://api.chapa.co/v1/transaction/initialize", payload)
    if err != nil {
        http.Error(w, "Error creating payment request", http.StatusInternalServerError)
        return
    }

    req.Header.Add("Authorization", "Bearer CHASECK_TEST-pjjzrYoVZs6KvEXJzvi04kYxx8UsHACN")
    req.Header.Add("Content-Type", "application/json")

    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        http.Error(w, "Failed to reach payment provider", http.StatusBadGateway)
        return
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        http.Error(w, "Failed to read payment response", http.StatusInternalServerError)
        return
    }

    // Forward the Chapa response directly to client
    w.WriteHeader(res.StatusCode)
    w.Write(body)
}

func handleCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}