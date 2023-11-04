package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/remieven/miam/model"
	"github.com/remieven/miam/pb-lite/rest"
	"github.com/remieven/miam/service"
)

// IngredientHandler is an ingredient handler
type IngredientHandler struct {
	ingredientService *service.IngredientService
}

func newIngredientHandler(ingredientService *service.IngredientService) *IngredientHandler {
	return &IngredientHandler{
		ingredientService,
	}
}

// GetIngredients returns all known ingredients
func (handler *IngredientHandler) GetIngredients(responseWriter http.ResponseWriter, request *http.Request) {
	ingredients, err := handler.ingredientService.GetAllIngredients(request.Context())
	if rest.HandleErrorCase(responseWriter, err) {
		return
	}
	rest.WriteOKResponse(responseWriter, ingredients)
}

// UpdateIngredient updates an ingredient
func (handler *IngredientHandler) UpdateIngredient(responseWriter http.ResponseWriter, request *http.Request) {
	var baseIngredient model.BaseIngredient
	if err := json.NewDecoder(request.Body).Decode(&baseIngredient); rest.HandleParseBodyErrorCase(responseWriter, err) {
		return
	}

	vars := mux.Vars(request)
	ingredient, err := handler.ingredientService.UpdateIngredient(request.Context(), vars["id"], baseIngredient)
	if rest.HandleErrorCase(responseWriter, err) {
		return
	}
	rest.WriteOKResponse(responseWriter, ingredient)
}

// DeleteIngredient deletes the ingredient with the given id
func (handler *IngredientHandler) DeleteIngredient(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	if err := handler.ingredientService.DeleteIngredient(request.Context(), vars["id"]); rest.HandleErrorCase(responseWriter, err) {
		return
	}
	rest.WriteNoContentResponse(responseWriter)
}
