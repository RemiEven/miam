package rest

import (
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/RemiEven/miam/model"
	"github.com/RemiEven/miam/service"
)

// RecipeHandler is a recipe handler
type RecipeHandler struct {
	recipeService *service.RecipeService
}

func newRecipeHandler(recipeService *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService,
	}
}

// GetRecipeByID handles a recipe request
func (handler *RecipeHandler) GetRecipeByID(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	recipe, err := handler.recipeService.GetRecipe(request.Context(), vars["id"])
	if err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(recipe)
}

// AddRecipe adds a recipe
func (handler *RecipeHandler) AddRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var recipe model.BaseRecipe
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&recipe)
	if err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := handler.recipeService.AddRecipe(request.Context(), recipe)
	if err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Add("Location", id)
	responseWriter.WriteHeader(http.StatusCreated)
}

// UpdateRecipe updates a recipe
func (handler *RecipeHandler) UpdateRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var recipe model.BaseRecipe
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&recipe)
	if err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	id := mux.Vars(request)["id"]
	updated, err := handler.recipeService.UpdateRecipe(request.Context(), id, recipe)
	if err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(updated)
}

// DeleteRecipe deletes a recipe
func (handler *RecipeHandler) DeleteRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	strID := vars["id"]

	if err := handler.recipeService.DeleteRecipe(request.Context(), strID); err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusNoContent)
}

// SearchRecipe searches for recipes
func (handler *RecipeHandler) SearchRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var search model.RecipeSearch
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&search)
	if err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	results, err := handler.recipeService.SearchRecipe(request.Context(), search)
	if err != nil {
		log.Error().Err(err).Msg("")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(results)
}
