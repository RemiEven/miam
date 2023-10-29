package service

import (
	"context"
	"fmt"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/model"
	"github.com/RemiEven/miam/pb-lite/failure"
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
		return nil, fmt.Errorf("failed to update ingredient: %w", err)
	}
	return &ingredient, nil
}

// DeleteIngredient deletes the ingredient with the given id
func (service *IngredientService) DeleteIngredient(ctx context.Context, ID string) error {
	used, err := service.recipeIngredientDao.IsUsedInRecipe(ctx, ID)
	if err != nil {
		return fmt.Errorf("failed to check if ingredient is used in recipe: %w", err)
	}
	if used {
		return &failure.InvalidValueError{
			Message: "cannot delete ingredient [" + ID + "] because it's still used in some recipes",
		}
	}

	if err := service.ingredientDao.DeleteIngredient(ctx, ID); err != nil {
		return fmt.Errorf("failed to delete ingredient: %w", err)
	}

	return nil
}
