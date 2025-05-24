package auth

import (
	"context"
	"errors"
	"fmt"
	"foodrecipe/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
    ErrPasswordHashing   = errors.New("failed to hash password")
    ErrTokenGeneration   = errors.New("failed to generate token")
    ErrUserCreation      = errors.New("failed to create user")
    // Add other error variables here if needed
)
type Service struct {
	userClient *UserClient
	config     *utils.Config
}
// NewService creates a new auth Service instance
func NewService(userClient *UserClient, config *utils.Config) *Service {
    return &Service{
        userClient: userClient,
        config:     config,
    }
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
    user, err := s.userClient.GetUserByEmail(ctx, input.Email)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, ErrInvalidCredentials
        }
        return nil, fmt.Errorf("login failed: %w", err)
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        return nil, ErrInvalidCredentials
    }

    token, err := utils.GenerateJWT(user.ID, user.Email, []byte(s.config.JWTSecret))
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrTokenGeneration, err)
    }

    return &LoginResponse{
        Token: token,
        User: User{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
        },
    }, nil
}


func (s *Service) Register(ctx context.Context, input RegisterInput) (*RegisterResponse, error) {
    // 1. Check if user exists
    _, err := s.userClient.GetUserByEmail(ctx, input.Email)
    if err == nil {
        return nil, errors.New("user with this email already exists")
    }
    if !errors.Is(err, ErrUserNotFound) {
        return nil, fmt.Errorf("failed to check user existence: %w", err)
    }

    // 2. Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    // 3. Create user
    user, err := s.userClient.CreateUser(ctx, input.Name, input.Email, string(hashedPassword))
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // 4. Generate token
    token, err := utils.GenerateJWT(user.ID, user.Email, []byte(s.config.JWTSecret))
    if err != nil {
        return nil, fmt.Errorf("failed to generate token: %w", err)
    }

    return &RegisterResponse{
        Token: token,
        User: User{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
        },
    }, nil
}
type RefreshResponse struct {
    Token string `json:"token"`
}
func (s *Service) RefreshToken(ctx context.Context, tokenString string) (*RefreshResponse, error) {
    claims, err := utils.ValidateJWT(tokenString, []byte(s.config.JWTSecret))
    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }

    newToken, err := utils.GenerateJWT(
        claims["sub"].(string),
        claims["email"].(string),
        []byte(s.config.JWTSecret),
    )
    if err != nil {
        return nil, err
    }

    return &RefreshResponse{Token: newToken}, nil
}