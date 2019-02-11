package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/model"
	"github.com/gorilla/mux"
)

// IngredientHandler is an ingredient handler
type IngredientHandler struct {
	ingredientDao *datasource.IngredientDao
}

// NewIngredientHandler creates a new ingredient handler
func NewIngredientHandler(dao *datasource.IngredientDao) *IngredientHandler {
	return &IngredientHandler{
		dao,
	}
}

// GetIngredients returns all known ingredients
func (handler *IngredientHandler) GetIngredients(responseWriter http.ResponseWriter, request *http.Request) {
	ingredients, err := handler.ingredientDao.GetAllIngredients()
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(responseWriter).Encode(ingredients)
	}
}

// UpdateIngredient updates an ingredient
func (handler *IngredientHandler) UpdateIngredient(responseWriter http.ResponseWriter, request *http.Request) {
	var baseIngredient model.BaseIngredient
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&baseIngredient)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(request)
	ingredient := model.Ingredient{
		ID:             vars["id"],
		BaseIngredient: baseIngredient,
	}
	if err = handler.ingredientDao.UpdateIngredient(ingredient); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(responseWriter).Encode(ingredient)
	}
}

// DeleteIngredient deletes the ingredient with the given id
func (handler *IngredientHandler) DeleteIngredient(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = handler.ingredientDao.DeleteIngredient(id); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		responseWriter.WriteHeader(http.StatusNoContent)
	}
}
