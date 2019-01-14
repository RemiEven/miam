package dao

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/RemiEven/miam/model"
	// Necessary to load sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// ProductDao is a product dao
type ProductDao struct {
	db *sql.DB
}

// NewProductDao return a new product dao
func NewProductDao() (*ProductDao, error) {
	db, err := sql.Open("sqlite3", "./product.db")
	if err != nil {
		return nil, err
	}
	initStatement := `
		create table if not exists product (name text);
	`
	_, err = db.Exec(initStatement)
	if err != nil {
		return nil, err
	}

	return &ProductDao{db}, nil
}

// GetProduct returns the product with the given ID or nil
func (dao *ProductDao) GetProduct(id int) (*model.Product, error) {
	rows, err := dao.db.Query("select oid, name from product where oid=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var oid int
		var name string
		err = rows.Scan(&oid, &name)
		if err != nil {
			return nil, err
		}
		strID := strconv.Itoa(oid)
		return &model.Product{
			ID:   strID,
			Name: name,
		}, nil
	} else if err := rows.Err(); err != nil {
		return nil, err
	} else {
		return nil, fmt.Errorf("No product found with id [%d]", id)
	}
}

// AddProduct adds the given product
func (dao *ProductDao) AddProduct(product *model.EditableProduct) (string, error) {
	transaction, err := dao.db.Begin()
	if err != nil {
		return "", err
	}
	insertStatement, err := transaction.Prepare("insert into product(name) values (?)")
	if err != nil {
		return "", err
	}
	defer insertStatement.Close()

	result, err := insertStatement.Exec(product.Name)
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

// Close closes the dao connection to the underlying database
func (dao *ProductDao) Close() error {
	log.Println("Closing product database connection")
	return dao.db.Close()
}
