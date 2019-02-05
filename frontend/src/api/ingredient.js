const backend = 'http://localhost:7040'

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
