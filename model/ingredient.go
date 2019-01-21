package model

type Ingredient struct {
	BaseIngredient `json:""`
	ID             string
}

type BaseIngredient struct {
	Name string `json:"name"`
}
