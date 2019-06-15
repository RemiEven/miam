const backend = process.env.VUE_APP_API_BASE_URL

export default {
    async getIngredients() {
        const response = await fetch(`${backend}/ingredient`);
        return await response.json();
    },

    async deleteIngredient(ingredientId) {
        const request = {
            method: 'DELETE',
        }
        return await fetch(`${backend}/ingredient/${ingredientId}`, request)
    }
}
