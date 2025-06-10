package recipe

import (
	"encoding/json"
	"foodrecipe/internal/middleware"
	"foodrecipe/internal/utils"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRecipeRequest struct {
	Title       string                  `json:"title"`
	Description *string                 `json:"description"`
	PrepTime    *int                    `json:"prep_time"`
	CookTime    *int                    `json:"cook_time"`
	Servings    *int                    `json:"servings"`
	Price       *int                    `json:"price"`
	CategoryIDs []string                `json:"category_ids"`
	Steps       []CreateStepInput       `json:"steps"`
	Ingredients []CreateIngredientInput `json:"ingredients"`
}

func (h *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok || userID == "" {
		msg := "missing user authentication"
		log.Println("[RecipeHandler] " + msg)
		utils.WriteJSONError(w, http.StatusUnauthorized, msg)
		return
	}
	log.Printf("[RecipeHandler] Authenticated user: %s", userID)


	var req createRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := CreateRecipeInput{
		Title:       req.Title,
		Description: req.Description,
		PrepTime:    req.PrepTime,
		CookTime:    req.CookTime,
		Servings:    req.Servings,
		Price:       req.Price,
		CategoryIDs: req.CategoryIDs,
		Steps:       req.Steps,
		Ingredients: req.Ingredients,
	}

	recipe, err := h.service.CreateRecipe(ctx, input, userID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, recipe)
}