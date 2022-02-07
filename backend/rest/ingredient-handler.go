package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/RemiEven/miam/model"
	"github.com/RemiEven/miam/service"
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
	if err != nil {
		log.Error().Err(err).Msg("")
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
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(request)
	ingredient, err := handler.ingredientService.UpdateIngredient(request.Context(), vars["id"], baseIngredient)
	if err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(ingredient)
}

// DeleteIngredient deletes the ingredient with the given id
func (handler *IngredientHandler) DeleteIngredient(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	if err := handler.ingredientService.DeleteIngredient(request.Context(), vars["id"]); err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusNoContent)
}
