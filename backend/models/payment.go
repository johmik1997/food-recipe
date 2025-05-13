package models

import "time"

// ChapaPaymentRequest represents the request payload for Chapa payment initialization
type ChapaPaymentRequest struct {
    Amount       string `json:"amount"`
    Currency     string `json:"currency"`
    Email        string `json:"email"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
    PhoneNumber  string `json:"phone_number,omitempty"`
    TxRef        string `json:"tx_ref"`
    CallbackURL  string `json:"callback_url"`
    ReturnURL    string `json:"return_url"`
    Customization struct {
        Title       string `json:"title,omitempty"`
        Description string `json:"description,omitempty"`
    } `json:"customization,omitempty"`
    Meta map[string]interface{} `json:"meta,omitempty"`
}

// ChapaPaymentResponse represents the response from Chapa payment initialization
type ChapaPaymentResponse struct {
    Message string `json:"message"`
    Status  string `json:"status"`
    Data    struct {
        CheckoutURL string `json:"checkout_url"`
    } `json:"data"`
}

// VerifyPaymentResponse represents the response from Chapa payment verification
type VerifyPaymentResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    struct {
        ID           string    `json:"id"`
        Amount       float64   `json:"amount"`
        Currency     string    `json:"currency"`
        Status       string    `json:"status"`
        TxRef        string    `json:"tx_ref"`
        ChargedAmount float64  `json:"charged_amount"`
        Fee          float64   `json:"fee"`
        CreatedAt    time.Time `json:"created_at"`
        Customer     struct {
            FirstName   string `json:"first_name"`
            LastName    string `json:"last_name"`
            Email       string `json:"email"`
            PhoneNumber string `json:"phone_number"`
        } `json:"customer"`
        Metadata map[string]interface{} `json:"meta"`
    } `json:"data"`
}

// Payment represents a payment record in the database
type Payment struct {
    ID            string    `json:"id"`
    UserID        string    `json:"user_id"`
    RecipeID      string    `json:"recipe_id"`
    Amount        float64   `json:"amount"`
    Currency      string    `json:"currency"`
    TxRef         string    `json:"tx_ref"`
    ChapaRefID    string    `json:"chapa_ref_id,omitempty"`
    CustomerEmail string    `json:"customer_email"`
    CustomerPhone string    `json:"customer_phone,omitempty"`
    FirstName     string    `json:"first_name,omitempty"`
    LastName      string    `json:"last_name,omitempty"`
    CallbackURL   string    `json:"callback_url"`
    ReturnURL     string    `json:"return_url"`
    Status        string    `json:"status"`
    AccessGranted bool      `json:"access_granted"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

// PaymentCallback represents the webhook payload from Chapa
type PaymentCallback struct {
    TxRef      string `json:"tx_ref"`
    Status     string `json:"status"`
    ChapaRefID string `json:"chapa_ref_id,omitempty"`
}

// PaymentResponse represents the response after payment initiation
type PaymentResponse struct {
    CheckoutURL string `json:"checkout_url"`
    TxRef      string `json:"tx_ref"`
    PaymentID  string `json:"payment_id"`
}

// PaymentCallbackResponse represents the response for payment callback
type PaymentCallbackResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

// Recipe represents recipe information needed for payment processing
type Recipe struct {
    ID       string  `json:"id"`
    Price    float64 `json:"price"`
}

// RecipeAccess represents granted access to a premium recipe
type RecipeAccess struct {
    UserID    string    `json:"user_id"`
    RecipeID  string    `json:"recipe_id"`
    PaymentID string    `json:"payment_id"`
    GrantedAt time.Time `json:"granted_at"`
}