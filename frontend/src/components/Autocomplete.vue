<template>
  <div>
    <div class="has-icon-right">
      <input @blur="displayDropdown = false" @focus="displayDropdown = true" type="text" placeholder="Nom" v-model.trim="text" class="form-input" />
      <i class="form-icon icon text-gray" v-bind:class="{'icon-bookmark': suggestionSelected, 'icon-plus': !suggestionSelected}"></i>
    </div>
    <ul v-if="displayDropdown" class="menu mi-suggestion-dropdown">
      <li class="menu-item" v-for="suggestion in filteredSuggestions" v-bind:key="suggestion.id" v-on:click="selectSuggestion(suggestion)">
        {{suggestion.name}}
      </li>
    </ul>
  </div>
</template>

<script>
function slug(string) {
  return string.toLowerCase().trim() // TODO: use better slug method
}

function matchSlugsStart(string1, string2) {
  return slug(string1).startsWith(slug(string2))
}

function matchSlugs(string1, string2) {
  return slug(string1) === slug(string2)
}

export default {
  name: "autocomplete",
  props: ["suggestions"],
  data() {
    return {
      displayDropdown: false,
      value: {
        name: "",
        id: ""
      },
    }
  },
  computed: {
    filteredSuggestions() {
      return this.suggestions
        .filter(suggestion => matchSlugsStart(suggestion.name, this.text))
        .sort((s1, s2) => (s1.name > s2.name ? 1 : -1))
        .slice(0, 10)
    },
    text: {
      get() {
        return this.value.name
      },
      set(text) {
        const matchedSuggestions = this.suggestions
          .filter(suggestion => matchSlugs(suggestion.name, text))
          .sort((s1, s2) => (s1.name > s2.name ? 1 : -1))
        if (matchedSuggestions.length !== 0) {
          this.selectSuggestion(matchedSuggestions[0])
        } else {
          this.value.name = text
          this.value.id = ""
        }
      },
    },
    suggestionSelected() {
      return !!this.value.id
    },
  },
  methods: {
    selectSuggestion({id, name}) {
      this.value = {id, name}
    },
  },
}
</script>

<style scoped>
.mi-suggestion-dropdown {
  position: absolute;
}
</style>
