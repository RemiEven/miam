<template>
  <div>
    <h1>{{ recipe.name }}</h1>
    <p>{{ recipe.howTo }}</p>
    <ul>
      <li v-for="ingredient in recipe.ingredients" :key="ingredient.id">
        {{ingredient.name}} {{ingredient.quantity ? `(${ingredient.quantity})` : ''}}
      </li>
    </ul>
  </div>
</template>

<script>
import recipeApi from '@/api/recipe'

export default {
  name: 'recipe',
  mounted() {
    this.getRecipe(this.$route.params.id)
  },
  watch: {
    '$route' (to, from) {
      this.getRecipe(to.params.id)
    }
  },
  computed: {
    recipe() {
      return this.$store.state.recipe
    }
  },
  methods: {
    getRecipe(recipeId) {
      this.$store.dispatch('setRecipe', {recipeId})
    }
  }
  // FIXME: https://router.vuejs.org/guide/essentials/dynamic-matching.html#reacting-to-params-changes
}
</script>

<style scoped>

</style>
