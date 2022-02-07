package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/rest"
	"github.com/RemiEven/miam/service"
)

const defaultPort = 7040

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

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Starting")

	databaseHolder, err := datasource.NewDatabaseHolder("./miam.db")
	if err != nil {
		appendError(fmt.Errorf("failed to create database holder: %w", err))
		return
	}
	defer func() { appendError(databaseHolder.Close()) }()

	ingredientDao, err := datasource.NewIngredientDao(databaseHolder)
	if err != nil {
		appendError(fmt.Errorf("failed to initialize ingredientDao: %w", err))
		return
	}
	recipeIngredientDao, err := datasource.NewRecipeIngredientDao(databaseHolder, ingredientDao)
	if err != nil {
		appendError(fmt.Errorf("failed to initialize recipeIngredientDao: %w", err))
		return
	}
	recipeDao, err := datasource.NewRecipeDao(databaseHolder, recipeIngredientDao)
	if err != nil {
		appendError(fmt.Errorf("failed to initialize recipeDao: %w", err))
		return
	}
	recipeSearchDao, err := datasource.NewRecipeSearchDao()
	if err != nil {
		appendError(fmt.Errorf("failed to initialize recipeSearchDao: %w", err))
		return
	}
	defer func() { appendError(recipeSearchDao.Close()) }()

	var (
		ingredientService = service.NewIngredientService(ingredientDao, recipeIngredientDao)
		recipeService     = service.NewRecipeService(recipeDao, recipeSearchDao)
	)

	ctx := context.Background()
	if err := recipeService.IndexAllExistingRecipes(ctx); err != nil {
		appendError(fmt.Errorf("failed to index recipes: %w", err))
		return
	}

	router := rest.CreateRouter(recipeService, ingredientService)

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
