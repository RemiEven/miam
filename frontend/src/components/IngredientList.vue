<template>
  <div>
    <div v-if="ingredients.length > 0" class="mb-2 mt-2">
      <div v-for="(ingredient, index) in ingredients" :key="ingredient.id">
        <div class="tile tile-centered">
          <div class="tile-content">
            <div class="tile-title text-bold">{{ingredient.name}}</div>
          </div>
          <div class="tile-action">
            <button type="button" v-on:click="deleteIngredient(ingredient.id)" class="btn btn btn-error btn-action btn-lg"><i class="icon icon-delete"></i></button>
          </div>
        </div>
        <div class="divider" v-if="index !== ingredients.length - 1"></div>
      </div>
    </div>
    <div v-else class="column col-12 empty">
      <div class="empty-icon icon-3x icon-resize-horiz icon"></div>
      <p class="empty-title h5">Aucun ingrédient</p>
      <p class="empty-subtitle">Ajoutez une recette pour ajouter des ingrédients</p>
      <button v-on:click="goToAddRecipeForm" class="empty-action btn btn-primary">Ajouter une recette</button>
    </div>
  </div>
</template>

<script>
export default {
  name: "IngredientList",
  computed: {
    ingredients() {
      return this.$store.state.allIngredients
        .sort((s1, s2) => (s1.name > s2.name ? 1 : -1))
    },
  },
  methods: {
    deleteIngredient(ingredientId) {
      this.$store.dispatch("deleteIngredient", { ingredientId })
    },
    goToAddRecipeForm() {
      this.$router.push({
        name: 'add-recipe-form',
      })
    },
  },
}
</script>
