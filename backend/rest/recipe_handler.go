package rest

import (
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/remieven/miam/model"
	"github.com/remieven/miam/pb-lite/rest"
	"github.com/remieven/miam/service"
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
	if rest.HandleErrorCase(responseWriter, err) {
		return
	}

	rest.WriteOKResponse(responseWriter, recipe)
}

// AddRecipe adds a recipe
func (handler *RecipeHandler) AddRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var recipe model.BaseRecipe
	if err := json.NewDecoder(request.Body).Decode(&recipe); rest.HandleParseBodyErrorCase(responseWriter, err) {
		return
	}

	id, err := handler.recipeService.AddRecipe(request.Context(), recipe)
	if rest.HandleErrorCase(responseWriter, err) {
		return
	}

	rest.WriteCreatedResponse(responseWriter, request, id)
}

// UpdateRecipe updates a recipe
func (handler *RecipeHandler) UpdateRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var recipe model.BaseRecipe
	if err := json.NewDecoder(request.Body).Decode(&recipe); rest.HandleParseBodyErrorCase(responseWriter, err) {
		return
	}

	id := mux.Vars(request)["id"]
	updated, err := handler.recipeService.UpdateRecipe(request.Context(), id, recipe)
	if rest.HandleErrorCase(responseWriter, err) {
		return
	}

	rest.WriteOKResponse(responseWriter, updated)
}

// DeleteRecipe deletes a recipe
func (handler *RecipeHandler) DeleteRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	strID := vars["id"]

	if err := handler.recipeService.DeleteRecipe(request.Context(), strID); rest.HandleErrorCase(responseWriter, err) {
		return
	}
	rest.WriteNoContentResponse(responseWriter)
}

// SearchRecipe searches for recipes
func (handler *RecipeHandler) SearchRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var search model.RecipeSearch
	if err := json.NewDecoder(request.Body).Decode(&search); rest.HandleParseBodyErrorCase(responseWriter, err) {
		return
	}

	results, err := handler.recipeService.SearchRecipe(request.Context(), search)
	if rest.HandleErrorCase(responseWriter, err) {
		return
	}

	rest.WriteOKResponse(responseWriter, results)
}
