package service

import (
	"context"

	"github.com/RemiEven/miam/common"
	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/model"
)

// IngredientService struct
type IngredientService struct {
	ingredientDao       *datasource.IngredientDao
	recipeIngredientDao *datasource.RecipeIngredientDao
}

// NewIngredientService creates a new ingredient service
func NewIngredientService(ingredientDao *datasource.IngredientDao, recipeIngredientDao *datasource.RecipeIngredientDao) *IngredientService {
	return &IngredientService{
		ingredientDao,
		recipeIngredientDao,
	}
}

// GetAllIngredients returns all known ingredients
func (service *IngredientService) GetAllIngredients(ctx context.Context) ([]model.Ingredient, error) {
	return service.ingredientDao.GetAllIngredients(ctx)
}

// UpdateIngredient updates an ingredient
func (service *IngredientService) UpdateIngredient(ctx context.Context, ID string, update model.BaseIngredient) (*model.Ingredient, error) {
	ingredient := model.Ingredient{
		ID:             ID,
		BaseIngredient: update,
	}
	if err := service.ingredientDao.UpdateIngredient(ctx, ingredient); err != nil {
		return nil, err
	}
	return &ingredient, nil
}

// DeleteIngredient deletes the ingredient with the given id
func (service *IngredientService) DeleteIngredient(ctx context.Context, ID string) error {
	used, err := service.recipeIngredientDao.IsUsedInRecipe(ctx, ID)
	if err != nil {
		return err
	}
	if used {
		return common.ErrInvalidOperation
	}
	return service.ingredientDao.DeleteIngredient(ctx, ID)
}
