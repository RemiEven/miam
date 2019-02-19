<template>
  <div>
    <search-recipe-form></search-recipe-form>
    <div class="divider"></div>

    <div v-if="numberOfResults > 0" class="mb-2 mt-2">
      <div v-for="(recipe, index) in firstResults" :key="recipe.id">
       <recipe-tile v-bind:recipe="recipe"></recipe-tile>
        <div class="divider" v-if="index !== firstResults.length - 1"></div>
      </div>
    </div>
    <div v-else class="column col-12 empty">
      <div class="empty-icon icon-3x icon-resize-horiz icon"></div>
      <p class="empty-title h5">Aucune recette {{ searching ? 'trouvée' : '' }}</p>
      <button v-on:click="goToAddRecipeForm" type="button" class="empty-action btn btn-primary">Ajouter une recette</button>
    </div>
  </div>
</template>

<script>
import SearchRecipeForm from '@/components/SearchRecipeForm.vue'
import RecipeTile from '@/components/RecipeTile.vue'

export default {
  name: 'HelloWorld',
  components: {
    SearchRecipeForm,
    RecipeTile,
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
        name: 'add-recipe-form',
      })
    },
  },
  mounted() {
    this.$store.dispatch('search')
  },
}
</script>
