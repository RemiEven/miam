<template>
  <div>
    <form onsubmit="return false" class="form-horizontal">
      <div class="input-group">
        <input type="text" id="searchTermInput" name="searchTerm" placeholder="Chercher" autofocus v-model.trim="searchTerm" class="form-input" />
        <button type="submit" v-on:click="search" class="input-group-btn btn btn-success"><i class="icon icon-search"></i></button>
      </div>
    </form>
    <div v-if="excludedIngredients.length > 0" class="divider text-center" data-content="Ingrédients exclus"></div>
    <span v-for="ingredient in excludedIngredients" :key="ingredient.id" v-on:click="includeIngredient(ingredient.id)" class="chip">
      {{ingredient.name}}
      <button class="btn btn-clear" role="button" />
    </span>
    <div v-if="excludedRecipes.length > 0" class="divider text-center" data-content="Recette exclues"></div>
    <span v-for="recipe in excludedRecipes" :key="recipe.id" v-on:click="includeRecipe(recipe.id)" class="chip">
      {{recipe.name}}
      <button class="btn btn-clear" role="button" />
    </span>
  </div>
</template>

<script>
export default {
  name: 'search-recipe-form',
  computed: {
    searchTerm: {
      get() {
        return this.$store.state.search.searchTerm
      },
      set(searchTerm) {
        return this.$store.commit('setSearchTerm', {searchTerm})
      },
    },
    excludedRecipes() {
      return this.$store.state.search.excludedRecipes
    },
    excludedIngredients() {
      return this.$store.state.search.excludedIngredients
    },
  },
  methods: {
    includeRecipe(recipeId) {
      this.$store.dispatch('includeRecipe', {recipeId})
    },
    includeIngredient(ingredientId) {
      this.$store.dispatch('includeIngredient', {ingredientId})
    },
    search() {
      this.$store.dispatch('search')
    },
  },
}
</script>
