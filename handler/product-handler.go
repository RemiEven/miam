package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/RemiEven/miam/model"

	"github.com/RemiEven/miam/dao"
	"github.com/gorilla/mux"
)

// ProductHandler is a product handler
type ProductHandler struct {
	productDao *dao.ProductDao
}

// NewProductHandler creates a new product handler
func NewProductHandler(dao *dao.ProductDao) *ProductHandler {
	return &ProductHandler{
		dao,
	}
}

// GetProductByID handles a product request
func (handler *ProductHandler) GetProductByID(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := handler.productDao.GetProduct(id)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(responseWriter).Encode(product)
	}
}

// AddProduct adds a product
func (handler *ProductHandler) AddProduct(responseWriter http.ResponseWriter, request *http.Request) {
	var product model.EditableProduct
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := handler.productDao.AddProduct(&product)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		responseWriter.Header().Add("Location", id)
		responseWriter.WriteHeader(http.StatusCreated)
	}
}
