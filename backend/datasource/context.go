package datasource

// Context struct
type Context struct {
	holder              *databaseHolder
	IngredientDao       *IngredientDao
	RecipeIngredientDao *RecipeIngredientDao
	RecipeDao           *RecipeDao
	RecipeSearchDao     *RecipeSearchDao
}

// NewContext creates a new datasource context
func NewContext() (*Context, error) {
	holder, err := newDatabaseHolder()
	if err != nil {
		return nil, err
	}
	ingredientDao, err := newIngredientDao(holder)
	if err != nil {
		return nil, err
	}
	recipeIngredientDao, err := newRecipeIngredientDao(holder, ingredientDao)
	if err != nil {
		return nil, err
	}
	recipeDao, err := newRecipeDao(holder, recipeIngredientDao)
	if err != nil {
		return nil, err
	}
	recipeSearchDao, err := newRecipeSearchDao()
	if err != nil {
		return nil, err
	}

	return &Context{
		holder,
		ingredientDao,
		recipeIngredientDao,
		recipeDao,
		recipeSearchDao,
	}, nil
}

// Close cleanly closes the context by releasing any resources it might have
func (context *Context) Close() error {
	return context.holder.Close()
}
