package datasource

import (
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
		create table if not exists recipe_ingredient (recipe_id text, ingredient_id text, quantity text)
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
	recipeIngredients := make([]model.RecipeIngredient, 8) // most recipes have 8 or less ingredients
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
