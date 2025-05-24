package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hasura/go-graphql-client"
)

var ErrUserNotFound = errors.New("user not found")

type UserClient struct {
    client *graphql.Client
}

func NewUserClient(client *graphql.Client) *UserClient {
	return &UserClient{
		client: client,
	}
}

// func (c *UserClient) GetUserByEmail(ctx context.Context, email string) (*User, error) {
//     // First try exact match
//     var query struct {
//         Users []struct {
//             ID       string `graphql:"id"`
//             Email    string
//             Password string
//             Name     string
//         } `graphql:"users(where: {email: {_eq: $email}}, limit: 1)"`
//     }

//     err := c.client.Query(ctx, &query, map[string]interface{}{
//         "email": graphql.String(email),
//     })
//     if err != nil {
//         return nil, fmt.Errorf("exact match query failed: %w", err)
//     }

//     if len(query.Users) > 0 {
//         return &User{
//             ID:       query.Users[0].ID,
//             Email:    query.Users[0].Email,
//             Password: query.Users[0].Password,
//             Name:     query.Users[0].Name,
//         }, nil
//     }

//     // Fallback to case-insensitive search if exact match fails
//     var ciQuery struct {
//         Users []struct {
//             ID       string `graphql:"id"`
//             Email    string
//             Password string
//             Name     string
//         } `graphql:"users(where: {email: {_ilike: $email}}, limit: 1)"`
//     }

//     err = c.client.Query(ctx, &ciQuery, map[string]interface{}{
//         "email": graphql.String(email),
//     })
//     if err != nil {
//         return nil, fmt.Errorf("case-insensitive query failed: %w", err)
//     }

//     if len(ciQuery.Users) > 0 {
//         return &User{
//             ID:       ciQuery.Users[0].ID,
//             Email:    ciQuery.Users[0].Email,
//             Password: ciQuery.Users[0].Password,
//             Name:     ciQuery.Users[0].Name,
//         }, nil
//     }

//     return nil, ErrUserNotFound
// }

// GetUserByEmail retrieves a user by email (case-sensitive exact match first, then case-insensitive fallback)
func (c *UserClient) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// First try exact match
	user, err := c.getUserExactMatch(ctx, email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("exact match query failed: %w", err)
	}
	if user != nil {
		return user, nil
	}

	// Fallback to case-insensitive search if exact match fails
	user, err = c.getUserCaseInsensitive(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("case-insensitive query failed: %w", err)
	}

	return user, nil
}

// getUserExactMatch performs a case-sensitive exact match query
func (c *UserClient) getUserExactMatch(ctx context.Context, email string) (*User, error) {
	var query struct {
		Users []struct {
			ID       string `graphql:"id"`
			Email    string `graphql:"email"`
			Name     string `graphql:"name"`
			Password string `graphql:"password"`
		} `graphql:"users(where: {email: {_eq: $email}}, limit: 1)"`
	}

	err := c.client.Query(ctx, &query, map[string]interface{}{
		"email": graphql.String(email),
	})
	if err != nil {
		return nil, fmt.Errorf("graphql query failed: %w", err)
	}

	if len(query.Users) == 0 {
		return nil, ErrUserNotFound
	}

	return &User{
		ID:       query.Users[0].ID,
		Email:    query.Users[0].Email,
		Name:     query.Users[0].Name,
		Password: query.Users[0].Password,
	}, nil
}

// getUserCaseInsensitive performs a case-insensitive search
func (c *UserClient) getUserCaseInsensitive(ctx context.Context, email string) (*User, error) {
	var query struct {
		Users []struct {
			ID       string `graphql:"id"`
			Email    string `graphql:"email"`
			Name     string `graphql:"name"`
			Password string `graphql:"password"`
		} `graphql:"users(where: {email: {_ilike: $email}}, limit: 1)"`
	}

	err := c.client.Query(ctx, &query, map[string]interface{}{
		"email": graphql.String(email),
	})
	if err != nil {
		return nil, fmt.Errorf("graphql query failed: %w", err)
	}

	if len(query.Users) == 0 {
		return nil, ErrUserNotFound
	}

	return &User{
		ID:       query.Users[0].ID,
		Email:    query.Users[0].Email,
		Name:     query.Users[0].Name,
		Password: query.Users[0].Password,
	}, nil
}

func (c *UserClient) CreateUser(ctx context.Context, name, email, password string) (*User, error) {
    var mutation struct {
        InsertUser struct {
            ID    string `graphql:"id"`
            Email string
            Name  string
        } `graphql:"insert_users_one(object: {name: $name, email: $email, password: $password})"`
    }

    normalizedEmail := strings.ToLower(strings.TrimSpace(email))
    err := c.client.Mutate(ctx, &mutation, map[string]interface{}{
        "name":     graphql.String(strings.TrimSpace(name)),
        "email":    graphql.String(normalizedEmail),
        "password": graphql.String(password), // Note: password should be hashed before calling this
    })
    if err != nil {
        if strings.Contains(err.Error(), "unique constraint") {
            return nil, errors.New("user with this email already exists")
        }
        return nil, fmt.Errorf("graphql mutation failed: %w", err)
    }

    return &User{
        ID:    mutation.InsertUser.ID,
        Email: mutation.InsertUser.Email,
        Name:  mutation.InsertUser.Name,
    }, nil
}