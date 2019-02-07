import Vue from 'vue'
import Vuex from 'vuex'

import ingredientApi from '@/api/ingredient'
import recipeApi from '@/api/recipe'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    allIngredients: [],
    addedRecipeId: '',
    recipe: null,
  },
  mutations: {
    setAllIngredients(state, {ingredients}) {
      state.allIngredients = ingredients
    },
    removeIngredient(state, {ingredientId}) {
      state.allIngredients = state.allIngredients
          .filter(ingredient => ingredient.id != ingredientId)
    },
    setRecipe(state, {recipe}) {
      state.recipe = recipe
    },
    setAddedRecipeId(state, {recipeId}) {
      state.addedRecipeId = recipeId
    },
  },
  actions: {
    async getAllIngredients({commit}) {
      const ingredients = await ingredientApi.getIngredients()
      commit('setAllIngredients', {ingredients})
    },
    async deleteIngredient({commit}, {ingredientId}) {
      await ingredientApi.deleteIngredient(ingredientId)
      commit('removeIngredient', {ingredientId})
    },
    async addRecipe(commit, {recipe}) {
      const recipeId = await recipeApi.addRecipe(recipe)
      commit.commit('setAddedRecipeId', {recipeId})
    },
    async setRecipe({commit}, {recipeId}) {
      const recipe = await recipeApi.getRecipe(recipeId)
      commit('setRecipe', {recipe})
    },
  },
})
