<template>
  <div>
    <form>
      <input type="text" id="nameInput" name="name" placeholder="Nom" autofocus v-model.trim="name" />
      <textarea id="how-to-input" name="how-to" placeholder="Instructions" v-model.trim="howTo" />
      <div v-for="ingredient in ingredients" :key="ingredient.localId">
        <input type="text" v-bind:id="'ingredientInput' + ingredient.localId" v-bind:name="'ingredient' + ingredient.localId" placeholder="Nom" v-model.trim="ingredient.name" />
        <input type="text" v-bind:id="'quantityInput' + ingredient.localId" v-bind:name="'quantity' + ingredient.localId" placeholder="Quantity" v-model.trim="ingredient.quantity" />
        <button type="button" v-on:click="removeIngredientInput(ingredient.localId)">Retirer l'ingrédient</button>
      </div>
      <button type="button" v-on:click="addIngredientInput">Ajouter un ingrédient</button>
      <button type="button" v-on:click="add" :disabled="!valid">Ajouter</button>
    </form>
  </div>
</template>

<script>
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
          // FIXME: also test that all ingredients either have non blank id or non blank name
          return !!this.name
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
        }
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
