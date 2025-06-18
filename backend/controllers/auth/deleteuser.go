package auth

import (
	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/requests"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var eventPayload requests.EventPayload

		if err := c.ShouldBindJSON(&eventPayload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		if eventPayload.Event.Op == "DELETE" && eventPayload.Event.Data.Old != nil {
			oldUserData := eventPayload.Event.Data.Old

			emailData := helpers.EmailData{
				Name:    oldUserData.UserName,
				Email:   oldUserData.Email,
				Subject: "Account Deletion Confirmation",
			}

			success, errorString := helpers.SendEmail(
				[]string{emailData.Email},
				"deletedAccount.html",
				emailData,
			)

			if !success {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send account deletion email", "details": errorString})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "Account removal email has been sent", "email": oldUserData.Email})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": "Something went wrong", "details": "Invalid input"})
	}
}
