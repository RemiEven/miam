package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RemiEven/miam/datasource"
	"github.com/RemiEven/miam/service"
	"github.com/RemiEven/miam/testutils"
)

func TestGetIngredients(t *testing.T) {
	tests := map[string]struct {
		prepareDatabase  func(*datasource.DatabaseHolder) error
		expectedStatus   int
		responseBodyTest func(string) (string, bool)
	}{
		"no ingredients": {
			expectedStatus:   http.StatusOK,
			responseBodyTest: testutils.JsonResponseBodyTest(`[]`),
		},
		"one ingredient": {
			prepareDatabase: func(holder *datasource.DatabaseHolder) error {
				_, err := holder.DB.Exec(`insert into ingredient(name) values("salade verte")`)
				return err
			},
			expectedStatus:   http.StatusOK,
			responseBodyTest: testutils.JsonResponseBodyTest(`[{"name":"salade verte","id":"1"}]`),
		},
		"several ingredients": {
			prepareDatabase: func(holder *datasource.DatabaseHolder) error {
				_, err := holder.DB.Exec(`insert into ingredient(name) values("salade verte")`)
				if err != nil {
					return err
				}
				_, err = holder.DB.Exec(`insert into ingredient(name) values("oignon rouge")`)
				return err
			},
			expectedStatus:   http.StatusOK,
			responseBodyTest: testutils.JsonResponseBodyTest(`[{"name":"salade verte","id":"1"},{"name":"oignon rouge","id":"2"}]`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dbFilePath := testutils.GetRandomDBFileName()
			databaseHolder, err := datasource.NewDatabaseHolder(dbFilePath)
			if err != nil {
				t.Error(err)
				return
			}

			defer func() {
				if err := databaseHolder.Close(); err != nil {
					t.Error(err)
				}
			}()

			ingredientDao, err := datasource.NewIngredientDao(databaseHolder)
			if err != nil {
				t.Error(fmt.Errorf("failed to initialize ingredientDao: %w", err))
				return
			}
			recipeIngredientDao, err := datasource.NewRecipeIngredientDao(databaseHolder, ingredientDao)
			if err != nil {
				t.Error(fmt.Errorf("failed to initialize recipeIngredientDao: %w", err))
				return
			}
			recipeDao, err := datasource.NewRecipeDao(databaseHolder, recipeIngredientDao)
			if err != nil {
				t.Error(fmt.Errorf("failed to initialize recipeDao: %w", err))
				return
			}
			recipeSearchDao, err := datasource.NewRecipeSearchDao()
			if err != nil {
				t.Error(fmt.Errorf("failed to initialize recipeSearchDao: %w", err))
				return
			}
			defer func() {
				if err := recipeSearchDao.Close(); err != nil {
					t.Error(err)
				}
			}()

			if test.prepareDatabase != nil {
				if err := test.prepareDatabase(databaseHolder); err != nil {
					t.Errorf("failed to prepare database: %v", err)
					return
				}
			}

			var (
				ingredientService = service.NewIngredientService(ingredientDao, recipeIngredientDao)
				recipeService     = service.NewRecipeService(recipeDao, recipeSearchDao)
			)

			ctx := context.Background()
			if err := recipeService.IndexAllExistingRecipes(ctx); err != nil {
				t.Error(fmt.Errorf("failed to index recipes: %w", err))
				return
			}

			router := CreateRouter(recipeService, ingredientService)

			request, err := http.NewRequest(http.MethodGet, "/ingredient", nil)
			if err != nil {
				t.Error(err)
				return
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, request)

			if rr.Result().StatusCode != test.expectedStatus {
				t.Errorf("unexpected statusCode: wanted [%d], got [%d]", test.expectedStatus, rr.Result().StatusCode)
			}

			responseBody, err := io.ReadAll(rr.Result().Body)
			if err != nil {
				t.Error(err)
				return
			}
			if msg, ok := test.responseBodyTest(string(responseBody)); !ok {
				t.Error(msg)
			}
		})
	}
}
