// package handlers

// import (
// 	"foodrecipe/config"
// 	"foodrecipe/models"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/jmoiron/sqlx"
// )

// type LoginRequest struct {
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=8"`
// }

// func LoginHandler(db *sqlx.DB, cfg *config.Config) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var creds LoginRequest

// 		if err := c.ShouldBindJSON(&creds); err != nil {
// 			c.JSON(400, gin.H{
// 				"error": "Invalid request format",
// 				"code": "invalid_request",
// 			})
// 			return
// 		}

// 		if err := validator.New().Struct(creds); err != nil {
// 			ErrorResponse(c, 400, "Validation failed: " + err.Error())
// 			return
// 		}

// 		user, err := models.VerifyUser(db, creds.Email, creds.Password)
// 		if err != nil {
// 			ErrorResponse(c, 401, "Invalid credentials")
// 			return
// 		}

// 		claims := jwt.MapClaims{
// 			"sub": user.ID,
// 			"exp": time.Now().Add(cfg.TokenExpiry).Unix(),
// 			"iss": "foodrecipe-backend",
// 			"https://hasura.io/jwt/claims": jwt.MapClaims{
// 				"x-hasura-allowed-roles": []string{"user"},
// 				"x-hasura-default-role":  "user",
// 				"x-hasura-user-id":       user.ID,
// 			},
// 		}

// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
// 		if err != nil {
// 			ErrorResponse(c, 500, "Failed to generate token")
// 			return
// 		}

// 		c.JSON(200, gin.H{
// 			"token": tokenString,
// 			"user": gin.H{
// 				"id":    user.ID,
// 				"name":  user.Name,
// 				"email": user.Email,
// 			},
// 		})
// 	}
// }

// func ErrorResponse(c *gin.Context, code int, message string) {
// 	c.JSON(code, gin.H{"error": message})
// }
// func RegisterHandler(db *sqlx.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var user models.User
// 		if err := c.ShouldBindJSON(&user); err != nil {
// 			ErrorResponse(c, 400, "Invalid request format")
// 			return
// 		}

// 		if err := validator.New().Struct(user); err != nil {
// 			ErrorResponse(c, 400, "Validation failed: " + err.Error())
// 			return
// 		}

// 		if _, err := models.CreateUser(db, user.Name, user.Email, user.Password); err != nil {
// 			ErrorResponse(c, 500, "Failed to create user: " + err.Error())
// 			return
// 		}

// 		c.JSON(201, gin.H{"message": "User created successfully"})
// 	}
// }

// // package handlers

// // import (
// // 	"net/http"

// // 	"foodrecipe/graphql"
// // 	"foodrecipe/models"

// // 	"github.com/gin-gonic/gin"
// // )

// // type AuthHandler struct {
// // 	gqlClient *graphql.Client
// // }

// // func NewAuthHandler(client *graphql.Client) *AuthHandler {
// // 	return &AuthHandler{gqlClient: client}
// // }

// // // Login using GraphQL mutation
// // func (h *AuthHandler) Login(c *gin.Context) {
// // 	var req struct {
// // 		Email    string `json:"email"`
// // 		Password string `json:"password"`
// // 	}

// // 	if err := c.ShouldBindJSON(&req); err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// // 		return
// // 	}

// // 	var resp struct {
// // 		Data struct {
// // 			Login struct {
// // 				Token string      `json:"token"`
// // 				User  models.User `json:"user"`
// // 			} `json:"login"`
// // 		} `json:"data"`
// // 	}

// // 	err := h.gqlClient.Run(
// // 		graphql.loginQuery,
// // 		map[string]interface{}{
// // 			"email":    req.Email,
// // 			"password": req.Password,
// // 		},
// // 		&resp,
// // 	)

// // 	if err != nil {
// // 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// // 		return
// // 	}

// // 	c.JSON(http.StatusOK, gin.H{
// // 		"token": resp.Data.Login.Token,
// // 		"user":  resp.Data.Login.User,
// // 	})
// // }

// // // Register using GraphQL mutation
// // func (h *AuthHandler) Register(c *gin.Context) {
// // 	var input struct {
// // 		Name     string `json:"name"`
// // 		Email    string `json:"email"`
// // 		Password string `json:"password"`
// // 	}

// // 	if err := c.ShouldBindJSON(&input); err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// // 		return
// // 	}

// // 	var resp struct {
// // 		Data struct {
// // 			Register struct {
// // 				Token string      `json:"token"`
// // 				User  models.User `json:"user"`
// // 			} `json:"register"`
// // 		} `json:"data"`
// // 	}

// // 	err := h.gqlClient.Run(
// // 		graphql.registerQuery,
// // 		map[string]interface{}{
// // 			"input": map[string]interface{}{
// // 				"name":     input.Name,
// // 				"email":    input.Email,
// // 				"password": input.Password,
// // 			},
// // 		},
// // 		&resp,
// // 	)

// // 	if err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
// // 		return
// // 	}

// // 	c.JSON(http.StatusCreated, gin.H{
// // 		"token": resp.Data.Register.Token,
// // 		"user":  resp.Data.Register.User,
// // 	})
// // // }
// // package handlers

// // import (
// // 	"foodrecipe/models"
// // 	"foodrecipe/utils"
// // 	"net/http"

// // 	"github.com/gin-gonic/gin"
// // 	"github.com/jmoiron/sqlx"
// // )

// // type LoginActionInput struct {
// // 	Email    string `json:"email"`
// // 	Password string `json:"password"`
// // }

// // type LoginActionResponse struct {
// // 	Token string     `json:"token"`
// // 	User  models.User `json:"user"`
// // }

// // func LoginAction(c *gin.Context) {
// // 	var input struct {
// // 		Input LoginActionInput `json:"input"`
// // 	}

// // 	if err := c.ShouldBindJSON(&input); err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// // 		return
// // 	}

// // 	// Get DB from context or package variable
// // 	db := c.MustGet("db").(*sqlx.DB)

// // 	user, err := models.VerifyUser(db, input.Input.Email, input.Input.Password)
// // 	if err != nil {
// // 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// // 		return
// // 	}

// // 	token, err := utils.GenerateToken(user)
// // 	if err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// // 		return
// // 	}

// // 	c.JSON(http.StatusOK, gin.H{
// // 		"token": token,
// // 		"user":  user,
// // 	})
// // }
// // func RegisterAction(c *gin.Context) {
// // 	var input struct {
// // 		Input models.User `json:"input"`
// // 	}

// // 	if err := c.ShouldBindJSON(&input); err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// // 		return
// // 	}

// // 	// Get DB from context or package variable
// // 	db := c.MustGet("db").(*sqlx.DB)
// // 	user, err := models.CreateUser(db, input.Input.Name, input.Input.Email, input.Input.PasswordHash)
// // 	if err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// // 		return
// // 	}

// // 	c.JSON(http.StatusCreated, gin.H{
// // 		"user": user,
// // 	})
// // }

// package handlers

// import (
// 	"golang.org/x/crypto/bcrypt"
// )

//	func HashPassword(password string) (string, error) {
//		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//		if err != nil {
//			return "", err
//		}
//		return string(hashedBytes), nil
//	}
package handlers