package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)


type PaymentRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	SuccessURL  string `json:"success_url"`
	ErrorURL    string `json:"error_url"`
}

func generateTxRef() string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(1000000)
	return fmt.Sprintf("chewatatest-%d", randomNum)
}

func verifyPayment(txRef string) (string, error) {
	url := "https://api.chapa.co/v1/transaction/verify/" + txRef
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer CHASECK_TEST-pjjzrYoVZs6KvEXJzvi04kYxx8UsHACN")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response["status"] == "success" {
		return "Payment successful", nil
	} else {
		return "Payment failed", nil
	}
}

func HandlePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		handleCORS(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the payment request from frontend
	var paymentReq PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if paymentReq.Amount == "" || paymentReq.Email == "" || paymentReq.FirstName == "" || paymentReq.LastName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Generate a unique transaction reference
	txRef := generateTxRef()

	// Prepare the payment payload with dynamic values
	payload := strings.NewReader(fmt.Sprintf(`{
		"amount": "%s",
		"currency": "%s",
		"email": "%s",
		"first_name": "%s",
		"last_name": "%s",
		"phone_number": "%s",
		"tx_ref": "%s",
		"callback_url": "https://webhook.site/077164d6-29cb-40df-ba29-8a00e59a7e60",
		"return_url": "%s",
		"customization": {
			"title": "Payment merchant",
			"description": "I love online payments"
		},
		"meta": {
			"hide_receipt": true
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
	))

	url := "https://api.chapa.co/v1/transaction/initialize"
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", "Bearer CHASECK_TEST-pjjzrYoVZs6KvEXJzvi04kYxx8UsHACN")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach Chapa", http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var callbackData struct {
			Status   string `json:"status"`
			RefID    string `json:"ref_id"`
			TxRef    string `json:"tx_ref"`
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&callbackData)
		if err != nil {
			http.Error(w, "Failed to decode callback data", http.StatusBadRequest)
			return
		}

		status, err := verifyPayment(callbackData.TxRef)
		if err != nil {
			http.Error(w, "Error verifying payment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Payment verification status:", status)

		if callbackData.Status == "successful" {
			fmt.Println("Payment successful:", callbackData)
		} else {
			fmt.Println("Payment failed:", callbackData)
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var webhookData struct {
			Status   string `json:"status"`
			RefID    string `json:"ref_id"`
			TxRef    string `json:"tx_ref"`
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&webhookData)
		if err != nil {
			http.Error(w, "Failed to decode webhook data", http.StatusBadRequest)
			return
		}

		status, err := verifyPayment(webhookData.TxRef)
		if err != nil {
			http.Error(w, "Error verifying payment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Payment verification status:", status)

		if webhookData.Status == "successful" {
			fmt.Println("Payment successful via webhook:", webhookData)
		} else {
			fmt.Println("Payment failed via webhook:", webhookData)
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handleCORS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}