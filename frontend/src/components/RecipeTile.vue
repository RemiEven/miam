<template>
 <div class="tile tile-centered" v-on:click="displayIngredients = !displayIngredients">
    <div class="tile-content">
    <div class="tile-title text-bold">{{recipe.name}}</div>
    <div class="tile-subtitle">
        <button v-on:click="goToRecipePage(recipe)" class="empty-action btn btn-secondary">Détails</button>
    </div>
    <div v-if="displayIngredients">
        <div class="divider text-center" data-content="Ingrédients"></div>
        <span v-for="ingredient in recipe.ingredients" :key="ingredient.id" v-on:click="excludeIngredient(ingredient)" class="chip">
            {{ingredient.name}}
            <button class="btn btn-clear" role="button" />
        </span>
    </div>
    </div>
    <div class="tile-action">
    <button class="btn btn-link btn-lg" v-on:click="excludeRecipe(recipe)"><i class="icon icon-cross"></i></button>
    </div>
</div>
</template>

<script>
export default {
  name: 'RecipeTile',
  props: ['recipe'],
  data() {
    return {
      displayIngredients: false,
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
