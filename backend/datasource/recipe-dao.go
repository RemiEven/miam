package datasource

import (
	"errors"
	"strconv"

	"github.com/RemiEven/miam/common"
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
		create table if not exists recipe (id integer primary key asc, name text, how_to text);
	`
	if _, err := holder.db.Exec(initStatement); err != nil {
		return nil, err
	}
	return &RecipeDao{holder, recipeIngredientDao}, nil
}

// GetRecipe returns the recipe with the given ID or nil
func (dao *RecipeDao) GetRecipe(ID int) (*model.Recipe, error) {
	rows, err := dao.holder.db.Query("select id, name, how_to from recipe where id=?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var id int
		var name, howTo string
		if err = rows.Scan(&id, &name, &howTo); err != nil {
			return nil, err
		}
		ingredients, err := dao.recipeIngredientDao.GetRecipeIngredients(ID)
		if err != nil {
			return nil, err
		}

		return &model.Recipe{
			ID: strconv.Itoa(id),
			BaseRecipe: model.BaseRecipe{
				Name:        name,
				HowTo:       howTo,
				Ingredients: ingredients,
			},
		}, nil
	} else if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, common.ErrNotFound
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

// DeleteRecipe deletes a recipe and its ingredients
func (dao *RecipeDao) DeleteRecipe(ID int) error {
	transaction, err := dao.holder.db.Begin()
	if err != nil {
		return err
	}
	deleteStatement, err := transaction.Prepare("delete from recipe where id=?")
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

// UpdateRecipe updates a recipe
func (dao *RecipeDao) UpdateRecipe(recipe model.Recipe) (*model.Recipe, error) {
	intRecipeID, err := strconv.Atoi(recipe.ID)
	if err != nil {
		return nil, err
	}
	transaction, err := dao.holder.db.Begin()
	if err != nil {
		return nil, err
	}
	updateStatement, err := transaction.Prepare("update recipe set (name, how_to) = (?2, ?3) where id=?1")
	if err != nil {
		return nil, err
	}

	result, err := updateStatement.Exec(intRecipeID, recipe.Name, recipe.HowTo)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	switch {
	case err != nil:
		transaction.Rollback() // TODO: log if fail to rollback
		return nil, err
	case rowsAffected == 0:
		transaction.Rollback() // TODO: log if fail to rollback
		return nil, common.ErrNotFound
	}

	currentIngredients, err := dao.recipeIngredientDao.GetRecipeIngredients(intRecipeID)
	if err != nil {
		transaction.Rollback() // TODO: log if fail to rollback
		return nil, err
	}
	for _, currentIngredient := range currentIngredients {
		stillThere, newIngredient := containsIngredient(currentIngredient, recipe.Ingredients)
		if !stillThere {
			intIngredientID, err := strconv.Atoi(currentIngredient.ID)
			if err != nil {
				return nil, err
			}
			if err = dao.recipeIngredientDao.DeleteRecipeIngredient(transaction, intRecipeID, intIngredientID); err != nil {
				transaction.Rollback() // TODO: log if fail to rollback
				return nil, err
			}
		} else if newIngredient.Quantity != currentIngredient.Quantity {
			if err = dao.recipeIngredientDao.UpdateRecipeIngredient(transaction, intRecipeID, newIngredient); err != nil {
				transaction.Rollback() // TODO: log if fail to rollback
				return nil, err
			}

		}
	}
	for i, newIngredient := range recipe.Ingredients {
		if alreadyThere, _ := containsIngredient(newIngredient, currentIngredients); !alreadyThere {
			ingredientID, err := dao.recipeIngredientDao.AddRecipeIngredient(transaction, recipe.ID, newIngredient)
			if err != nil {
				transaction.Rollback() // TODO: log if fail to rollback
				return nil, err
			}
			recipe.Ingredients[i].ID = ingredientID
		}
	}

	transaction.Commit()
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

// SearchRecipes search for recipes according to given search criteria
func (dao *RecipeDao) SearchRecipes(search model.RecipeSearch) (*model.RecipeSearchResult, error) {
	results, err := dao.getRandomRecipes(5) // 5 is arbitrary
	if err != nil {
		return nil, err
	}
	total, err := dao.getRecipeCount()
	if err != nil {
		return nil, err
	}

	return &model.RecipeSearchResult{
		FirstResults: results,
		Total:        total,
	}, nil
}

// getRandomRecipes returns a given number of randomly selected recipes
func (dao *RecipeDao) getRandomRecipes(numberWanted int) ([]model.Recipe, error) {
	rows, err := dao.holder.db.Query("select id, name, how_to from recipe where id in (select id from recipe order by random() limit ?)", numberWanted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]model.Recipe, 0, numberWanted)
	for rows.Next() {
		var id int
		var name, howTo string
		if err = rows.Scan(&id, &name, &howTo); err != nil {
			return nil, err
		}
		ingredients, err := dao.recipeIngredientDao.GetRecipeIngredients(id)
		if err != nil {
			return nil, err
		}

		results = append(results, model.Recipe{
			ID: strconv.Itoa(id),
			BaseRecipe: model.BaseRecipe{
				Name:        name,
				HowTo:       howTo,
				Ingredients: ingredients,
			},
		})
	}
	return results, rows.Err()
}

// getRecipeCount returns the number of saved recipes
func (dao *RecipeDao) getRecipeCount() (int, error) {
	rows, err := dao.holder.db.Query("select count(*) from recipe")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	} else if err := rows.Err(); err != nil {
		return 0, err
	}
	return 0, errors.New("No row after select count SQL request")
}
