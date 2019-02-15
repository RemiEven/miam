package service

import (
	"github.com/RemiEven/miam/common"
	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/model"
)

// IngredientService struct
type IngredientService struct {
	ingredientDao       *datasource.IngredientDao
	recipeIngredientDao *datasource.RecipeIngredientDao
}

// newIngredientService creates a new ingredient service
func newIngredientService(ingredientDao *datasource.IngredientDao, recipeIngredientDao *datasource.RecipeIngredientDao) *IngredientService {
	return &IngredientService{
		ingredientDao,
		recipeIngredientDao,
	}
}

// GetAllIngredients returns all known ingredients
func (service *IngredientService) GetAllIngredients() ([]model.Ingredient, error) {
	return service.ingredientDao.GetAllIngredients()
}

// UpdateIngredient updates an ingredient
func (service *IngredientService) UpdateIngredient(ID string, update model.BaseIngredient) (*model.Ingredient, error) {
	ingredient := model.Ingredient{
		ID:             ID,
		BaseIngredient: update,
	}
	if err := service.ingredientDao.UpdateIngredient(ingredient); err != nil {
		return nil, err
	}
	return &ingredient, nil
}

// DeleteIngredient deletes the ingredient with the given id
func (service *IngredientService) DeleteIngredient(ID string) error {
	used, err := service.recipeIngredientDao.IsUsedInRecipe(ID)
	if err != nil {
		return err
	}
	if used {
		return common.ErrInvalidOperation
	}
	return service.ingredientDao.DeleteIngredient(ID)
}
