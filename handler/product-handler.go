package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RemiEven/miam/dao"
	"github.com/gorilla/mux"
)

// ProductHandler handles a product request
func ProductHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	product := dao.GetProduct(vars["id"])
	json.NewEncoder(responseWriter).Encode(product)
}
