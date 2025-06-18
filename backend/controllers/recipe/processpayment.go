package recipe

import (
	"context"

	"foodrecipe/controllers/libs"

	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)


func ProcessPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var req requests.PaymentProcessRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}
		log.Println("tx_ref_id:", req.Input.Id)

		isVerified, err := helpers.VerifyPayment(req.Input.TxRef)
		if err != nil || !isVerified {
			log.Println("Payment verification failed:", err)
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
			} `graphql:"update_payments(where: {id: {_eq: $id}}, _set: {payment_status: \"paid\"})"`
		}

		mutationVars := map[string]interface{}{
			"id": graphql.Int(req.Input.Id),
		}

		err = client.Mutate(ctx, &mutation, mutationVars)
		if err != nil {
			log.Println("Failed to update payment status:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update payment status"})
			return
		}

		log.Println("Payment status updated successfully for Payment ID:", req.Input.Id)
		log.Println("Payment status updated successfully for Payment txref:", req.Input.TxRef)
		res := response.ProcessPaymentOutput{
			Message: "payment processed successfully and status is updated",
			Status:  "success",
		}
		c.JSON(http.StatusOK, res)
	}
}
