package datasource

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RemiEven/miam/common"
	"github.com/RemiEven/miam/model"
)

// IngredientDao struct
type IngredientDao struct {
	holder *DatabaseHolder
}

func NewIngredientDao(holder *DatabaseHolder) (*IngredientDao, error) {
	initStatement := `
		create table if not exists ingredient (id integer primary key asc, name text)
	`
	if _, err := holder.DB.Exec(initStatement); err != nil {
		return nil, err
	}
	return &IngredientDao{holder}, nil
}

// GetIngredient returns the ingredient with the given ID or nil
func (dao *IngredientDao) GetIngredient(ctx context.Context, ID string) (*model.Ingredient, error) {
	oid, err := toSqliteID(ID)
	if err != nil {
		return nil, err
	}
	row := dao.holder.DB.QueryRowContext(ctx, "select name from ingredient where id=?", oid)
	var name string

	if err := row.Scan(&name); errors.Is(err, sql.ErrNoRows) {
		return nil, common.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &model.Ingredient{
		ID: ID,
		BaseIngredient: model.BaseIngredient{
			Name: name,
		},
	}, nil
}

// AddIngredient adds a new ingredient
func (dao *IngredientDao) AddIngredient(ctx context.Context, transaction *sql.Tx, name string) (string, error) {
	insertStatement, err := transaction.PrepareContext(ctx, "insert into ingredient(name) values(?)")
	if err != nil {
		return "", err
	}
	defer insertStatement.Close()

	result, err := insertStatement.ExecContext(ctx, name)
	if err != nil {
		return "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return fromSqliteID(sqliteID(id)), nil
}

// DeleteIngredient delete the ingredient with the given id if present.
// It is up to the caller to ensure no recipe uses the ingredient.
func (dao *IngredientDao) DeleteIngredient(ctx context.Context, ID string) error {
	oid, err := toSqliteID(ID)
	if err != nil {
		return err
	}
	deleteStatement, err := dao.holder.DB.PrepareContext(ctx, "delete from ingredient where id=?")
	if err != nil {
		return err
	}
	defer deleteStatement.Close()

	_, err = deleteStatement.ExecContext(ctx, oid)
	return err
}

// UpdateIngredient updates the name of an ingredient
func (dao *IngredientDao) UpdateIngredient(ctx context.Context, ingredient model.Ingredient) error {
	updateStatement, err := dao.holder.DB.PrepareContext(ctx, "update ingredient set name=?2 where id=?1")
	if err != nil {
		return err
	}
	defer updateStatement.Close()

	result, err := updateStatement.ExecContext(ctx, ingredient.ID, ingredient.Name)
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

// GetAllIngredients returns all ingredients
func (dao *IngredientDao) GetAllIngredients(ctx context.Context) ([]model.Ingredient, error) {
	rows, err := dao.holder.DB.QueryContext(ctx, "select id, name from ingredient")
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
			ID: fromSqliteID(id),
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
