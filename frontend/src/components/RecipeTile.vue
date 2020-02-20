<template>
<div class="card">
  <div class="card-header">
        <button class="btn btn-link btn-lg float-right" v-on:click="excludeRecipe(recipe)"><i class="icon icon-cross"></i></button>
    <div class="card-title h5">{{recipe.name}}</div>
  </div>
<div class="card-image">
  <img class="img-responsive" src="http://localhost:7040/static/1.jpg">
</div>
  <div class="card-body">

    <span v-for="ingredient in recipe.ingredients" :key="ingredient.id" v-on:click="excludeIngredient(ingredient)" class="chip">
        {{ingredient.name}}
        <button class="btn btn-clear" role="button" />
    </span>
  </div>
    <div class="card-footer">
        <button v-on:click="goToRecipePage(recipe)" class="empty-action btn btn-secondary">Détails</button>
    </div>
</div>
</template>

<script>
export default {
  name: 'RecipeTile',
  props: ['recipe'],
  data() {
    return {
    }
  },
  methods: {
    goToRecipePage() {
      this.$router.push({
        name: 'recipe',
        params: {
          id: this.recipe.id,
        },
      })
    },
    excludeRecipe() {
      this.$store.dispatch('excludeRecipe', {
        id: this.recipe.id,
        name: this.recipe.name,
      })
    },
    excludeIngredient(ingredient) {
      this.$store.dispatch('excludeIngredient', {
        id: ingredient.id,
        name: ingredient.name,
      })
    },
  }
}

</script>
