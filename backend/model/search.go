package model

// RecipeSearch is a search for a recipe
type RecipeSearch struct {
	SearchTerm          string   `json:"searchTerm,omitempty"`
	ExcludedRecipes     []string `json:"excludedRecipes,omitempty"`
	ExcludedIngredients []string `json:"excludedIngredients,omitempty"`
}

// IsEmpty returns true if the search contains no criteria
func (search RecipeSearch) IsEmpty() bool {
	return len(search.SearchTerm) == 0 &&
		(search.ExcludedRecipes == nil || len(search.ExcludedRecipes) == 0) &&
		(search.ExcludedIngredients == nil || len(search.ExcludedIngredients) == 0)
}

// RecipeSearchResult is the result of a recipe search
type RecipeSearchResult struct {
	Total        int      `json:"total"`
	FirstResults []Recipe `json:"firstResults"`
}
