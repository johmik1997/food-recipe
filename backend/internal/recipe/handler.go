package recipe

import (
	"encoding/json"
	"foodrecipe/internal/utils"
	"net/http"
)

type RecipeHandler struct {
	recipeService RecipeService
}

func NewRecipeHandler(recipeService RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
	}
}

type createRecipeRequest struct {
	Input struct {
		Object CreateRecipeInput `json:"object"`
	} `json:"input"`
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var req createRecipeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	// Get user ID from auth context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	req.Input.Object.UserID = &userID

	recipe, err := h.recipeService.CreateRecipe(r.Context(), req.Input.Object)
	if err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(*utils.ValidationError); ok {
			status = http.StatusBadRequest
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, recipe)
}