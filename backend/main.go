package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/handler"
	"github.com/RemiEven/miam/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const defaultPort = 7040

var defaultAllowedHosts = []string{"http://localhost:8080"}

var (
	recipeDao *datasource.RecipeDao
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.Info("Starting")

	viper.SetConfigName("configuration")
	viper.AddConfigPath(".")
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("allowedHosts", defaultAllowedHosts)
	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatal("Failed to read configuration file")
	}

	serviceContext, err := service.NewContext()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if err := serviceContext.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close context")
		}
	}()

	ids, err := serviceContext.GetDatasourceContext().RecipeDao.StreamRecipeIds()
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	for _, id := range ids {
		recipe, err := serviceContext.GetDatasourceContext().RecipeDao.GetRecipe(id)
		if err != nil {
			logrus.WithError(err).Fatal()
		}
		if recipe != nil {
			if err := serviceContext.GetDatasourceContext().RecipeSearchDao.IndexRecipe(*recipe); err != nil {
				logrus.WithError(err).Fatal()
			} else {
				logrus.WithField("id", id).Debug("indexed recipe")
			}
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

	port := viper.GetInt("port")

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		logrus.WithField("port", port).Info("Will try to start http server")
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logrus.Info("Closed http server")
			} else {
				logrus.Error(err)
				serviceContext.Close()
				os.Exit(1)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	var wait = time.Second * 15
	context, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(context)

	logrus.Info("Shutting down")
}

func configureCORS(router *mux.Router) {
	router.Use(handlers.CORS(
		handlers.AllowedOrigins(viper.GetStringSlice("allowedHosts")),
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
