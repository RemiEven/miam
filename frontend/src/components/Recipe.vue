<template>
  <div>
    <h1>{{ recipe.name }}</h1>
    <button type="button" v-on:click="deleteRecipe" class="btn btn-error">Supprimer</button> <!-- TODO: maybe display delete only on edit page -->
    <h2>Ingrédients</h2>
    <ul>
      <li v-for="ingredient in recipe.ingredients" :key="ingredient.id">
        {{ingredient.name}} {{ingredient.quantity ? `(${ingredient.quantity})` : ''}}
      </li>
    </ul>
    <span v-if="compiledHowTo == ''" class="text-gray">Pas d'instructions disponibles</span>
    <div v-else>
      <h2>Instructions</h2> <!-- TODO: hide if no instructions -->
      <div v-html="compiledHowTo" />
    </div>
  </div>
</template>

<script>
import marked from 'marked'

export default {
  name: 'recipe',
  mounted() {
    this.getRecipe(this.$route.params.id)
  },
  watch: {
    '$route' (to) {
      this.getRecipe(to.params.id)
    }
  },
  computed: {
    recipe() {
      return this.$store.state.recipe
    },
    compiledHowTo() {
      return marked(((this.$store.state.recipe || {}).howTo) || '')
    },
  },
  methods: {
    getRecipe(recipeId) {
      this.$store.dispatch('setRecipe', {recipeId})
    },
    async deleteRecipe() {
      await this.$store.dispatch('deleteRecipe', {
        recipeId: this.$store.state.recipe.id
      })
      this.$router.push({
        name: 'home'
      })
    }
  }
  // FIXME: https://router.vuejs.org/guide/essentials/dynamic-matching.html#reacting-to-params-changes
}
</script>

<style scoped>

</style>
