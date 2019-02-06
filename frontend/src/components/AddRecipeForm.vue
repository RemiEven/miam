<template>
  <div>
    <form>
      <input type="text" id="nameInput" name="name" placeholder="Nom" autofocus v-model.trim="name">
      <textarea id="how-to-input" name="how-to" placeholder="Instructions" v-model.trim="howTo"></textarea>
      <button type="button" v-on:click="add" :disabled="!valid">Ajouter</button>
    </form>
  </div>
</template>

<script>

export default {
  name: 'add-recipe-form',
  data() {
      return {
          name: "",
          howTo: "",
      }
  },
  computed: {
      valid() {
          return !!this.name
      }
  },
  methods: {
    add() {
        this.$store.dispatch('addRecipe', {
          recipe: {
            name: this.name,
            howTo: this.howTo,
          }
        })
    }
  }
}
</script>
