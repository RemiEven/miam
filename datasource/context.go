package datasource

type Context struct {
	holder              *databaseHolder
	IngredientDao       *IngredientDao
	RecipeIngredientDao *RecipeIngredientDao
	RecipeDao           *RecipeDao
}

func NewContext() (*Context, error) {
	holder, err := newDatabaseHolder()
	if err != nil {
		return nil, err
	}
	// TODO also register process end hook to stop connection to db
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
	return &Context{
		holder,
		ingredientDao,
		recipeIngredientDao,
		recipeDao,
	}, nil
}

// Close cleanly closes the context by releasing any resources it might have
func (context *Context) Close() error {
	return context.holder.Close()
}
