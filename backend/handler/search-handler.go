package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/model"
)

// SearchHandler struct
type SearchHandler struct {
	recipeDao *datasource.RecipeDao
	searchDao *datasource.RecipeSearchDao
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(recipeDao *datasource.RecipeDao, searchDao *datasource.RecipeSearchDao) *SearchHandler {
	return &SearchHandler{
		recipeDao,
		searchDao,
	}
}

// SearchRecipe search for recipes
func (handler *SearchHandler) SearchRecipe(responseWriter http.ResponseWriter, request *http.Request) {
	var search model.RecipeSearch
	defer request.Body.Close()
	err := json.NewDecoder(request.Body).Decode(&search)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	results, err := handler.searchRecipe(search)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(responseWriter).Encode(results)
	}
}

func (handler *SearchHandler) searchRecipe(recipeSearch model.RecipeSearch) (*model.RecipeSearchResult, error) {
	if recipeSearch.IsEmpty() {
		return handler.recipeDao.GetRandomRecipes(recipeSearch)
	}
	strIDs, total, err := handler.searchDao.SearchRecipes(recipeSearch)
	if err != nil {
		return nil, err
	}
	IDs := make([]int, len(strIDs))
	for i := range strIDs {
		ID, err := strconv.Atoi(strIDs[i])
		if err != nil {
			return nil, err
		}
		IDs[i] = ID
	}
	recipes, err := handler.recipeDao.GetRecipes(IDs)
	if err != nil {
		return nil, err
	}
	return &model.RecipeSearchResult{
		FirstResults: recipes,
		Total:        total,
	}, nil
}
