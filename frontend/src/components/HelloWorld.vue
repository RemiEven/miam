<template>
  <div>
    <search-recipe-form></search-recipe-form>
    <div class="divider"></div>

    <div v-if="!!displayedRecipe" class="mb-2 mt-2">
       <recipe-tile v-bind:recipe="displayedRecipe"></recipe-tile>
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
    displayedRecipe() {
      return this.$store.state.recipe
    },
    searching() {
      const search = this.$store.state.search
      return search.searchTerm || search.excludedRecipes.length || search.excludedIngredients.length
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
    this.$store.dispatch('displayNewRecipe')
  },
}
</script>
