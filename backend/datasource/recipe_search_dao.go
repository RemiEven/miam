package datasource

import (
	"fmt"
	"log/slog"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/fr"
	"github.com/blevesearch/bleve/v2/mapping"

	"github.com/remieven/miam/model"
)

// RecipeSearchDao struct
type RecipeSearchDao struct {
	index bleve.Index
}

// NewRecipeSearchDao creates a new recipe search dao
func NewRecipeSearchDao() (*RecipeSearchDao, error) {
	mapping := buildIndexMapping()
	recipeIndex, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bleve search engine: %w", err)
	}

	return &RecipeSearchDao{
		index: recipeIndex,
	}, nil
}

func buildIndexMapping() mapping.IndexMapping {
	frenchTextFieldMapping := bleve.NewTextFieldMapping()
	frenchTextFieldMapping.Analyzer = fr.AnalyzerName
	frenchTextFieldMapping.Store = false

	idTextFieldMapping := bleve.NewTextFieldMapping()
	idTextFieldMapping.Analyzer = keyword.Name
	idTextFieldMapping.IncludeInAll = false
	idTextFieldMapping.Store = false

	ingredientMapping := bleve.NewDocumentStaticMapping()
	ingredientMapping.AddFieldMappingsAt("id", idTextFieldMapping)
	ingredientMapping.AddFieldMappingsAt("name", frenchTextFieldMapping)

	recipeMapping := bleve.NewDocumentStaticMapping()
	recipeMapping.AddFieldMappingsAt("name", frenchTextFieldMapping)
	recipeMapping.AddFieldMappingsAt("howTo", frenchTextFieldMapping)
	recipeMapping.AddSubDocumentMapping("ingredients", ingredientMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("recipe", recipeMapping)
	indexMapping.DefaultMapping = recipeMapping

	return indexMapping
}

// IndexRecipe indexes a new or already existing recipe in the search engine
func (dao *RecipeSearchDao) IndexRecipe(recipe model.Recipe) error {
	return dao.index.Index(recipe.ID, recipe)
}

// DeleteRecipe deletes a recipe from the search engine
func (dao *RecipeSearchDao) DeleteRecipe(recipeID string) error {
	return dao.index.Delete(recipeID)
}

// SearchRecipes searches for recipes according to the given criteria
func (dao *RecipeSearchDao) SearchRecipes(search model.RecipeSearch) ([]string, int, error) {
	query := bleve.NewConjunctionQuery()
	if search.SearchTerm != "" {
		matchQuery := bleve.NewMatchPhraseQuery(search.SearchTerm)
		matchQuery.Analyzer = fr.AnalyzerName
		query.AddQuery(matchQuery)
	}
	if search.ExcludedRecipes != nil && len(search.ExcludedRecipes) != 0 {
		exclusionQuery := bleve.NewBooleanQuery()
		exclusionQuery.AddMustNot(bleve.NewDocIDQuery(search.ExcludedRecipes))
		query.AddQuery(exclusionQuery)
	}
	if search.ExcludedIngredients != nil && len(search.ExcludedIngredients) != 0 {
		exclusionQuery := bleve.NewBooleanQuery()
		for _, excluded := range search.ExcludedIngredients {
			excludeIngredientQuery := bleve.NewTermQuery(excluded)
			excludeIngredientQuery.SetField("ingredients.id")
			exclusionQuery.AddMustNot(excludeIngredientQuery)
		}
		query.AddQuery(exclusionQuery)
	}

	searchResults, err := dao.index.Search(bleve.NewSearchRequest(query))
	if err != nil {
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}

	ids := make([]string, len(searchResults.Hits))
	for i := range searchResults.Hits {
		ids[i] = searchResults.Hits[i].ID
	}

	return ids, int(searchResults.Total), nil
}

// Close closes the index used to search recipes
func (dao *RecipeSearchDao) Close() error {
	slog.Info("closing bleve search engine index")
	return dao.index.Close()
}
