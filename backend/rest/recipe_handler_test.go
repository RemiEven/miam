package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remieven/miam/datasource"
	"github.com/remieven/miam/pb-lite/failure"
	"github.com/remieven/miam/pb-lite/fixture"
	"github.com/remieven/miam/pb-lite/testutils"
	"github.com/remieven/miam/service"
)

func TestGetRecipe(t *testing.T) {
	tests := map[string]struct {
		prepareDatabase  func(*datasource.DatabaseHolder) error
		recipeId         string
		expectedStatus   int
		responseBodyTest func(string) (string, bool)
	}{
		"recipe not found": {
			recipeId:         "1",
			expectedStatus:   http.StatusNotFound,
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.ResourceNotFoundErrorCode),
		},
		"invalid ID": {
			recipeId:         "not_a_valid_id",
			expectedStatus:   http.StatusBadRequest,
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.InvalidArgumentErrorCode),
		},
		"nominal case": {
			prepareDatabase: fixture.PrepareDatabase(
				`insert into ingredient(id, name) values
					(1, "riz"),
					(2, "haricots rouges"),
					(3, "purée de piment"),
					(4, "émincés de soja")
				`,
				`insert into recipe(id, name, how_to) values (1, "riz aux haricots rouges", "just prepare it")`,
				`insert into recipe_ingredient (recipe_id, ingredient_id, quantity) values
					(1, 1, ""),
					(1, 2, ""),
					(1, 3, "not too much"),
					(1, 4, "")
				`,
			),
			recipeId:       "1",
			expectedStatus: http.StatusOK,
			responseBodyTest: testutils.JsonResponseBodyTest(`{
				"id": "1",
				"name": "riz aux haricots rouges",
				"howTo": "just prepare it",
				"ingredients": [
					{
						"id": "1",
						"name": "riz"
					},
					{
						"id": "2",
						"name": "haricots rouges"
					},
					{
						"id": "3",
						"name": "purée de piment",
						"quantity": "not too much"
					},
					{
						"id": "4",
						"name": "émincés de soja"
					}
				]
			}`),
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

			request, err := http.NewRequest(http.MethodGet, "/recipe/"+test.recipeId, nil)
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
