<template>
  <div id="app" class="container">
    <ul class="tab tab-block">
      <router-link to="/" tag="li" class="tab-item" exact><a>Recettes</a></router-link>
      <router-link to="/recipe/add" tag="li" class="tab-item"><a>Ajouter une recette</a></router-link>
      <router-link to="/ingredients" tag="li" class="tab-item"><a>Ingrédients</a></router-link>
    </ul>
    <router-view/>
  </div>
</template>

<style>

</style>
