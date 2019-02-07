<template>
    <div>
        <ul>
            <li v-for="ingredient in ingredients" :key="ingredient.id">
                {{ingredient.name}}
                <button type="button" v-on:click="deleteIngredient(ingredient.id)">Supprimer</button>
            </li>
        </ul>
    </div>
</template>

<script>
import ingredientApi from '@/api/ingredient'

export default {
    name: 'IngredientList',
    computed: {
        ingredients() {
            return this.$store.state.allIngredients
        }
    },
    methods: {
        async deleteIngredient(ingredientId) {
            this.$store.dispatch("deleteIngredient", {ingredientId})
        }
    }
}
</script>
