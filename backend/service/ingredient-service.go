package service

import (
	"strconv"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/model"
)

// IngredientService struct
type IngredientService struct {
	ingredientDao *datasource.IngredientDao
}

// newIngredientService creates a new ingredient service
func newIngredientService(ingredientDao *datasource.IngredientDao) *IngredientService {
	return &IngredientService{
		ingredientDao,
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
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}
	return service.ingredientDao.DeleteIngredient(intID)
}
