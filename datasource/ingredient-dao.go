package datasource

import (
	"fmt"
	"strconv"

	"github.com/RemiEven/miam/model"
)

type IngredientDao struct {
	holder *DatabaseHolder
}

func NewIngredientDao(holder *DatabaseHolder) (*IngredientDao, error) {
	initStatement := `
		create table if not exists ingredient (name text)
	`
	if _, err := holder.db.Exec(initStatement); err != nil {
		return nil, err
	}
	return &IngredientDao{holder}, nil
}

// GetIngredient returns the ingredient with the given ID or nil
func (dao *IngredientDao) GetIngredient(ID int) (*model.Ingredient, error) {
	rows, err := dao.holder.db.Query("select oid, name from ingredient where oid=?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var oid int
		var name string
		if err := rows.Scan(&oid, &name); err != nil {
			return nil, err
		}
		return &model.Ingredient{
			ID: strconv.Itoa(oid),
			BaseIngredient: model.BaseIngredient{
				Name: name,
			},
		}, nil
	} else if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("No ingredient found with id [%d]", ID)
}
