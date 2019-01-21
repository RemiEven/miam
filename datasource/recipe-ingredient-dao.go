package datasource

import (
	"strconv"

	"github.com/RemiEven/miam/model"
)

type RecipeIngredientDao struct {
	holder *DatabaseHolder
}

// NewRecipeIngredientDao returns a new recipe ingredient dao
func NewRecipeIngredientDao(holder *DatabaseHolder) (*RecipeIngredientDao, error) {
	initStatement := `
		create table if not exists recipe_ingredient (recipe_id text, ingredient_id text, quantity text)
	`
	if _, err := holder.db.Exec(initStatement); err != nil {
		return nil, err
	}
	return &RecipeIngredientDao{holder}, nil
}

func (dao *RecipeIngredientDao) GetRecipeIngredients(recipeID int) ([]model.RecipeIngredient, error) {
	rows, err := dao.holder.db.Query("select recipe_id, ingredient_id, quantity from recipe_ingredient where recipe_id=?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	recipeIngredients := make([]model.RecipeIngredient, 8) // most recipes have 8 or less ingredients
	for rows.Next() {
		var recipeID int
		var ingredientID int
		var quantity string
		if err := rows.Scan(&recipeID, &ingredientID, &quantity); err != nil {
			return nil, err
		}
		recipeIngredients = append(recipeIngredients, model.RecipeIngredient{
			Ingredient: model.Ingredient{
				ID:   strconv.Itoa(ingredientID),
				BaseIngredient: model.BaseIngredient{
					Name: "name", // TODO
				}
			},
			Quantity: quantity,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return recipeIngredients, nil
}
