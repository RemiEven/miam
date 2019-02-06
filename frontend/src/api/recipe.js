const backend = 'http://localhost:7040' // TODO: put this in config

// Proper way would be sending a GET request to the server, but this will do
const extractRecipeIdFromLocation = location =>  location.split('/').pop();

export default {
    async addRecipe(recipe) {
        const request = {
            method: 'POST',
            body: JSON.stringify(recipe),
            headers: new Headers({
                "Content-Type": "application/json:charset=UTF_8"
            }),
        }
        const response = await fetch(`${backend}/recipe`, request)
        return extractRecipeIdFromLocation(response.headers.get("Location"))
    }
}
