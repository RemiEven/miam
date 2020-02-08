package service

import "github.com/RemiEven/miam/datasource"

// Context struct
type Context struct {
	datasourceContext *datasource.Context
	RecipeService     *RecipeService
	IngredientService *IngredientService
}

// NewContext creates a new service context
func NewContext() (*Context, error) {
	datasourceContext, err := datasource.NewContext()
	if err != nil {
		return nil, err
	}
	return &Context{
		datasourceContext,
		newRecipeService(datasourceContext.RecipeDao, datasourceContext.RecipeSearchDao),
		newIngredientService(datasourceContext.IngredientDao, datasourceContext.RecipeIngredientDao),
	}, nil
}

func (context *Context) GetDatasourceContext() *datasource.Context {
	return context.datasourceContext
}

// Close cleanly closes the context by releasing any resources it might hold
func (context *Context) Close() error {
	return context.datasourceContext.Close()
}
