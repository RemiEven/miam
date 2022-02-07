package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/model"
)

// RecipeService struct
type RecipeService struct {
	recipeDao *datasource.RecipeDao
	searchDao *datasource.RecipeSearchDao
}

// NewRecipeService creates a new recipe service
func NewRecipeService(recipeDao *datasource.RecipeDao, searchDao *datasource.RecipeSearchDao) *RecipeService {
	return &RecipeService{
		recipeDao,
		searchDao,
	}
}

// IndexAllExistingRecipes lists all recipes that are in the database and index them so that they are searchable
func (service *RecipeService) IndexAllExistingRecipes(ctx context.Context) error {
	ids, err := service.recipeDao.ListRecipeIds(ctx)
	if err != nil {
		return fmt.Errorf("failed to list recipes: %w", err)
	}

	for _, id := range ids {
		recipe, err := service.recipeDao.GetRecipe(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get recipe for id [%s]: %w", id, err)
		}
		if recipe == nil {
			continue
		}
		if err := service.searchDao.IndexRecipe(*recipe); err != nil {
			return fmt.Errorf("failed to index recipe with id [%s]: %w", id, err)
		}
		log.Debug().Str("id", id).Msg("indexed recipe")
	}

	return nil
}

// SearchRecipe searches for recipes
func (service *RecipeService) SearchRecipe(ctx context.Context, search model.RecipeSearch) (*model.RecipeSearchResult, error) {
	if search.IsEmpty() {
		return service.recipeDao.GetRandomRecipes(ctx, search)
	}
	IDs, total, err := service.searchDao.SearchRecipes(search)
	if err != nil {
		return nil, err
	}
	recipes, err := service.recipeDao.GetRecipes(ctx, IDs)
	if err != nil {
		return nil, fmt.Errorf("failed to hydrate matching recipes: *%w", err)
	}
	return &model.RecipeSearchResult{
		FirstResults: recipes,
		Total:        total,
	}, nil
}

// GetRecipe gets a recipe by its ID
func (service *RecipeService) GetRecipe(ctx context.Context, ID string) (*model.Recipe, error) {
	return service.recipeDao.GetRecipe(ctx, ID)
}

// AddRecipe adds a new recipe
func (service *RecipeService) AddRecipe(ctx context.Context, recipe model.BaseRecipe) (string, error) {
	id, err := service.recipeDao.AddRecipe(ctx, &recipe)
	if err != nil {
		return "", err
	}
	addedRecipe, err := service.recipeDao.GetRecipe(ctx, id)
	if err != nil {
		return "", err
	}
	if err = service.searchDao.IndexRecipe(*addedRecipe); err != nil {
		return "", err
	}
	return addedRecipe.ID, nil
}

// UpdateRecipe updates an existing recipe
func (service *RecipeService) UpdateRecipe(ctx context.Context, ID string, recipe model.BaseRecipe) (*model.Recipe, error) {
	updated, err := service.recipeDao.UpdateRecipe(ctx, model.Recipe{
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
func (service *RecipeService) DeleteRecipe(ctx context.Context, id string) error {
	if err := service.recipeDao.DeleteRecipe(ctx, id); err != nil {
		return err
	}
	return service.searchDao.DeleteRecipe(id)
}
