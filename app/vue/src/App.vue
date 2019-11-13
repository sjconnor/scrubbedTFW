<template>
  <div
    id="app"
    class="scroll-background"
  >
    <Header />

    <!-- page-specific components -->

    <Login v-if="state.user == null" />
    <CharacterSelect v-if="state.jwt != null && state.characterID == null" />
    <Character v-if="state.characterID != null" />

    <!-- end page-specific components -->
  </div>
</template>

<style>
  @import './assets/style/style.css';
</style>

<script>

import Character from './components/Character.vue'
import CharacterSelect from './components/CharacterSelect.vue'
import Login from './components/Login.vue'
import Header from './components/Header.vue'

// note: App component doesn't have data; data is in components further in tree
export default {
  name: 'App',
  components: {
    Character,
    CharacterSelect,
    Login,
    Header,
  },
  data() {
    return this.$root.$data
  },
  methods: {
    toggleCollapse: function() {
      const bar = document.getElementById("navbarResponsive");
      if (bar.className.includes("show")) {
          bar.classList.remove("show");
      } else {
          bar.classList.add("show");
      }
    }
  }

}
</script>
