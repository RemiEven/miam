package dao

import (
	"github.com/RemiEven/miam/model"
)

var products = map[string]model.Product{
	"1994": {
		ID:   "1994",
		Name: "jojo",
	},
	"2007": {
		ID:   "2007",
		Name: "jij8",
	},
}

// GetProduct returns the product with the given ID or nil
func GetProduct(ID string) model.Product {
	return products[ID]
}
