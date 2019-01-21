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
		create table if not exists recipe (name text);
	`
	if _, err := holder.db.Exec(initStatement); err != nil {
		return nil, err
	}
	return &RecipeDao{holder, recipeIngredientDao}, nil
}

// GetRecipe returns the recipe with the given ID or nil
func (dao *RecipeDao) GetRecipe(ID int) (*model.Recipe, error) {
	rows, err := dao.holder.db.Query("select oid, name from recipe where oid=?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var oid int
		var name string
		if err = rows.Scan(&oid, &name); err != nil {
			return nil, err
		}
		strID := strconv.Itoa(oid)
		return &model.Recipe{
			ID: strID,
			BaseRecipe: model.BaseRecipe{
				Name: name,
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
	insertStatement, err := transaction.Prepare("insert into recipe(name) values (?)")
	if err != nil {
		return "", err
	}
	defer insertStatement.Close()

	result, err := insertStatement.Exec(recipe.Name)
	if err != nil {
		return "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	transaction.Commit()
	return strconv.Itoa(int(id)), nil
}
