package datasource

import (
	"github.com/RemiEven/miam/model"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/fr"
	"github.com/blevesearch/bleve/mapping"
)

const indexPath = "miam.bleve"

// RecipeSearchDao struct
type RecipeSearchDao struct {
	index bleve.Index
}

func newRecipeSearchDao() (*RecipeSearchDao, error) {
	beerIndex, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := buildIndexMapping()
		beerIndex, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &RecipeSearchDao{
		index: beerIndex,
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

func (dao *RecipeSearchDao) IndexRecipe(recipe model.Recipe) error {
	return dao.index.Index(recipe.ID, recipe)
}

func (dao *RecipeSearchDao) DeleteRecipe(recipeID string) error {
	return dao.index.Delete(recipeID)
}

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
		return nil, 0, err
	}

	ids := make([]string, len(searchResults.Hits))
	for i := range searchResults.Hits {
		ids[i] = searchResults.Hits[i].ID
	}

	return ids, int(searchResults.Total), nil
}