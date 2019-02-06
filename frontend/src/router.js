import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/ingredients',
      name: 'ingredients-admin',
      component: () => import('./views/IngredientAdmin.vue'),
    },
    {
      path: '/recipe/add',
      name: 'add-recipe-form',
      component: () => import('./components/AddRecipeForm.vue'),
    }
  ]
})
