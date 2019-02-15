package service

import (
	"strconv"

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
	strIDs, total, err := service.searchDao.SearchRecipes(search)
	if err != nil {
		return nil, err
	}
	IDs := make([]int, len(strIDs))
	for i := range strIDs {
		ID, err := strconv.Atoi(strIDs[i])
		if err != nil {
			return nil, err
		}
		IDs[i] = ID
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
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return nil, err
	}
	return service.recipeDao.GetRecipe(intID)
}

// AddRecipe adds a new recipe
func (service *RecipeService) AddRecipe(recipe model.BaseRecipe) (string, error) {
	id, err := service.recipeDao.AddRecipe(&recipe)
	if err != nil {
		return "", err
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	addedRecipe, err := service.recipeDao.GetRecipe(intID)
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
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if err = service.recipeDao.DeleteRecipe(intID); err != nil {
		return err
	}
	return service.searchDao.DeleteRecipe(id)
}
