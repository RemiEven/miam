package model

// Ingredient is an ingredient with id/name
type Ingredient struct {
	BaseIngredient `json:""`
	ID             string `json:"id"`
}

// BaseIngredient is an editable ingredient
type BaseIngredient struct {
	Name string `json:"name"`
}
