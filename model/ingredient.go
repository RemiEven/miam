package model

type Ingredient struct {
	BaseIngredient `json:""`
	ID             string `json:"id"`
}

type BaseIngredient struct {
	Name string `json:"name"`
}
