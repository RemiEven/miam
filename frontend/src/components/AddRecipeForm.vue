<template>
  <form class="form-horizontal" onsubmit="return false">
    <input type="text" id="nameInput" name="name" placeholder="Nom" autofocus v-model.trim="name" class="form-group form-input" />
    <textarea id="how-to-input" name="how-to" placeholder="Instructions" v-model.trim="howTo" class="form-group form-input" rows="5" />
    <div v-for="ingredient in ingredients" :key="ingredient.localId" class="form-group columns">
      <span class="column col-5">
        <input type="text" v-bind:id="'ingredientInput' + ingredient.localId" v-bind:name="'ingredient' + ingredient.localId" placeholder="Nom" v-model.trim="ingredient.name" class="form-input" />
      </span>
      <span class="column col-5">
        <input type="text" v-bind:id="'quantityInput' + ingredient.localId" v-bind:name="'quantity' + ingredient.localId" placeholder="Quantity" v-model.trim="ingredient.quantity" class="form-input" />
      </span>
      <span class="column col-2">
        <button type="button" v-on:click="removeIngredientInput(ingredient.localId)" class="btn btn-action btn-large btn-error"><i class="icon icon-cross"></i></button>
      </span>
    </div>
    <button type="button" v-on:click="addIngredientInput" class="btn btn-secondary form-group mi-btn">Ajouter un ingrédient</button>
    <button type="submit" v-on:click="add" :disabled="!valid" class="btn btn-success form-group mi-btn">Ajouter la recette</button>
  </form>
</template>

<script>
import { notBlank } from '@/utils'

// TODO: also handle already existing ingredients with their id
export default {
  name: 'add-recipe-form',
  data() {
    return {
      name: "",
      howTo: "",
      ingredients: [],
      lastIngredientLocalId: 0,
    }
  },
  computed: {
    valid() {
      return notBlank(this.name) && this.ingredients.every(({id, name}) => notBlank(id) || notBlank(name))
    }
  },
  methods: {
    async add() {
      await this.$store.dispatch('addRecipe', {
        recipe: {
          name: this.name,
          howTo: this.howTo,
          ingredients: this.ingredients
              .map(({id, name, quantity}) => ({id, name, quantity})),
        }
      })
      this.$router.push({
        name: 'recipe',
        params: {
          id: this.$store.state.addedRecipeId
        },
      })
    },
    addIngredientInput() {
      this.ingredients.push({
        name: "",
        quantity: "",
        localId: this.lastIngredientLocalId,
      })
      this.lastIngredientLocalId++
    },
    removeIngredientInput(ingredientLocalId) {
      this.ingredients = this.ingredients
          .filter(ingredient => ingredient.localId != ingredientLocalId)
    },
  }
}
</script>

<style scoped>
.mi-btn {
  width: 100%;
  justify-content: center;
}
</style>
