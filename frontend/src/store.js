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
    search: {
      searchTerm: '',
      excludedRecipes: [],
      excludedIngredients: [],
    },
    searchResults: {
      total: 0,
      firstResults: [],
    },
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
    resetSearchResults(state) {
      state.searchResults.total = 0
      state.searchResults.firstResults = []
    },
    mergeSearchResults(state, {searchResults}) {
      state.searchResults.total = searchResults.total
      const currentResults = state.searchResults.firstResults
      const newResults = searchResults.firstResults
      state.searchResults.firstResults = [
        ...currentResults,
        ...newResults.filter(result => currentResults.every(currentResult => currentResult.id !== result.id))
      ]
    },
    setSearchTerm(state, {searchTerm}) {
      state.search.searchTerm = searchTerm
    },
    excludeIngredient(state, {id, name}) {
      state.search.excludedIngredients.push({id, name})
      state.searchResults.firstResults = state.searchResults.firstResults
          .filter(recipe => recipe.ingredients.every(ingredient => ingredient.id !== id))
    },
    includeIngredient(state, {ingredientId}) {
      state.search.excludedIngredients = state.search.excludedIngredients
          .filter(ingredient => ingredient.id != ingredientId)
    },
    excludeRecipe(state, {id, name}) {
      state.search.excludedRecipes.push({id, name})
      state.searchResults.firstResults = state.searchResults.firstResults
          .filter(recipe => recipe.id !== id)
    },
    includeRecipe(state, {recipeId}) {
      state.search.excludedRecipes = state.search.excludedRecipes
          .filter(recipe => recipe.id != recipeId)
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
    async deleteRecipe(_, {recipeId}) {
      await recipeApi.deleteRecipe(recipeId)
    },
    async search({state, commit}) {
      const searchRequest = {
        searchTerm: state.search.searchTerm,
        excludedRecipes: state.search.excludedRecipes.map(excluded => excluded.id),
        excludedIngredients: state.search.excludedIngredients.map(excluded => excluded.id),
      }
      const searchResults = await recipeApi.searchRecipe(searchRequest)
      commit('mergeSearchResults', {searchResults})
    },
    async excludeIngredient({commit, dispatch}, {id, name}) {
      commit('excludeIngredient', {id, name})
      dispatch('search')
    },
    async includeIngredient({commit, dispatch}, {ingredientId}) {
      commit('includeIngredient', {ingredientId})
      dispatch('search')
    },
    async excludeRecipe({commit, dispatch}, {id, name}) {
      commit('excludeRecipe', {id, name})
      dispatch('search')
    },
    async includeRecipe({commit, dispatch}, {recipeId}) {
      commit('includeRecipe', {recipeId})
      dispatch('search')
    },
  },
})
