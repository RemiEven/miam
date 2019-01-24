package datasource

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/RemiEven/miam/model"
)

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

func (dao *RecipeIngredientDao) GetRecipeIngredients(recipeID int) ([]model.RecipeIngredient, error) {
	rows, err := dao.holder.db.Query(`select
		recipe_ingredient.ingredient_id, recipe_ingredient.quantity, ingredient.name
		from recipe_ingredient
		inner join ingredient
		on recipe_ingredient.ingredient_id=ingredient.oid
		where recipe_ingredient.recipe_id=?`, recipeID)
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
				ID: strconv.Itoa(ingredientID),
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

func (dao *RecipeIngredientDao) AddRecipeIngredient(transaction *sql.Tx, recipeID string, recipeIngredient model.RecipeIngredient) (string, error) {
	if len(recipeIngredient.ID) == 0 {
		// New ingredient
		ingredientID, err := dao.ingredientDao.AddIngredient(transaction, recipeIngredient.Name)
		if err != nil {
			return "", err
		}
		recipeIngredient.ID = ingredientID
	}
	intRecipeID, err := strconv.Atoi(recipeID)
	if err != nil {
		return "", err
	}
	intIngredientID, err := strconv.Atoi(recipeIngredient.ID)
	if err != nil {
		return "", err
	}

	insertStatement, err := transaction.Prepare("insert into recipe_ingredient(recipe_id, ingredient_id, quantity) values(?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer insertStatement.Close()

	_, err = insertStatement.Exec(intRecipeID, intIngredientID, recipeIngredient.Quantity)
	if err != nil {
		return "", err
	}
	return recipeIngredient.ID, nil
}

func (dao *RecipeIngredientDao) IsUsedInRecipe(ingredientID int) (bool, error) {
	rows, err := dao.holder.db.Query("select exists(select 1 from recipe_ingredient where ingredient_id=?)", ingredientID)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	if rows.Next() {
		var exists bool
		if err := rows.Scan(&exists); err != nil {
			return false, err
		}
		return exists, nil
	} else if err := rows.Err(); err != nil {
		return false, err
	}
	return false, fmt.Errorf("Fail to query for recipes using ingredient [%q]", ingredientID)
}

func (dao *RecipeIngredientDao) DeleteRecipeIngredients(transaction *sql.Tx, recipeID int) error {
	deleteStatement, err := transaction.Prepare("delete from recipe_ingredient where recipe_id=?")
	if err != nil {
		return err
	}
	defer deleteStatement.Close()

	_, err = deleteStatement.Exec(recipeID)
	return err
}
