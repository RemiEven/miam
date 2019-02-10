package model

// RecipeSearch is a search for a recipe
type RecipeSearch struct {
	SearchTerm          string   `json:"searchTerm,omitempty"`
	ExcludedRecipes     []string `json:"excludedRecipes,omitempty"`
	ExcludedIngredients []string `json:"excludedIngredients,omitempty"`
}

// RecipeSearchResult is the result of a recipe search
type RecipeSearchResult struct {
	Total        int      `json:"total"`
	FirstResults []Recipe `json:"firstResults"`
}
