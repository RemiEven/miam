<template>
  <div>
  <form onsubmit="return false" class="form-horizontal">
    <input type="text" id="searchTermInput" name="searchTerm" placeholder="Chercher" autofocus v-model.trim="searchTerm" class="form-group form-input" />
    <button type="submit" v-on:click="search" class="btn btn-success form-group form-input mi-btn">Chercher</button>
  </form>
    <div class="divider text-center" data-content="Ingrédients exclus"></div>
    <!-- TODO if no ingredient is excluded, state it here -->
        <span v-for="ingredient in excludedIngredients" :key="ingredient.id" class="chip">
            {{ingredient.name}}
            <button class="btn btn-clear" role="button" v-on:click="includeIngredient(ingredient.id)" />
        </span>
    <div class="divider text-center" data-content="Recette exclues"></div>
        <span v-for="recipe in excludedRecipes" :key="recipe.id" class="chip">
            {{recipe.name}}
            <button class="btn btn-clear" role="button" v-on:click="includeRecipe(recipe.id)" />
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
      this.$store.commit('includeRecipe', {recipeId})
    },
    includeIngredient(ingredientId) {
      this.$store.commit('includeIngredient', {ingredientId})
    },
    search() {
      this.$store.dispatch('search')
    },
  },
}

</script>

<style scoped>
.mi-btn {
  width: 100%;
  justify-content: center;
}
</style>
