package model

// Recipe is a recipe with id/name
type Recipe struct {
	BaseRecipe `json:""`
	ID         string `json:"id"`
}

// BaseRecipe is an editable recipe
type BaseRecipe struct {
	Name        string             `json:"name"`
	HowTo       string             `json:"howTo,omitempty"`
	Ingredients []RecipeIngredient `json:"ingredients,omitempty"`
}

// RecipeIngredient is a recipe ingredient with an optional quantity
type RecipeIngredient struct {
	Quantity   string `json:"quantity,omitempty"`
	Ingredient `json:""`
}
