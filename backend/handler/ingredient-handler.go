package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RemiEven/miam/model"
	"github.com/RemiEven/miam/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// IngredientHandler is an ingredient handler
type IngredientHandler struct {
	ingredientService *service.IngredientService
}

// NewIngredientHandler creates a new ingredient handler
func NewIngredientHandler(ingredientService *service.IngredientService) *IngredientHandler {
	return &IngredientHandler{
		ingredientService,
	}
}

// GetIngredients returns all known ingredients
func (handler *IngredientHandler) GetIngredients(responseWriter http.ResponseWriter, request *http.Request) {
	ingredients, err := handler.ingredientService.GetAllIngredients()
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(ingredients)
}

// UpdateIngredient updates an ingredient
func (handler *IngredientHandler) UpdateIngredient(responseWriter http.ResponseWriter, request *http.Request) {
	var baseIngredient model.BaseIngredient
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&baseIngredient)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(request)
	ingredient, err := handler.ingredientService.UpdateIngredient(vars["id"], baseIngredient)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(ingredient)
}

// DeleteIngredient deletes the ingredient with the given id
func (handler *IngredientHandler) DeleteIngredient(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	if err := handler.ingredientService.DeleteIngredient(vars["id"]); err != nil {
		logrus.Error(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusNoContent)
}
