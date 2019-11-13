<template>
  <div class="headerBar">
    <h1 class="title">
      Tinfoil Wizard
    </h1>
    <div>
      <button
      v-if="state.characterID"
      class="icon-button right-space"
      title="back to character select"
      @click="goBack"
    >
      <img src="@/assets/icons/back.svg">
    </button>
      <Profile />
    </div>
  </div>
</template>

<script>
import Profile from "./Profile.vue"

export default {
  name: 'Header',
  components: {
    Profile,
  },
  data() {
    return {
      state: this.$root.$data.state,
    }
  },
  methods: {
    goBack: function() {
      this.state.characterID = null;
    }
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.headerBar {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.title {
  font-family: 'MrEavesSmallCaps';
  font-size: 32px;
}

.right-space {
  margin-right: 8px;
}
</style>