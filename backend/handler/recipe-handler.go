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

// FIXME: content-type application/json is not set in response where it should be

// RecipeHandler is a recipe handler
type RecipeHandler struct {
	recipeDao *datasource.RecipeDao
	searchDao *datasource.RecipeSearchDao
}

// NewRecipeHandler creates a new recipe handler
func NewRecipeHandler(recipeDao *datasource.RecipeDao, searchDao *datasource.RecipeSearchDao) *RecipeHandler {
	return &RecipeHandler{
		recipeDao,
		searchDao,
	}
}

// GetRecipeByID handles a recipe request
func (handler *RecipeHandler) GetRecipeByID(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	recipe, err := handler.recipeDao.GetRecipe(id)
	if err != nil { // TODO: handle errNotFound differently
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(responseWriter).Encode(recipe)
	}
}

// AddRecipe adds a recipe
func (handler *RecipeHandler) AddRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var recipe model.BaseRecipe
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&recipe)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := handler.recipeDao.AddRecipe(&recipe)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	addedRecipe, err := handler.recipeDao.GetRecipe(intID)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = handler.searchDao.IndexRecipe(*addedRecipe); err != nil {
		log.Println(err)
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
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, err := handler.recipeDao.UpdateRecipe(model.Recipe{
		ID:         mux.Vars(request)["id"],
		BaseRecipe: recipe,
	})
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = handler.searchDao.IndexRecipe(*updated); err != nil {
		log.Println(err)
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
	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = handler.recipeDao.DeleteRecipe(id); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = handler.searchDao.DeleteRecipe(strID); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusNoContent)
}
