package datasource

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RemiEven/miam/model"
	"github.com/RemiEven/miam/pb-lite/failure"
)

// RecipeIngredientDao struct
type RecipeIngredientDao struct {
	holder        *DatabaseHolder
	ingredientDao *IngredientDao
}

// NewRecipeIngredientDao returns a new recipe ingredient dao
func NewRecipeIngredientDao(holder *DatabaseHolder, ingredientDao *IngredientDao) (*RecipeIngredientDao, error) {
	initStatement := `
		create table if not exists recipe_ingredient (recipe_id int, ingredient_id int, quantity text);
		create index if not exists recipe_id_index on recipe_ingredient(recipe_id);
		create index if not exists ingredient_id_index on recipe_ingredient(ingredient_id);
	`
	if _, err := holder.DB.Exec(initStatement); err != nil {
		return nil, fmt.Errorf("failed to create recipe_ingredient table and/or its indices: %w", err)
	}
	return &RecipeIngredientDao{holder, ingredientDao}, nil
}

// GetRecipeIngredients returns the ingredient of a recipe
func (dao *RecipeIngredientDao) GetRecipeIngredients(ctx context.Context, recipeID string) ([]model.RecipeIngredient, error) {
	intRecipeID, err := toSqliteID(recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert [%s] to sqlite ID: %w", recipeID, err)
	}
	rows, err := dao.holder.DB.QueryContext(ctx, `select
		recipe_ingredient.ingredient_id, recipe_ingredient.quantity, ingredient.name
		from recipe_ingredient
		inner join ingredient
		on recipe_ingredient.ingredient_id=ingredient.id
		where recipe_ingredient.recipe_id=?`, intRecipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipe ingredients: %w", err)
	}
	defer rows.Close()
	recipeIngredients := make([]model.RecipeIngredient, 0, 8) // most recipes have 8 or less ingredients
	for rows.Next() {
		var ingredientID int
		var quantity string
		var name string
		if err := rows.Scan(&ingredientID, &quantity, &name); err != nil {
			return nil, fmt.Errorf("failed to scan recipe ingredient row: %w", err)
		}
		recipeIngredients = append(recipeIngredients, model.RecipeIngredient{
			Ingredient: model.Ingredient{
				ID: fromSqliteID(ingredientID),
				BaseIngredient: model.BaseIngredient{
					Name: name,
				},
			},
			Quantity: quantity,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("got an error while iterating on recipe ingredients rows: %w", err)
	}
	return recipeIngredients, nil
}

// AddRecipeIngredient adds a recipe ingredient.
// If the ingredient is new (ie. it has not yet got an ID) it is also added.
func (dao *RecipeIngredientDao) AddRecipeIngredient(ctx context.Context, transaction *sql.Tx, recipeID string, recipeIngredient model.RecipeIngredient) (string, error) {
	if len(recipeIngredient.ID) == 0 {
		// New ingredient
		ingredientID, err := dao.ingredientDao.AddIngredient(ctx, transaction, recipeIngredient.Name)
		if err != nil {
			return "", fmt.Errorf("failed to add new ingredient: %w", err)
		}
		recipeIngredient.ID = ingredientID
	}
	intRecipeID, err := toSqliteID(recipeID)
	if err != nil {
		return "", fmt.Errorf("failed to convert [%s] to sqlite ID: %w", recipeID, err)
	}
	intIngredientID, err := toSqliteID(recipeIngredient.ID)
	if err != nil {
		return "", fmt.Errorf("failed to convert [%s] to sqlite ID: %w", recipeIngredient.ID, err)
	}

	insertStatement, err := transaction.PrepareContext(ctx, "insert into recipe_ingredient(recipe_id, ingredient_id, quantity) values(?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer insertStatement.Close()

	_, err = insertStatement.ExecContext(ctx, intRecipeID, intIngredientID, recipeIngredient.Quantity)
	if err != nil {
		return "", fmt.Errorf("failed to execute insert statement: %w", err)
	}
	return recipeIngredient.ID, nil
}

// IsUsedInRecipe returns whether an ingredient is used by at least one recipe
func (dao *RecipeIngredientDao) IsUsedInRecipe(ctx context.Context, ingredientID string) (bool, error) {
	intID, err := toSqliteID(ingredientID)
	if err != nil {
		return false, fmt.Errorf("failed to convert [%s] to sqlite ID: %w", ingredientID, err)
	}
	row := dao.holder.DB.QueryRowContext(ctx, "select exists(select 1 from recipe_ingredient where ingredient_id=?)", intID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to scan exist row: %w", err)
	}
	return exists, nil
}

// DeleteRecipeIngredients deletes the ingredients of a recipe
func (dao *RecipeIngredientDao) DeleteRecipeIngredients(ctx context.Context, transaction *sql.Tx, recipeID string) error {
	intID, err := toSqliteID(recipeID)
	if err != nil {
		return fmt.Errorf("failed to convert [%s] to sqlite ID: %w", recipeID, err)
	}
	deleteStatement, err := transaction.PrepareContext(ctx, "delete from recipe_ingredient where recipe_id=?")
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer deleteStatement.Close()

	if _, err := deleteStatement.ExecContext(ctx, intID); err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}
	return nil
}

// DeleteRecipeIngredient removes one ingredient from a recipe
func (dao *RecipeIngredientDao) DeleteRecipeIngredient(ctx context.Context, transaction *sql.Tx, recipeID, ingredientID string) error {
	intRecipeID, err := toSqliteID(recipeID)
	if err != nil {
		return fmt.Errorf("failed to convert [%s] to sqlite ID: %w", recipeID, err)
	}
	intIngredientID, err := toSqliteID(ingredientID)
	if err != nil {
		return fmt.Errorf("failed to convert [%s] to sqlite ID: %w", ingredientID, err)
	}
	deleteStatement, err := transaction.PrepareContext(ctx, "delete from recipe_ingredient where recipe_id=? and ingredient_id=?")
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer deleteStatement.Close()

	if _, err := deleteStatement.Exec(ctx, intRecipeID, intIngredientID); err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}
	return nil
}

// UpdateRecipeIngredient updates an ingredient of a recipe
func (dao *RecipeIngredientDao) UpdateRecipeIngredient(ctx context.Context, transaction *sql.Tx, recipeID string, recipeIngredient model.RecipeIngredient) error {
	intID, err := toSqliteID(recipeID)
	if err != nil {
		return fmt.Errorf("failed to convert [%s] to sqlite ID: %w", recipeID, err)
	}
	updateStatement, err := transaction.PrepareContext(ctx, "update recipe_ingredient set quantity=?3 where recipe_id=?1 and ingredient_id=?2")
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer updateStatement.Close()

	result, err := updateStatement.ExecContext(ctx, intID, recipeIngredient.ID, recipeIngredient.Quantity)
	if err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	switch {
	case err != nil:
		return fmt.Errorf("failed to retrieve number of rows affected by update statement: %w", err)
	case rowsAffected == 0:
		return &failure.ResourceNotFoundError{
			Message: "ingredient [" + recipeIngredient.ID + "] for recipe [" + recipeID + "] not found",
		}
	}
	return nil
}
