package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/handler"
	"github.com/gorilla/mux"
)

const port = 8080

var (
	recipeDao *datasource.RecipeDao
)

func main() {
	log.Println("Starting") // FIXME this writes to stderr apparently

	datasourceContext, err := datasource.NewContext()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// defer datasourceContext.Close() // TODO this clashes with the os.Exit; get rid of those by extracting a method

	recipeHandler := handler.NewRecipeHandler(datasourceContext.RecipeDao)
	ingredientHandler := handler.NewIngredientHandler(datasourceContext.IngredientDao)

	router := mux.NewRouter()
	router.HandleFunc("/recipe", recipeHandler.AddRecipe).Methods(http.MethodPost)
	router.HandleFunc("/recipe/{id}", recipeHandler.GetRecipeByID).Methods(http.MethodGet)
	router.HandleFunc("/recipe/{id}", recipeHandler.UpdateRecipe).Methods(http.MethodPut)
	router.HandleFunc("/recipe/{id}", recipeHandler.DeleteRecipe).Methods(http.MethodDelete)
	router.HandleFunc("/ingredient/", ingredientHandler.GetIngredients).Methods(http.MethodGet)
	router.HandleFunc("/ingredient/{id}", ingredientHandler.UpdateIngredient).Methods(http.MethodPut)
	router.HandleFunc("/ingredient/{id}", ingredientHandler.DeleteIngredient).Methods(http.MethodDelete)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))) // FIXME: this is not suitable for SPA as is

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Printf("Will try to listen on port %d", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	var wait = time.Second * 15
	context, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(context)

	log.Println("Shutting down")
}
