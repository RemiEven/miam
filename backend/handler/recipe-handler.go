package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/RemiEven/miam/model"

	"github.com/RemiEven/miam/datasource"
	"github.com/gorilla/mux"
)

// RecipeHandler is a recipe handler
type RecipeHandler struct {
	recipeDao *datasource.RecipeDao
}

// NewRecipeHandler creates a new recipe handler
func NewRecipeHandler(dao *datasource.RecipeDao) *RecipeHandler {
	return &RecipeHandler{
		dao,
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
	} else {
		responseWriter.Header().Add("Location", id)
		responseWriter.WriteHeader(http.StatusCreated)
	}
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
	} else {
		json.NewEncoder(responseWriter).Encode(updated)
	}
}

func (handler *RecipeHandler) DeleteRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = handler.recipeDao.DeleteRecipe(id); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		responseWriter.WriteHeader(http.StatusNoContent)
	}
}

func (handler *RecipeHandler) SearchRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var search model.RecipeSearch
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&search)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := handler.recipeDao.SearchRecipes(search)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(responseWriter).Encode(results)
	}
}
