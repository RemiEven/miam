import Vue from 'vue'
import Vuex from 'vuex'

import ingredientApi from '@/api/ingredient'
import recipeApi from '@/api/recipe'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    allIngredients: []
  },
  mutations: {
    setAllIngredients(state, {ingredients}) {
      state.allIngredients = ingredients
    },
    removeIngredient(state, {ingredientId}) {
      state.allIngredients = state.allIngredients
          .filter(ingredient => ingredient.id != ingredientId)
    }
  },
  actions: {
    async getAllIngredients({commit}) {
      const ingredients = await ingredientApi.getIngredients()
      commit('setAllIngredients', {ingredients})
    },
    async deleteIngredient({commit}, {ingredientId}) {
      await ingredientApi.deleteIngredient(ingredientId)
      commit("removeIngredient", {ingredientId})
    },
    async addRecipe({commit}, {recipe}) {
      console.log(recipe)
      const recipeId = await recipeApi.addRecipe(recipe)
    },
  },
})
