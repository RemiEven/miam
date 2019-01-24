package datasource

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/RemiEven/miam/model"
)

type IngredientDao struct {
	holder *databaseHolder
}

func newIngredientDao(holder *databaseHolder) (*IngredientDao, error) {
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

func (dao *IngredientDao) AddIngredient(transaction *sql.Tx, name string) (string, error) {
	insertStatement, err := transaction.Prepare("insert into ingredient(name) values(?)")
	if err != nil {
		return "", err
	}
	defer insertStatement.Close()

	result, err := insertStatement.Exec(name)
	if err != nil {
		return "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(id)), nil
}

// DeleteIngredient deletes the ingredient with the given id if present.
// It is up to the caller to ensure no recipe uses the ingredient.
func (dao *IngredientDao) DeleteIngredient(ID int) error {
	deleteStatement, err := dao.holder.db.Prepare("delete from ingredient where oid=?")
	if err != nil {
		return err
	}
	defer deleteStatement.Close()

	_, err = deleteStatement.Exec(ID)
	return err
}
