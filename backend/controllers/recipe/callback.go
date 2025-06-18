package recipe

import (
	"context"

	"foodrecipe/controllers/libs"

	"foodrecipe/controllers/helpers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)

func PaymentCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()
		// Extract query parameters
		txRef := c.Query("trx_ref")
		status := c.Query("status")

		log.Println("Received payment callback - TxRef:", txRef, "Status:", status)

		if status != "success" {
			log.Println("Payment failed or not completed:", status)
			c.JSON(http.StatusOK, gin.H{"message": "payment not successful"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		isVerified, err := helpers.VerifyPayment(txRef)
		if err != nil || !isVerified {
			log.Println("Payment verification failed from backend callback:", err)
			c.JSON(http.StatusOK, gin.H{"message": "payment verification failed", "details": err.Error()})
			return
		}

		type UpdatePaymentMutation struct {
			Status graphql.String `graphql:"payment_status"`
			TxRef  graphql.String `graphql:"tx_ref"`
		}

		var mutation struct {
			UpdatePayment struct {
				Returning []UpdatePaymentMutation `graphql:"returning"`
			} `graphql:"update_payments(where: {tx_ref: {_eq: $txRef}}, _set: {payment_status: \"paid\"})"`
		}
		mutationVars := map[string]interface{}{
			"txRef": graphql.String(txRef),
		}
		err = client.Mutate(ctx, &mutation, mutationVars)
		if err != nil {
			log.Println("Failed to update payment status:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update payment status"})
			return
		}

		log.Println("Payment status updated successfully for TxRef:", txRef)
		c.JSON(http.StatusOK, gin.H{"message": "payment processed successfully"})
	}
}