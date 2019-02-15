package service

import (
	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/model"
)

// RecipeService struct
type RecipeService struct {
	recipeDao *datasource.RecipeDao
	searchDao *datasource.RecipeSearchDao
}

// newRecipeService creates a new recipe service
func newRecipeService(recipeDao *datasource.RecipeDao, searchDao *datasource.RecipeSearchDao) *RecipeService {
	return &RecipeService{
		recipeDao,
		searchDao,
	}
}

// SearchRecipe searches for recipes
func (service *RecipeService) SearchRecipe(search model.RecipeSearch) (*model.RecipeSearchResult, error) {
	if search.IsEmpty() {
		return service.recipeDao.GetRandomRecipes(search)
	}
	IDs, total, err := service.searchDao.SearchRecipes(search)
	if err != nil {
		return nil, err
	}
	recipes, err := service.recipeDao.GetRecipes(IDs)
	if err != nil {
		return nil, err
	}
	return &model.RecipeSearchResult{
		FirstResults: recipes,
		Total:        total,
	}, nil
}

// GetRecipe gets a recipe by its ID
func (service *RecipeService) GetRecipe(ID string) (*model.Recipe, error) {
	return service.recipeDao.GetRecipe(ID)
}

// AddRecipe adds a new recipe
func (service *RecipeService) AddRecipe(recipe model.BaseRecipe) (string, error) {
	id, err := service.recipeDao.AddRecipe(&recipe)
	if err != nil {
		return "", err
	}
	addedRecipe, err := service.recipeDao.GetRecipe(id)
	if err != nil {
		return "", err
	}
	if err = service.searchDao.IndexRecipe(*addedRecipe); err != nil {
		return "", err
	}
	return addedRecipe.ID, nil
}

// UpdateRecipe updates an existing recipe
func (service *RecipeService) UpdateRecipe(ID string, recipe model.BaseRecipe) (*model.Recipe, error) {
	updated, err := service.recipeDao.UpdateRecipe(model.Recipe{
		ID:         ID,
		BaseRecipe: recipe,
	})
	if err != nil {
		return nil, err
	}
	if err = service.searchDao.IndexRecipe(*updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// DeleteRecipe deletes a recipe
func (service *RecipeService) DeleteRecipe(id string) error {
	if err := service.recipeDao.DeleteRecipe(id); err != nil {
		return err
	}
	return service.searchDao.DeleteRecipe(id)
}
