<template>
  <div>
    <ul v-if="ingredients.length > 0">
      <li v-for="ingredient in ingredients" :key="ingredient.id" class="chip">
        {{ingredient.name}}
        <button type="button" v-on:click="deleteIngredient(ingredient.id)" class="btn btn-clear"></button>
      </li>
    </ul>
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
      return this.$store.state.allIngredients;
    }
  },
  methods: {
    deleteIngredient(ingredientId) {
      this.$store.dispatch("deleteIngredient", { ingredientId });
    },
    goToAddRecipeForm() {
      this.$router.push({
        name: "add-recipe-form"
      });
    },
  }
};
</script>
