package datasource

import (
	"fmt"
	"strconv"

	"github.com/RemiEven/miam/model"
)

// RecipeDao is a recipe dao
type RecipeDao struct {
	holder              *databaseHolder
	recipeIngredientDao *RecipeIngredientDao
}

// NewRecipeDao returns a new recipe dao
func newRecipeDao(holder *databaseHolder, recipeIngredientDao *RecipeIngredientDao) (*RecipeDao, error) {
	initStatement := `
		create table if not exists recipe (name text, how_to text);
	`
	if _, err := holder.db.Exec(initStatement); err != nil {
		return nil, err
	}
	return &RecipeDao{holder, recipeIngredientDao}, nil
}

// GetRecipe returns the recipe with the given ID or nil
func (dao *RecipeDao) GetRecipe(ID int) (*model.Recipe, error) {
	rows, err := dao.holder.db.Query("select oid, name, how_to from recipe where oid=?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var oid int
		var name, howTo string
		if err = rows.Scan(&oid, &name, &howTo); err != nil {
			return nil, err
		}
		strID := strconv.Itoa(oid)
		ingredients, err := dao.recipeIngredientDao.GetRecipeIngredients(ID)
		if err != nil {
			return nil, err
		}

		return &model.Recipe{
			ID: strID,
			BaseRecipe: model.BaseRecipe{
				Name:        name,
				HowTo:       howTo,
				Ingredients: ingredients,
			},
		}, nil
	} else if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("No recipe found with id [%d]", ID)
}

// AddRecipe adds the given recipe
func (dao *RecipeDao) AddRecipe(recipe *model.BaseRecipe) (string, error) {
	transaction, err := dao.holder.db.Begin()
	if err != nil {
		return "", err
	}
	insertStatement, err := transaction.Prepare("insert into recipe(name, how_to) values (?, ?)")
	if err != nil {
		return "", err
	}
	defer insertStatement.Close()

	result, err := insertStatement.Exec(recipe.Name, recipe.HowTo)
	if err != nil {
		transaction.Rollback() // TODO: log if fail to rollback
		return "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback() // TODO: log if fail to rollback
		return "", err
	}
	recipeID := strconv.Itoa(int(id))

	for _, recipeIngredient := range recipe.Ingredients {
		_, err = dao.recipeIngredientDao.AddRecipeIngredient(transaction, recipeID, recipeIngredient)
		if err != nil {
			transaction.Rollback() // TODO: log if fail to rollback
		}
	}

	return recipeID, transaction.Commit()
}

func (dao *RecipeDao) DeleteRecipe(ID int) error {
	transaction, err := dao.holder.db.Begin()
	if err != nil {
		return err
	}
	deleteStatement, err := transaction.Prepare("delete from recipe where oid=?")
	if err != nil {
		return err
	}
	defer deleteStatement.Close()

	if _, err := deleteStatement.Exec(ID); err != nil {
		transaction.Rollback() // TODO: log if fail to rollback
		return err
	}
	if err = dao.recipeIngredientDao.DeleteRecipeIngredients(transaction, ID); err != nil {
		transaction.Rollback() // TODO: log if fail to rollback
		return err
	}

	return transaction.Commit()
}
