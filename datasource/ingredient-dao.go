package datasource

import (
	"database/sql"
	"strconv"

	"github.com/RemiEven/miam/common"
	"github.com/RemiEven/miam/model"
)

type IngredientDao struct {
	holder *databaseHolder
}

func newIngredientDao(holder *databaseHolder) (*IngredientDao, error) {
	initStatement := `
		create table if not exists ingredient (id integer primary key asc, name text)
	`
	if _, err := holder.db.Exec(initStatement); err != nil {
		return nil, err
	}
	return &IngredientDao{holder}, nil
}

// GetIngredient returns the ingredient with the given ID or nil
func (dao *IngredientDao) GetIngredient(ID int) (*model.Ingredient, error) {
	rows, err := dao.holder.db.Query("select id, name from ingredient where id=?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		return &model.Ingredient{
			ID: strconv.Itoa(id),
			BaseIngredient: model.BaseIngredient{
				Name: name,
			},
		}, nil
	} else if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, common.ErrNotFound
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

// DeleteIngredient delete the ingredient with the given id if present.
// It is up to the caller to ensure no recipe uses the ingredient.
func (dao *IngredientDao) DeleteIngredient(ID int) error {
	deleteStatement, err := dao.holder.db.Prepare("delete from ingredient where id=?")
	if err != nil {
		return err
	}
	defer deleteStatement.Close()

	_, err = deleteStatement.Exec(ID)
	return err
}

func (dao *IngredientDao) UpdateIngredient(ingredient model.Ingredient) error {
	updateStatement, err := dao.holder.db.Prepare("update ingredient set name=?2 where id=?1")
	if err != nil {
		return err
	}
	defer updateStatement.Close()

	result, err := updateStatement.Exec(ingredient.ID, ingredient.Name)
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

func (dao *IngredientDao) GetAllIngredients() ([]model.Ingredient, error) {
	rows, err := dao.holder.db.Query("select id, name from ingredient")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ingredients := make([]model.Ingredient, 0, 50) // 50 is arbitrary

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, model.Ingredient{
			ID: strconv.Itoa(id),
			BaseIngredient: model.BaseIngredient{
				Name: name,
			},
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ingredients, nil
}
