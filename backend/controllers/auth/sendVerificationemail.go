package auth
import (
	"fmt"
	"os"
	"strconv"
	"foodrecipe/controllers/helpers"
	
)
func sendVerificationEmail(userID int, name, email, token string) error {
	verificationLink := os.Getenv("VERIFICATION_URL") + "?verification_token=" + token + "&user_id=" + strconv.Itoa(userID)

	emailForm := helpers.EmailData{
		Name:    name,
		Email:   email,
		Link:    verificationLink,
		Subject: "Email Verification",
	}

	res, errMsg := helpers.SendEmail(
		[]string{emailForm.Email},
		"verifyEmail.html",
		emailForm,
	)
	if !res {
		return fmt.Errorf("failed to send email: %s", errMsg)
	}

	return nil
}