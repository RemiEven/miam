<template>
  <div>
    <search-recipe-form></search-recipe-form>
    <div class="divider"></div>

    <div v-if="numberOfResults > 0" class="mb-2 mt-2">
      <div v-for="(recipe, index) in firstResults" :key="recipe.id">
        <div class="tile tile-centered">
          <div class="tile-content">
            <router-link :to="linkToRecipe(recipe)" class="tile-title text-bold" tag="div">{{recipe.name}}</router-link>
          </div>
          <div class="tile-action">
            <button class="btn btn-link btn-action btn-lg"><i class="icon icon-more-vert"></i></button>
          </div>
        </div>
        <div class="divider" v-if="index !== firstResults.length - 1"></div>
      </div>
    </div>
    <div v-else class="column col-12 empty">
      <div class="empty-icon icon-3x icon-resize-horiz icon"></div>
      <p class="empty-title h5">Aucune recette {{ searching ? 'trouvée' : '' }}</p>
      <button v-on:click="goToAddRecipeForm" class="empty-action btn btn-primary">Ajouter une recette</button>
    </div>
  </div>
</template>

<script>
import SearchRecipeForm from '@/components/SearchRecipeForm.vue'

export default {
  name: 'HelloWorld',
  components: {
    SearchRecipeForm,
  },
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
