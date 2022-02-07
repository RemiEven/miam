package rest

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/RemiEven/miam/service"
)

var defaultAllowedHosts = []string{"http://localhost:8080"}

// CreateRouter creates a new HTTP router
func CreateRouter(recipeService *service.RecipeService, ingredientService *service.IngredientService) http.Handler {
	router := mux.NewRouter()

	var (
		recipeHandler     = newRecipeHandler(recipeService)
		ingredientHandler = newIngredientHandler(ingredientService)
	)

	router.Use(handlers.CompressHandler)
	configureCORS(router)

	router.HandleFunc("/recipe", recipeHandler.AddRecipe).Methods(http.MethodPost)
	router.HandleFunc("/recipe/{id}", recipeHandler.GetRecipeByID).Methods(http.MethodGet)
	router.HandleFunc("/recipe/{id}", recipeHandler.UpdateRecipe).Methods(http.MethodPut)
	router.HandleFunc("/recipe/{id}", recipeHandler.DeleteRecipe).Methods(http.MethodDelete)
	router.HandleFunc("/recipe/search", recipeHandler.SearchRecipe).Methods(http.MethodPost)
	router.HandleFunc("/ingredient", ingredientHandler.GetIngredients).Methods(http.MethodGet)
	router.HandleFunc("/ingredient/{id}", ingredientHandler.UpdateIngredient).Methods(http.MethodPut)
	router.HandleFunc("/ingredient/{id}", ingredientHandler.DeleteIngredient).Methods(http.MethodDelete)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", SpaHandler{})).Methods(http.MethodGet)

	return router
}

func configureCORS(router *mux.Router) {
	router.Use(handlers.CORS(
		handlers.AllowedOrigins(defaultAllowedHosts),
		handlers.AllowedMethods([]string{
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		}),
		handlers.AllowedHeaders([]string{
			"Content-Type",
		}),
		handlers.ExposedHeaders([]string{
			"Location",
		}),
	))
	router.Use(mux.CORSMethodMiddleware(router))
	router.PathPrefix("/").HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {}).Methods(http.MethodOptions)
}
