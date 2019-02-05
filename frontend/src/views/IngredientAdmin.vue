<template>
  <div>
    Ingredients :
    <IngredientList />
  </div>
</template>

<script>
import IngredientList from '@/components/IngredientList.vue'

export default {
  name: 'ingredient-admin',
  components: {
    IngredientList
  },
  mounted() {
    this.getIngredients()
  },
  methods: {
    getIngredients() {
        this.$store.dispatch('getAllIngredients')
    }
  }
}
</script>
