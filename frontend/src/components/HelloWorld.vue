<template>
  <div>
    <ul v-if="numberOfResults !== 0">
      <router-link :to="linkToRecipe(recipe)" v-for="recipe in firstResults" :key="recipe.id" tag="li"><a>{{recipe.name}}</a></router-link>
    </ul>
    <div v-if="numberOfResults === 0" class="column col-12 empty">
      <div class="empty-icon icon-3x icon-resize-horiz icon"></div>
      <p class="empty-title h5">Aucune recette {{ searching ? 'trouvée' : '' }}</p>
      <button v-on:click="goToAddRecipeForm" class="empty-action btn btn-primary">Ajouter une recette</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  computed: {
    numberOfResults() {
      return this.$store.state.searchResults.total
    },
    searching() {
      const search = this.$store.state.search
      return search.searchTerm || search.excludedRecipes.length || search.excludedIngredients.length
    },
    firstResults() {
      return this.$store.state.searchResults.firstResults
    },
  },
  methods: {
    goToAddRecipeForm() {
      this.$router.push({
        name: "add-recipe-form",
      })
    },
    linkToRecipe(recipe) {
      return `/recipe/${recipe.id}`
    },
  },
  mounted() {
    this.$store.dispatch('search')
  },
}
</script>
