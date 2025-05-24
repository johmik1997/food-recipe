package auth

import (
	"encoding/json"
	"errors"
	"foodrecipe/internal/utils"
	"net/http"
	"strings"
)

// Config holds configuration for the auth handler, such as JWTSecret.
type Config struct {
    JWTSecret string
}

var ErrInvalidCredentials = errors.New("invalid credentials")

type Handler struct {
    service *Service
    config  *Config
}
func NewHandler(service *Service, config *Config) *Handler {
    return &Handler{service: service, config: config}
}


func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Input LoginInput `json:"input"`
    }

    if err := utils.ParseJSON(r, &req); err != nil {
        utils.WriteJSONError(w, http.StatusBadRequest, "invalid request format")
        return
    }

        response, err := h.service.Login(r.Context(), req.Input)
        if err != nil {
            status := http.StatusInternalServerError
            if errors.Is(err, ErrInvalidCredentials) {
                status = http.StatusUnauthorized
            }
            utils.WriteJSONError(w, status, err.Error())
            return
        }
    
        utils.WriteJSON(w, http.StatusOK, response)
        return
    }

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
//     var req struct {
//         Input RegisterInput `json:"input"`
//     }
  var req struct {
        Input RegisterInput `json:"input"`
    }
    if err := utils.ParseJSON(r, &req); err != nil {
        utils.WriteJSONError(w, http.StatusBadRequest, "invalid request format")
        return
    }

    response, err := h.service.Register(r.Context(), req.Input)
    if err != nil {
        status := http.StatusInternalServerError
        if errors.Is(err, ErrUserExists) {
            status = http.StatusConflict
        }
        utils.WriteJSONError(w, status, err.Error())
        return
    }

    utils.WriteJSON(w, http.StatusOK, response)}

// func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
//     w.Header().Set("Content-Type", "application/json")

//     // Parse the nested Hasura action input
//     var requestBody struct {
//         Input struct {
//             Email    string `json:"email"`
//             Name     string `json:"name"`
//             Password string `json:"password"`
//         } `json:"input"`
//     }

//     if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
//         writeHasuraError(w, "Invalid request format", http.StatusBadRequest)
//         return
//     }

//     // Validate input
//     if requestBody.Input.Email == "" {
//         writeHasuraError(w, "email is required", http.StatusBadRequest)
//         return
//     }
//     if requestBody.Input.Password == "" {
//         writeHasuraError(w, "password is required", http.StatusBadRequest)
//         return
//     }
//     if requestBody.Input.Name == "" {
//         writeHasuraError(w, "name is required", http.StatusBadRequest)
//         return
//     }

//     // Process registration
//     response, err := h.service.Register(r.Context(), RegisterInput{
//         Email:    requestBody.Input.Email,
//         Name:     requestBody.Input.Name,
//         Password: requestBody.Input.Password,
//     })
//     if err != nil {
//         status := http.StatusInternalServerError
//         if strings.Contains(err.Error(), "already exists") {
//             status = http.StatusConflict
//         }
//         writeHasuraError(w, err.Error(), status)
//         return
//     }

//     // Return success response in Hasura action format
//     w.WriteHeader(http.StatusOK)
//     json.NewEncoder(w).Encode(map[string]interface{}{
//         "token": response.Token,
//         "user": map[string]interface{}{
//             "id":    response.User.ID,
//             "email": response.User.Email,
//             "name":  response.User.Name,
//         },
//     })
// }

// Helper function for consistent GraphQL error responses
func writeGraphQLError(w http.ResponseWriter, message string, statusCode int) {
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "errors": []map[string]interface{}{
            {
                "message": message,
                "extensions": map[string]interface{}{
                    "code": strings.ToLower(http.StatusText(statusCode)),
                },
            },
        },
    })
}

// HandleRefreshToken remains similar but ensure it's at the correct endpoint
func (h *Handler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var req struct {
        Token string `json:"token"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeGraphQLError(w, "invalid request format", http.StatusBadRequest)
        return
    }
    
    response, err := h.service.RefreshToken(r.Context(), req.Token)
    if err != nil {
        writeGraphQLError(w, err.Error(), http.StatusUnauthorized)
        return
    }
    
    result := map[string]interface{}{
        "data": map[string]interface{}{
            "refreshToken": response,
        },
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(result)
}


// Hasura-compatible error response
type HasuraError struct {
    Message string `json:"message"`
}

// func writeHasuraError(w http.ResponseWriter, message string, statusCode int) {
//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(statusCode)
//     json.NewEncoder(w).Encode(HasuraError{
//         Message: message,
//     })
// }