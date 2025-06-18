package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shurcooL/graphql"
)

var validate = validator.New()

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags users Registration
// @Accept json
// @Produce json
// @Param input body requests.RegisterRequest true "User registration details"
// @Success 200 {object} response.SignedUpUserOutput
// @Failure 400 {object} gin.H "Invalid input data"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/register [post]
func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Declare a variable to hold the input data from the json req body
		var request requests.RegisterRequest

		// Bind the JSON request body to the request struct
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("Validation error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data"})
			return
		}

		// Validate the input data
		if err := validate.Struct(request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "details": err.Error()})
			return
		}

		// Check if the user already exists
		var query struct {
			User []struct {
				ID       graphql.Int    `graphql:"id"`
				Name     graphql.String `graphql:"username"`
				Email    graphql.String `graphql:"email"`
				Password graphql.String `graphql:"password"`
				Role     graphql.String `graphql:"role"`
			} `graphql:"users(where: {email: {_eq: $email}})"`
		}

		variable := map[string]interface{}{
			"email": graphql.String(request.Input.Email),
		}

		if err := client.Query(ctx, &query, variable); err != nil {
			log.Println("Error querying existing user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query the existing user", "details": err.Error()})
			return
		}

		if len(query.User) != 0 {
			log.Println("User already exists")
			c.JSON(http.StatusBadRequest, gin.H{"message": "User already exists"})
			return

		}

		imageUrl, exists := c.Get("imageUrl")
		if !exists {
			imageUrl = ""
		}

		profilePictureURL, ok := imageUrl.(string)
		if !ok {
			profilePictureURL = ""
		}
		// hash the password
		password := helpers.HashPassword(request.Input.Password)

		var mutation struct {
			CreateUser struct {
				ID       graphql.Int    `graphql:"id"`
				UserName graphql.String `graphql:"username"`
				Email    graphql.String `graphql:"email"`
				Profile  graphql.String `graphql:"profile"`
				Role     graphql.String `graphql:"role"`
			} `graphql:"insert_users_one(object: {username: $userName, email: $email, password: $password, profile: $profile, role:$role, phone:$phone})"`
		}
		mutationVariables := map[string]interface{}{
			"userName": graphql.String(request.Input.UserName),
			"email":    graphql.String(request.Input.Email),
			"password": graphql.String(password),
			"phone":    graphql.String(request.Input.Phone),
			"profile":  graphql.String(profilePictureURL),
			"role":     graphql.String(request.Input.Role + "user"),
		}
		err := client.Mutate(context.Background(), &mutation, mutationVariables)
		if err != nil {
			log.Println("Failed to register user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user"})
			return
		}
		// query newly created user
		var regesteredUserQuery struct {
			User []struct {
				ID       graphql.Int    `graphql:"id"`
				Name     graphql.String `graphql:"username"`
				Email    graphql.String `graphql:"email"`
				Password graphql.String `graphql:"password"`
				Role     graphql.String `graphql:"role"`
				TokenId  graphql.String `graphql:"tokenId"`
			} `graphql:"users(where: {email: {_eq: $email}})"`
		}

		regesteredUserVariable := map[string]interface{}{
			"email": graphql.String(request.Input.Email),
		}
		err = client.Query(context.Background(), &regesteredUserQuery, regesteredUserVariable)
		if err != nil {
			log.Println("Error querying the registered user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query the registered user"})
			return
		}
		if len(regesteredUserQuery.User) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "user not found after registration maybe he flead to another table:)"})
		}
		user := regesteredUserQuery.User[0]
		// generate email verification token
		emailVerficationToken, err := helpers.GenerateToken()

		if err != nil {
			log.Println("Error generating email verification token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate email verification token", "detail": err.Error()})
			return
		}

		var VerficationEmailMutation struct {
			InsertedToken struct {
				ID    graphql.Int    `graphql:"id"`
				Token graphql.String `graphql:"token"`
			} `graphql:"insert_email_verification_tokens_one(object: {token: $token, user_id: $userId})"`
		}

		// "insert_email_verification_tokens_one(object: {token: $token, user_id: $userId})"

		verficatonEmailVariable := map[string]interface{}{
			"token":  graphql.String(emailVerficationToken),
			"userId": graphql.Int(user.ID),
		}

		err = client.Mutate(context.Background(), &VerficationEmailMutation, verficatonEmailVariable)
		if err != nil {
			log.Printf("Error inserting email verification token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register verfication token"})
			return
		}

		var UpdateUserTokenMutation struct {
			UpdatedUser struct {
				ID      graphql.Int `graphql:"id"`
				TokenID graphql.Int `graphql:"tokenId"`
			} `graphql:"update_users_by_pk(pk_columns: {id: $userId}, _set: {tokenId: $tokenId})"`
		}

		updatedUserVariables := map[string]interface{}{
			"userId":  graphql.Int(user.ID),
			"tokenId": VerficationEmailMutation.InsertedToken.ID,
		}

		err = client.Mutate(context.Background(), &UpdateUserTokenMutation, updatedUserVariables)

		if err != nil {
			log.Printf("error updating user with tokenId: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user tokenid", "details": err.Error()})
			return
		}
		log.Printf("User %d successfully updated with tokenId %d", user.ID, UpdateUserTokenMutation.UpdatedUser.TokenID)

		verificationLink := os.Getenv("RETURN_URL") + "?verification_token=" + emailVerficationToken + "&user_id=" + strconv.Itoa(int(user.ID))

		emailForm := helpers.EmailData{
			Name:    string(user.Name),
			Email:   string(user.Email),
			Link:    verificationLink,
			Subject: "Verifying your email",
		}

		res, errString := helpers.SendEmail(
			[]string{emailForm.Email},
			"verifyEmail.html",
			emailForm,
		)
		if !res {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send email verification email, please contact support.", "details": errString})
			return
		}

		token, refreshToken, err := helpers.GenerateAllTokens(string(user.Email), string(user.Name), string(user.Role), fmt.Sprintf("%d", user.ID), int(user.ID))
		fmt.Println("token id", user.Email, user.Name, user.Role, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to  generate jwt token", "details": err.Error()})
			return

		}

		response := response.SignedUpUserOutput{
			ID:           int(user.ID),
			UserName:     string(user.Name),
			Email:        string(user.Email),
			Token:        token,
			RefreshToken: refreshToken,
			Role:         string(user.Role),
		}
		c.JSON(http.StatusOK, response)

	}

}
