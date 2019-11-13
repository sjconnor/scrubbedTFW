<template>
  <div>
    <h2>Characters</h2>
    <div
      v-if="loading"
      class="loading-message"
    >
      Summoning your characters<span class="dot1">.</span><span class="dot2">.</span><span class="dot3">.</span>
    </div>
    <transition name="fade">
      <div v-if="!loading">
        <div
          v-for="(character, index) in state.characters"
          :key="character.id"
          class="char-row"
          @click="setChar(character.id)"
        >
          <div class="portrait-container">
            <img
              class="char-thumb"
              :src="character.portraitPath"
            >
            <img
              class="char-thumb-overlay"
              src="@/assets/icons/portraitFrame.svg"
            >
          </div>
          <div class="char-details">
            <div>{{ character.name }}</div>
            <div>Level {{ character.level }} {{ character.class.name }}</div>
            <div><span v-if="character.subclass">{{ character.subclass.name }}</span></div>
          </div>
          <div
            class="delete-container"
            @click.prevent="setChar"
          >
            <div class="left-adjust">
              <DeleteCharDialog :char-index="index" />
            </div>
          </div>
        </div>
        <div
          v-if="state.characters == null || state.characters.length < 10"
          class="new-char-block"
        >
          <NewChar />
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
import axios from "axios"
import NewChar from "./NewChar.vue"
import DeleteCharDialog from "./DeleteCharDialog.vue"

export default {
  name: 'CharacterSelect',
  components: {
    NewChar,
    DeleteCharDialog,
  },
  data() {
    return {
      loading: true,
      showNewCharDialog: false,
      state: this.$root.$data.state,
    }
  },
  mounted() {
    axios.get(`/api/characters`, {
      headers: {
        Authorization: this.state.jwt,
      },
    })
    .then(response => {
      this.state.characters = response.data;
      this.loading = false;
    })
    .catch(e => {
      window.console.log(e);
    })
  },
  methods: {
    setChar: function(chId) {
      this.state.characterID = chId;
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h2 {
  font-family: 'MrEavesSmallCaps';
  font-size: 28px;
  text-align: center;
  margin: 12px;
}

h3 {
  font-family: 'ScalySansCaps';
  margin: 0;
  height: 28px;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}

.char-details {
  display: flex;
  flex-direction: column;
  justify-content:space-around;
  box-sizing: content-box;
  overflow-y: hidden;
  padding: 12px 0 12px 24px;
  font-family: 'ScalySansCaps';
  font-size: 20px;
  flex-grow: 1;
}

.char-details div {
  line-height: 27px;
  height: 26px;
  margin: 4px 0;
  text-overflow:ellipsis;
  overflow: hidden;
  white-space: nowrap;
  }

.char-details div:first-child {
  font-size: 26px;
  font-weight:500;
}

.char-row {
  display: flex;
  flex-direction: row;
  cursor: pointer;
  margin-left: auto;
  margin-top: 24px;
  margin-right: auto;
  border: solid;
  border-image: url('../assets/icons/frameBorder.svg') repeat;
  border-image-slice: 60 fill;
  border-image-width: 60px;
  padding: 24px;
  max-width: 540px;
  overflow: hidden;
}

.portrait-container {
  position: relative;
}

.char-thumb {
  width: 128px;
  height: 128px;
  margin: 4px 0;
  border-radius: 50%;
}

.char-thumb-overlay {
  position: absolute;
  top: -4px;
  left: -8px;
}

.delete-container {
  width: 0px;
  position: relative;
}

.delete-container .left-adjust {
  position: absolute;
  left: -24px;
  top: -4px;
}
</style>

<style>
.char-row:hover .hover-hide {
  opacity: 1;
}

.hover-hide {
  opacity: 0;
  transition: opacity .5s;
}
</style>
