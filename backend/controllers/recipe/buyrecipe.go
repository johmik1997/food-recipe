package recipe

import (
	"context"
	"fmt"
	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)

func BuyRecipe() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var req requests.BuyRecipeRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Println("comming request from frontend", req)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid json input", "details": err.Error()})
			println(err.Error())
			return
		}
		validate := validator.New()
		if validationError := validate.Struct(req); validationError != nil {
			log.Println("validation failed:", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": validationError.Error()})
			return
		}

		var userQuery struct {
			User struct {
				Name  string `graphql:"username"`
				Phone string `graphql:"phone"`
				Email string `graphql:"email"`
			} `graphql:"users_by_pk(id: $id)"`
		}

		userQueryVars := map[string]interface{}{
			"id": graphql.Int(req.Input.BuyerId),
		}

		if err := client.Query(ctx, &userQuery, userQueryVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query user data aye", "details": err.Error()})
			return
		}
		log.Println("user data after query", userQuery.User)

		var RecipeQuery struct {
			Recipes []struct {
				ID    int `graphql:"id"`
				Price int `graphql:"price"`
				User  struct {
				ID int `graphql:"id"`
				}
			} `graphql:"recipes(where: {id: {_eq: $recipeId}})"`
		}
		var recipeQueryVars = map[string]interface{}{
			"recipeId": graphql.Int(req.Input.RecipeId),
		}
		if err := client.Query(ctx, &RecipeQuery, recipeQueryVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query recipe data", "details": err.Error()})
			return
		}
		log.Println("recipe data after query", RecipeQuery.Recipes)
		if len(RecipeQuery.Recipes) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Recipe not found"})
			return
		}

		Price := RecipeQuery.Recipes[0].Price

		// check if the user already bought the recipe
		var soldRecipeQuery struct {
			SoldRecipes []struct {
				ID       int `graphql:"id"`
				BuyerId  int `graphql:"buyer_id"`
				RecipeId int `graphql:"recipe_id"`
				SellerId int `graphql:"seller_id"`
			} `graphql:"sold_recipes(where: {buyer_id: {_eq: $buyer_id}, recipe_id: {_eq: $recipe_id}})"`
		}
		var soldRecipeQueryVars = map[string]interface{}{
			"buyer_id":  graphql.Int(req.Input.BuyerId),
			"recipe_id": graphql.Int(req.Input.RecipeId),
		}
		if err := client.Query(ctx, &soldRecipeQuery, soldRecipeQueryVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query sold recipe data", "details": err.Error()})
			return
		}
		log.Println("sold recipe data after query", soldRecipeQuery.SoldRecipes)
		if len(soldRecipeQuery.SoldRecipes) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "you already bought the recipe"})
			return
		}

		var muation struct {
			BuyRecipe struct {
				ID       graphql.Int `graphql:"id"`
				BuyerId  graphql.Int `graphql:"buyer_id"`
				RecipeId graphql.Int `graphql:"recipe_id"`
				Price    graphql.Int `graphql:"price"`
				SellerId graphql.Int `graphql:"seller_id"`
			} `graphql:"insert_sold_recipes_one(object: {buyer_id: $buyer_id, price: $price, recipe_id: $recipe_id, seller_id: $seller_id})"`
		}
		var mutationVars = map[string]interface{}{
			"buyer_id":  graphql.Int(req.Input.BuyerId),
			"recipe_id": graphql.Int(req.Input.RecipeId),
			"price":     graphql.Int(Price),
			"seller_id": graphql.Int(RecipeQuery.Recipes[0].User.ID),
		}
		log.Println("a user with id " + strconv.Itoa(int(req.Input.BuyerId)) + " is trying to buy a recipe " + strconv.Itoa(int(req.Input.RecipeId)))
		if err := client.Mutate(ctx, &muation, mutationVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to buy recipe", "details": err.Error()})
			return
		}

		var chapaForm helpers.PaymentRequest
		chapaForm.Amount = Price

		chapaForm.PhoneNumber = userQuery.User.Phone
		chapaForm.Email = userQuery.User.Email
		chapaForm.FirstName = userQuery.User.Name
		chapaForm.Currency = "ETB"
		chapaForm.TxRef = fmt.Sprintf("purchesed-recipe-%d", muation.BuyRecipe.ID)
		// chapaForm.ReturnURL = os.Getenv("CHAPA_REDIRECT_URL")
		chapaForm.ReturnURL = fmt.Sprintf("%s/%d", os.Getenv("CHAPA_REDIRECT_URL"), muation.BuyRecipe.RecipeId)
		chapaForm.CallbackURL = os.Getenv("CHAPA_CALLBACK_URL")
		chapaForm.CustomizationTitle = "Buying a recipe"
		chapaForm.CustomizationDesc = "Buying a recipe"

		log.Println("recipe bought successfully", muation.BuyRecipe)

		paymentResponse, err := helpers.InitPayment(&chapaForm)
		fmt.Println("chapa form to get the user detail..", chapaForm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to buy recipe", "details": err.Error()})
			return
		}
		log.Println("payment initiated successfully", paymentResponse)
		log.Println("ChapaResponseeeee: ", paymentResponse.ChapaResponse)
		

		var paymentID graphql.Int
		var checkoutUrl graphql.String

		if paymentResponse.Status {
			data, ok := paymentResponse.ChapaResponse["data"].(map[string]interface{})
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve checkout_url from ChapaResponse"})
				return
			}

			checkoutURL, ok := data["checkout_url"].(string)
			if !ok {
				log.Panicln("checkout url not found in ChapaResponse", data)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve checkout_url from data"})
				return
			}

			log.Println("checkout url: ", checkoutURL)

			paymentMutation := struct {
				InsertPayment struct {
					ID            graphql.Int `json:"id"`
					SoldRecipeId  graphql.Int `graphql:"sold_recipe_id"`
					TxRef         string      `graphql:"tx_ref"`
					CheckoutURL   string      `graphql:"checkout_url"`
					Amount        graphql.Int `graphql:"amount"`
					Currency      string      `graphql:"currency"`
					PaymentMethod string      `graphql:"payment_method"`
					Status        string      `graphql:"payment_status"`
					BuyerId       graphql.Int `graphql:"buyer_id"`
					RecipeId      graphql.Int `graphql:"recipe_id"`
				} `graphql:"insert_payments_one(object: {sold_recipe_id: $soldRecipeId, tx_ref: $txRef, checkout_url: $checkoutUrl, amount: $amount, currency: $currency, payment_method: $paymentMethod, payment_status: $status, buyer_id: $buyerId, recipe_id: $recipeId})"`
			}{}

			paymentVars := map[string]interface{}{
				"soldRecipeId":  graphql.Int(muation.BuyRecipe.ID),
				"txRef":         graphql.String(chapaForm.TxRef),
				"checkoutUrl":   graphql.String(checkoutURL),
				"amount":        graphql.Int(Price),
				"currency":      graphql.String(chapaForm.Currency),
				"paymentMethod": graphql.String("chapa"),
				"status":        graphql.String("pending"),
				"buyerId":       graphql.Int(req.Input.BuyerId),
				"recipeId":      graphql.Int(req.Input.RecipeId),
			}

			if err := client.Mutate(ctx, &paymentMutation, paymentVars); err != nil {
				log.Println("Failed to insert payment record:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to insert payment record", "details": err.Error()})
				return
			}
			log.Println("payment inserted successfully", paymentMutation.InsertPayment)
			paymentID = paymentMutation.InsertPayment.ID
			checkoutUrl = graphql.String(paymentMutation.InsertPayment.CheckoutURL)
			log.Println("checccccclkout url", paymentMutation.InsertPayment.CheckoutURL)
			log.Println("Payment record inserted successfully and checkout url:.", checkoutURL)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "paymnent initation failed", "details": paymentResponse.Message})
		}
		res := response.BuyRecipeOutput{
			Message:     "recipe bought successfully",
			PaymentId:   graphql.Int(paymentID),
			CheckOutUrl: checkoutUrl,
			BuyerId:     graphql.Int(muation.BuyRecipe.BuyerId),
			RecipeId:    graphql.Int(muation.BuyRecipe.RecipeId),
			Price:       graphql.Int(Price),
			SellerId:    graphql.Int(muation.BuyRecipe.SellerId),
		}
		log.Println("Returning response:", res)
		log.Printf("Response object before sending: %+v\n", res)

		c.JSON(http.StatusOK, res)
	}
}
