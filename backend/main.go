package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/RemiEven/miam/handler"
	"github.com/RemiEven/miam/service"
)

const defaultPort = 7040

var defaultAllowedHosts = []string{"http://localhost:8080"}

func main() {
	for _, err := range startApplication() {
		log.Error().Err(err).Msg("execution failed")
	}
}

func startApplication() (errors []error) {
	appendError := func(err error) {
		if err != nil {
			errors = append(errors, err)
		}
	}
	ctx := context.Background()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Starting")

	serviceContext, err := service.NewContext()
	if err != nil {
		appendError(err)
		return
	}
	defer func() { appendError(serviceContext.Close()) }()

	ids, err := serviceContext.GetDatasourceContext().RecipeDao.StreamRecipeIds(ctx)
	if err != nil {
		appendError(err)
		return
	}
	for _, id := range ids {
		recipe, err := serviceContext.GetDatasourceContext().RecipeDao.GetRecipe(ctx, id)
		if err != nil {
			appendError(err)
			return
		}
		if recipe == nil {
			continue
		}
		if err := serviceContext.GetDatasourceContext().RecipeSearchDao.IndexRecipe(*recipe); err != nil {
			appendError(err)
			return
		} else {
			log.Debug().Str("id", id).Msg("indexed recipe")
		}
	}

	recipeHandler := handler.NewRecipeHandler(serviceContext.RecipeService)
	ingredientHandler := handler.NewIngredientHandler(serviceContext.IngredientService)

	router := mux.NewRouter()
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

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", handler.SpaHandler{})).Methods(http.MethodGet)

	port := defaultPort

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Info().Int("port", port).Msg("will try to start http server")
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info().Msg("closed http server")
			} else {
				log.Error().Err(err)
				serviceContext.Close()
				os.Exit(1)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	wait := time.Second * 15
	shutdownCtx, cancel := context.WithTimeout(ctx, wait)
	defer cancel()

	srv.Shutdown(shutdownCtx)

	log.Info().Msg("shutting down")

	return
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
