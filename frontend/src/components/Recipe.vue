<template>
  <div>
    <h1>{{ recipe.name }}</h1>
    <h2>Ingrédients</h2>
    <ul>
      <li v-for="ingredient in recipe.ingredients" :key="ingredient.id">
        {{ingredient.name}} {{ingredient.quantity ? `(${ingredient.quantity})` : ''}}
      </li>
    </ul>
    <h2>Instructions</h2>
    <div v-html="compiledHowTo" />
  </div>
</template>

<script>
import marked from 'marked'

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
    },
    compiledHowTo() {
      return marked((this.$store.state.recipe || {}).howTo)
    },
  },
  methods: {
    getRecipe(recipeId) {
      this.$store.dispatch('setRecipe', {recipeId})
    },
  }
  // FIXME: https://router.vuejs.org/guide/essentials/dynamic-matching.html#reacting-to-params-changes
}
</script>

<style scoped>

</style>
