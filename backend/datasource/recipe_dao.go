package datasource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/remieven/miam/model"
	"github.com/remieven/miam/pb-lite/failure"
)

// RecipeDao is a recipe dao
type RecipeDao struct {
	holder              *DatabaseHolder
	recipeIngredientDao *RecipeIngredientDao
}

// NewRecipeDao returns a new recipe dao
func NewRecipeDao(holder *DatabaseHolder, recipeIngredientDao *RecipeIngredientDao) (*RecipeDao, error) {
	initStatement := `
		create table if not exists recipe (id integer primary key asc, name text, how_to text);
	`
	if _, err := holder.DB.Exec(initStatement); err != nil {
		return nil, fmt.Errorf("failed to create recipe table: %w", err)
	}
	return &RecipeDao{holder, recipeIngredientDao}, nil
}

// GetRecipe returns the recipe with the given ID or nil
func (dao *RecipeDao) GetRecipe(ctx context.Context, ID string) (*model.Recipe, error) {
	oid, err := toSqliteID(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert [%s] to sqlite ID: %w", ID, err)
	}
	row := dao.holder.DB.QueryRowContext(ctx, "select name, how_to from recipe where id=?", oid)
	var name, howTo string

	if err := row.Scan(&name, &howTo); errors.Is(err, sql.ErrNoRows) {
		return nil, &failure.ResourceNotFoundError{
			Message: "recipe [" + ID + "] not found",
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve recipe: %w", err)
	}
	ingredients, err := dao.recipeIngredientDao.GetRecipeIngredients(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recipe ingredients: %w", err)
	}

	return &model.Recipe{
		ID: ID,
		BaseRecipe: model.BaseRecipe{
			Name:        name,
			HowTo:       howTo,
			Ingredients: ingredients,
		},
	}, nil
}

// GetRecipes returns the recipes with the given IDs or an empty slice
func (dao *RecipeDao) GetRecipes(ctx context.Context, IDs []string) ([]model.Recipe, error) {
	if len(IDs) == 0 {
		return []model.Recipe{}, nil
	}
	queryParamPlaceholders := "?" + strings.Repeat(",?", len(IDs)-1)
	queryParams := make([]interface{}, len(IDs))
	var err error
	for i := range IDs {
		if queryParams[i], err = toSqliteID(IDs[i]); err != nil {
			return nil, fmt.Errorf("failed to convert [%s] to sqlite ID: %w", IDs[i], err)
		}
	}
	rows, err := dao.holder.DB.QueryContext(ctx, "select id, name, how_to from recipe where id in ("+queryParamPlaceholders+")", queryParams...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recipes: %w", err)
	}
	defer rows.Close()
	results := make([]model.Recipe, 0, len(IDs))
	for rows.Next() {
		var id sqliteID
		var name, howTo string
		if err = rows.Scan(&id, &name, &howTo); err != nil {
			return nil, fmt.Errorf("failed to scan recipe row: %w", err)
		}
		recipeID := fromSqliteID(id)
		ingredients, err := dao.recipeIngredientDao.GetRecipeIngredients(ctx, recipeID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve recipe ingredients: %w", err)
		}

		results = append(results, model.Recipe{
			ID: recipeID,
			BaseRecipe: model.BaseRecipe{
				Name:        name,
				HowTo:       howTo,
				Ingredients: ingredients,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("got an error while iterating on recipe rows: %w", err)
	}

	return results, nil
}

// AddRecipe adds the given recipe
func (dao *RecipeDao) AddRecipe(ctx context.Context, recipe *model.BaseRecipe) (string, error) {
	transaction, err := dao.holder.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to init transaction: %w", err)
	}
	insertStatement, err := transaction.PrepareContext(ctx, "insert into recipe(name, how_to) values (?, ?)")
	if err != nil {
		return "", fmt.Errorf("failed to prepare recipe statement: %w", err)
	}
	defer insertStatement.Close()

	result, err := insertStatement.ExecContext(ctx, recipe.Name, recipe.HowTo)
	if err != nil {
		rollback(transaction)
		return "", fmt.Errorf("failed to execute insert recipe statement: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		rollback(transaction)
		return "", fmt.Errorf("failed to retrieve ID of inserted recipe: %w", err)
	}
	recipeID := fromSqliteID(sqliteID(id))

	for _, recipeIngredient := range recipe.Ingredients {
		if _, err := dao.recipeIngredientDao.AddRecipeIngredient(ctx, transaction, recipeID, recipeIngredient); err != nil {
			rollback(transaction)
			return "", fmt.Errorf("failed to add ingredient: %w", err)
		}
	}

	if err := transaction.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return recipeID, nil
}

// DeleteRecipe deletes a recipe and its ingredients
func (dao *RecipeDao) DeleteRecipe(ctx context.Context, ID string) error {
	oid, err := toSqliteID(ID)
	if err != nil {
		return fmt.Errorf("failed to convert [%s] to sqlite ID: %w", ID, err)
	}
	transaction, err := dao.holder.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to init transaction: %w", err)
	}
	deleteStatement, err := transaction.PrepareContext(ctx, "delete from recipe where id=?")
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer deleteStatement.Close()

	if _, err := deleteStatement.ExecContext(ctx, oid); err != nil {
		rollback(transaction)
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}
	if err = dao.recipeIngredientDao.DeleteRecipeIngredients(ctx, transaction, ID); err != nil {
		rollback(transaction)
		return fmt.Errorf("failed to delete recipe ingredients: %w", err)
	}

	return transaction.Commit()
}

// UpdateRecipe updates a recipe
func (dao *RecipeDao) UpdateRecipe(ctx context.Context, recipe model.Recipe) (*model.Recipe, error) {
	transaction, err := dao.holder.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to init transaction: %w", err)
	}
	updateStatement, err := transaction.PrepareContext(ctx, "update recipe set (name, how_to) = (?2, ?3) where id=?1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}

	result, err := updateStatement.ExecContext(ctx, recipe.ID, recipe.Name, recipe.HowTo)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	switch {
	case err != nil:
		rollback(transaction)
		return nil, fmt.Errorf("failed to retrieve number of rows affected by update statement: %w", err)
	case rowsAffected == 0:
		rollback(transaction)
		return nil, &failure.ResourceNotFoundError{
			Message: "recipe [" + recipe.ID + "] not found",
		}
	}

	currentIngredients, err := dao.recipeIngredientDao.GetRecipeIngredients(ctx, recipe.ID)
	if err != nil {
		rollback(transaction)
		return nil, fmt.Errorf("failed to retrieve recipe ingredients: %w", err)
	}
	for _, currentIngredient := range currentIngredients {
		stillThere, newIngredient := containsIngredient(currentIngredient, recipe.Ingredients)
		if !stillThere {
			if err = dao.recipeIngredientDao.DeleteRecipeIngredient(ctx, transaction, recipe.ID, currentIngredient.ID); err != nil {
				rollback(transaction)
				return nil, fmt.Errorf("failed to remove ingredient from recipe: %w", err)
			}
		} else if newIngredient.Quantity != currentIngredient.Quantity {
			if err = dao.recipeIngredientDao.UpdateRecipeIngredient(ctx, transaction, recipe.ID, newIngredient); err != nil {
				rollback(transaction)
				return nil, fmt.Errorf("failed to update quantity of recipe ingredient: %w", err)
			}

		}
	}
	for i, newIngredient := range recipe.Ingredients {
		if alreadyThere, _ := containsIngredient(newIngredient, currentIngredients); !alreadyThere {
			ingredientID, err := dao.recipeIngredientDao.AddRecipeIngredient(ctx, transaction, recipe.ID, newIngredient)
			if err != nil {
				rollback(transaction)
				return nil, fmt.Errorf("failed to add recipe ingredient: %w", err)
			}
			recipe.Ingredients[i].ID = ingredientID
		}
	}

	if err := transaction.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &recipe, nil
}

// containsIngredient returns whether a given recipe ingredient is present in a slice of recipe ingredients
func containsIngredient(searched model.RecipeIngredient, ingredients []model.RecipeIngredient) (bool, model.RecipeIngredient) {
	for _, ingredient := range ingredients {
		if ingredient.ID == searched.ID {
			return true, ingredient
		}
	}
	return false, model.RecipeIngredient{}
}

// GetRandomRecipes search for recipes according to given search criteria
func (dao *RecipeDao) GetRandomRecipes(ctx context.Context, search model.RecipeSearch) (*model.RecipeSearchResult, error) {
	results, err := dao.getRandomRecipes(ctx, 10) // 10 is arbitrary
	if err != nil {
		return nil, fmt.Errorf("failed to get random recipes: %w", err)
	}
	total, err := dao.getRecipeCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count recipes: %w", err)
	}

	return &model.RecipeSearchResult{
		FirstResults: results,
		Total:        total,
	}, nil
}

// getRandomRecipes returns a given number of randomly selected recipes
func (dao *RecipeDao) getRandomRecipes(ctx context.Context, numberWanted int) ([]model.Recipe, error) {
	rows, err := dao.holder.DB.QueryContext(ctx, "select id, name, how_to from recipe where id in (select id from recipe order by random() limit ?)", numberWanted)
	if err != nil {
		return nil, fmt.Errorf("failed to query random recipes: %w", err)
	}
	defer rows.Close()
	results := make([]model.Recipe, 0, numberWanted)
	for rows.Next() {
		var id int
		var name, howTo string
		if err = rows.Scan(&id, &name, &howTo); err != nil {
			return nil, fmt.Errorf("failed to scan recipe row: %w", err)
		}
		recipeID := fromSqliteID(id)
		ingredients, err := dao.recipeIngredientDao.GetRecipeIngredients(ctx, recipeID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve recipe ingredients: %w", err)
		}

		results = append(results, model.Recipe{
			ID: recipeID,
			BaseRecipe: model.BaseRecipe{
				Name:        name,
				HowTo:       howTo,
				Ingredients: ingredients,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("got an error while iterating on recipe rows: %w", err)
	}
	return results, nil
}

// getRecipeCount returns the number of saved recipes
func (dao *RecipeDao) getRecipeCount(ctx context.Context) (int, error) {
	rows, err := dao.holder.DB.QueryContext(ctx, "select count(*) from recipe")
	if err != nil {
		return 0, fmt.Errorf("query to count recipes failed: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return 0, fmt.Errorf("failed to scan count of recipes: %w", err)
		}
		return count, nil
	} else if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("got an error before scanning count of recipes: %w", err)
	}
	return 0, errors.New("no row after select count SQL request")
}

func rollback(transaction *sql.Tx) {
	if err := transaction.Rollback(); err != nil {
		slog.With("error", err).Error("transaction rollback failed")
	}
}

// ListRecipeIds returns all recipe ids
func (dao *RecipeDao) ListRecipeIds(ctx context.Context) ([]string, error) {
	rows, err := dao.holder.DB.QueryContext(ctx, "select id from recipe")
	if err != nil {
		return nil, fmt.Errorf("failed to query recipe ids: %w", err)
	}

	defer rows.Close()
	ids := make([]string, 0)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan recipe id row: %w", err)
		}
		ids = append(ids, fromSqliteID(id))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("got an error while iterating on recipe id rows: %w", err)
	}

	return ids, nil
}
