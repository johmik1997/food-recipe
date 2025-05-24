package auth

// Matches your GraphQL action schema exactly
type (
    LoginInput struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    RegisterInput struct {
        Email    string `json:"email"`
        Password string `json:"password"`
        Name     string `json:"name"`
    }

    User struct {
        ID       string `json:"id"`
        Email    string `json:"email"`
        Name     string `json:"name"`
        Password string `json:"-"` // Never marshal to JSON
    }

    LoginResponse struct {
        Token string `json:"token"`
        User  User   `json:"user"`
    }

    RegisterResponse struct {
        Token string `json:"token"`
        User  User   `json:"user"`
    }
)