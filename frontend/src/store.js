import Vue from 'vue'
import Vuex from 'vuex'

import ingredientApi from '@/api/ingredient'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    allIngredients: []
  },
  mutations: {
    setAllIngredients(state, {ingredients}) {
        state.allIngredients = ingredients
    }
  },
  actions: {
    async getAllIngredients({commit}) {
      const ingredients = await ingredientApi.getIngredients()
      commit('setAllIngredients', {ingredients})
    }
  }
})
