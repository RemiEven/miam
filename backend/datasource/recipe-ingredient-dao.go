package datasource

import (
	"context"
	"database/sql"

	"github.com/RemiEven/miam/common"

	"github.com/RemiEven/miam/model"
)

// RecipeIngredientDao struct
type RecipeIngredientDao struct {
	holder        *databaseHolder
	ingredientDao *IngredientDao
}

// NewRecipeIngredientDao returns a new recipe ingredient dao
func newRecipeIngredientDao(holder *databaseHolder, ingredientDao *IngredientDao) (*RecipeIngredientDao, error) {
	initStatement := `
		create table if not exists recipe_ingredient (recipe_id int, ingredient_id int, quantity text);
		create index if not exists recipe_id_index on recipe_ingredient(recipe_id);
		create index if not exists ingredient_id_index on recipe_ingredient(ingredient_id);
	`
	if _, err := holder.db.Exec(initStatement); err != nil {
		return nil, err
	}
	return &RecipeIngredientDao{holder, ingredientDao}, nil
}

// GetRecipeIngredients returns the ingredient of a recipe
func (dao *RecipeIngredientDao) GetRecipeIngredients(ctx context.Context, recipeID string) ([]model.RecipeIngredient, error) {
	intRecipeID, err := toSqliteID(recipeID)
	if err != nil {
		return nil, err
	}
	rows, err := dao.holder.db.QueryContext(ctx, `select
		recipe_ingredient.ingredient_id, recipe_ingredient.quantity, ingredient.name
		from recipe_ingredient
		inner join ingredient
		on recipe_ingredient.ingredient_id=ingredient.id
		where recipe_ingredient.recipe_id=?`, intRecipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	recipeIngredients := make([]model.RecipeIngredient, 0, 8) // most recipes have 8 or less ingredients
	for rows.Next() {
		var ingredientID int
		var quantity string
		var name string
		if err := rows.Scan(&ingredientID, &quantity, &name); err != nil {
			return nil, err
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
		return nil, err
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
			return "", err
		}
		recipeIngredient.ID = ingredientID
	}
	intRecipeID, err := toSqliteID(recipeID)
	if err != nil {
		return "", err
	}
	intIngredientID, err := toSqliteID(recipeIngredient.ID)
	if err != nil {
		return "", err
	}

	insertStatement, err := transaction.PrepareContext(ctx, "insert into recipe_ingredient(recipe_id, ingredient_id, quantity) values(?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer insertStatement.Close()

	_, err = insertStatement.ExecContext(ctx, intRecipeID, intIngredientID, recipeIngredient.Quantity)
	if err != nil {
		return "", err
	}
	return recipeIngredient.ID, nil
}

// IsUsedInRecipe returns whether an ingredient is used by at least one recipe
func (dao *RecipeIngredientDao) IsUsedInRecipe(ctx context.Context, ingredientID string) (bool, error) {
	intID, err := toSqliteID(ingredientID)
	if err != nil {
		return false, err
	}
	row := dao.holder.db.QueryRowContext(ctx, "select exists(select 1 from recipe_ingredient where ingredient_id=?)", intID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// DeleteRecipeIngredients deletes the ingredients of a recipe
func (dao *RecipeIngredientDao) DeleteRecipeIngredients(ctx context.Context, transaction *sql.Tx, recipeID string) error {
	intID, err := toSqliteID(recipeID)
	if err != nil {
		return err
	}
	deleteStatement, err := transaction.PrepareContext(ctx, "delete from recipe_ingredient where recipe_id=?")
	if err != nil {
		return err
	}
	defer deleteStatement.Close()

	_, err = deleteStatement.ExecContext(ctx, intID)
	return err
}

// DeleteRecipeIngredient removes one ingredient from a recipe
func (dao *RecipeIngredientDao) DeleteRecipeIngredient(ctx context.Context, transaction *sql.Tx, recipeID, ingredientID string) error {
	intRecipeID, err := toSqliteID(recipeID)
	if err != nil {
		return err
	}
	intIngredientID, err := toSqliteID(ingredientID)
	if err != nil {
		return err
	}
	deleteStatement, err := transaction.PrepareContext(ctx, "delete from recipe_ingredient where recipe_id=? and ingredient_id=?")
	if err != nil {
		return err
	}
	defer deleteStatement.Close()

	_, err = deleteStatement.Exec(ctx, intRecipeID, intIngredientID)
	return err
}

// UpdateRecipeIngredient updates an ingredient of a recipe
func (dao *RecipeIngredientDao) UpdateRecipeIngredient(ctx context.Context, transaction *sql.Tx, recipeID string, recipeIngredient model.RecipeIngredient) error {
	intID, err := toSqliteID(recipeID)
	if err != nil {
		return err
	}
	updateStatement, err := transaction.PrepareContext(ctx, "update recipe_ingredient set quantity=?3 where recipe_id=?1 and ingredient_id=?2")
	if err != nil {
		return err
	}
	defer updateStatement.Close()

	result, err := updateStatement.ExecContext(ctx, intID, recipeIngredient.ID, recipeIngredient.Quantity)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	switch {
	case err != nil:
		return err
	case rowsAffected == 0:
		return common.ErrNotFound
	}
	return nil
}
