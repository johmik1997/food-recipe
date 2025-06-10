package recipe

type Recipe struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Description      *string `json:"description"`
	PrepTime         *int    `json:"prep_time"`
	CookTime         *int    `json:"cook_time"`
	TotalTime        *int    `json:"total_time"`
	Servings         *int    `json:"servings"`
	FeaturedImageURL *string `json:"featured_image_url"`
	Price            *int    `json:"price"`
	UserID           string  `json:"user_id"`
}

type Step struct {
	ID          string  `json:"id"`
	RecipeID    string  `json:"recipe_id"`
	StepNumber  int     `json:"step_number"`
	Instruction string  `json:"instruction"`
	ImageURL    *string `json:"image_url"`
}

type Ingredient struct {
	ID       string   `json:"id"`
	RecipeID string   `json:"recipe_id"`
	Name     string   `json:"name"`
	Quantity *float64 `json:"quantity"`
	Unit     *string  `json:"unit"`
}

type RecipeCategory struct {
	RecipeID   string `json:"recipe_id"`
	CategoryID string `json:"category_id"`
}

type CreateRecipeInput struct {
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

type CreateStepInput struct {
	StepNumber  int     `json:"step_number"`
	Instruction string  `json:"instruction"`
	ImageURL    *string `json:"image_url"`
}

type CreateIngredientInput struct {
	Name     string   `json:"name"`
	Quantity *float64 `json:"quantity"`
	Unit     *string  `json:"unit"`
}