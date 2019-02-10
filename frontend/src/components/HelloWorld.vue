<template>
  <div class="column col-12 empty">
    <div class="empty-icon icon-3x icon-resize-horiz icon"></div>
    <p class="empty-title h5">Aucune recette</p>
    <button v-on:click="goToAddRecipeForm" class="empty-action btn btn-primary">Ajouter une recette</button>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  methods: {
    goToAddRecipeForm() {
      this.$router.push({
        name: "add-recipe-form",
      })
    }
  },
}
</script>
