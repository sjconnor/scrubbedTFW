<template>
  <div>
    <div
      class="new-char-button"
      @click="openModal"
    >
      <div>
        + new character
      </div>
    </div>
    <b-modal
      id="new_char"
      centered
      hide-header
      body-class="phbDialogDialogOverride"
      content-class="phbDialogContentOverride"
      dialog-class="dialogClass"
      hide-footer
    >
      <div>
        <div class="phbDialogScrollContainer">
          <svg
            class="phbDialogCornerScroll"
            height="10"
            width="20"
          >
            <polygon
              points="0,10 20,0 20,10"
              style="fill:#212529;"
            />
          </svg>
          <svg
            class="phbDialogCornerScroll"
            height="10"
            width="20"
          >
            <polygon
              points="0,10 0,0 20,10"
              style="fill:#212529;"
            />
          </svg>
        </div>
        <div class="phbDialogBody">
          <div class="phbDialogHead">
            <h2>Create character</h2>
            <button
              class="phbIconButton"
              title="cancel"
              @click="cancel"
            >
              <img src="@/assets/icons/clear.svg">
            </button>
          </div>
          <div class="inputs">
            <div
              class="portrait-container"
              @click="selectFile"
            >
              <div
                v-if="portrait == null"
                class="char-thumb"
              >
                <div class="file-select-instr">
                  select<br>portrait
                </div>
              </div>
              <img
                v-else
                id="portraitPreview"
                class="char-thumb"
                :src="portrait"
              >
              <img
                class="char-thumb-overlay"
                src="@/assets/icons/portraitFrame.svg"
              >
            </div>
            <div class="new-char-fields">
              <div class="phbTextInputContainer">
                <div class="phbInputLabel min-label">
                  Name:
                </div>
                <input
                  v-model="newChar.name"
                  class="phbTextInput"
                  tabindex="3"
                >
              </div>
              <div class="phbTextInputContainer">
                <div class="phbInputLabel min-label">
                  Class:
                </div>
                <select
                  v-model="newChar.class"
                  class="phbSelect"
                  @change="selectClass"
                  tabindex="4"
                >
                  <option
                    v-for="option in classOptions"
                    :key="option.id"
                    :value="option"
                  >
                    {{ option.name }}
                  </option>
                </select>
              </div>
              <div
                v-if="newChar.class !== null"
                class="phbTextInputContainer"
              >
                <div class="phbInputLabel min-label">
                  Sub:
                </div>
                <select
                  v-model="newChar.subclass"
                  class="phbSelect"
                  tabindex="5"
                >
                  <option
                    v-for="subclass in newChar.class.subclasses"
                    :key="subclass.id"
                    :value="subclass"
                  >
                    {{ subclass.name }}
                  </option>
                </select>
              </div>
              <div>
                <div
                  v-if="newChar.class !== null"
                  class="phbTextInputContainer"
                >
                  <div class="phbInputLabel min-label">
                    {{ newChar.class.spellcastingAbility }}:
                  </div>
                  <input
                    v-model.number="newChar.abilityScore"
                    class="phbTextInput"
                    tabindex="6"
                  >
                </div>
              </div>
            </div>
          </div>
          <div class="action-button">
            <button
              class="phbButton"
              @click="create"
              :disabled="!newChar.class || newChar.name == '' || !newChar.subclass || !newChar.abilityScore || !portrait"
            >
              create
            </button>
          </div>
        </div>
        <div class="phbDialogScrollContainer">
          <svg
            class="phbDialogCornerScroll"
            height="10"
            width="20"
          >
            <polygon
              points="0,0 20,0 20,10"
              style="fill:#212529;"
            />
          </svg>
          <svg
            class="phbDialogCornerScroll"
            height="10"
            width="20"
          >
            <polygon
              points="0,0 20,0 0,10"
              style="fill:#212529;"
            />
          </svg>
        </div>
      </div>
    </b-modal>
  </div>
</template>

<script>
import axios from "axios"

function baseState(state) {
  return {
      newChar: {
        name: "",
        class: null,
        subclass: null,
        abilityScore: null,
      },
      portrait: null,
      classOptions: null,
      state: state,
    };
}

export default {
  name: 'NewChar',
  props: {
    display: {
      type: String,
      default: null
    }
  },
  data() {
    return baseState(this.$root.$data.state);
  },
  mounted() {
    axios.get(`/api/classes?subclasses=true`)
    .then(response => {
      this.classOptions = response.data;
    })
    .catch(e => {
      window.console.log(e);
    })
  },
  methods: {
    cancel: function() {
      this.$bvModal.hide('new_char');
    },
    create: function() {
      this.newChar.level = 1;
      axios.post('/api/characters',
      {
        "character": this.newChar,
        "portrait": this.portrait,
      },
      {
        headers: {
          "Authorization": this.$root.$data.state.jwt,
          "Content-Type": "application/json",
        },
      }).then(res => {
        this.state.characterID = res.data.id;
      }).catch(err => {
        this.$bvToast.toast(`Couldn't create character: ${err}`,
        {
          toaster: 'b-toaster-bottom-center',
          variant: 'warning',
          noCloseButton: true,
        });
        window.console.log(err);
      })
    },
    openModal: function() {
      let base = baseState(this.$root.$data.state);
      base.classOptions = this.classOptions;
      Object.assign(this.$data, base);
      this.$bvModal.show('new_char');
    },
    selectClass: function() {
      this.newChar.subclass = null;
    },
    selectFile: function() {
      // from s/o: https://stackoverflow.com/questions/16215771/how-open-select-file-dialog-via-js/16215950
      const input = document.createElement('input');
      input.type = 'file';
      input.accept = 'image/*';

      input.onchange = e => {
        const file = e.target.files[0];

        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = e => {
          this.portrait =  e.target.result;
        };
        reader.onerror = e => {
          this.$bvToast.toast(`Couldn't load portrait. Try again.`,
              {
                toaster: 'b-toaster-bottom-center',
                variant: 'warning',
                noCloseButton: true,
              });
            window.console.error(err);
        };
      }
      input.click();
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h2 {
  font-family: 'MrEavesSmallCaps';
  font-size: 24px !important;
}

.action-button {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  margin-top: 12px;
}

.new-char-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  line-height: 40px;
  cursor: pointer;
  margin-left: auto;
  margin-top: 24px;
  margin-right: auto;
  border: solid;
  border-image: url('../assets/icons/frameBorderFilled.svg') repeat;
  border-image-slice: 40 60 fill;
  border-image-width: 40px 60px;
  padding: 24px;
  max-width: 540px;
  overflow: hidden;
  font-family: 'ScalySansCaps';
  font-size: 26px;
  min-height: 80px;
}

.new-char-button div {
  height: 40px;
}

.inputs {
  display: flex;
  flex-direction: row;
}

.new-char-fields {
  flex-grow: 2;
}

.min-label {
  min-width: 52px !important;
}

.portrait-container {
  position: relative;
  margin-right: 12px;
  cursor: pointer;
}

.char-thumb {
  width: 128px;
  height: 128px;
  margin: 4px 0;
  border-radius: 50%;
  background-color: #e0e5c1;
}

.char-thumb-overlay {
  position: absolute;
  top: -4px;
  left: -8px;
}

.file-select-instr {
  padding-top:38px;
  text-align: center;
}
</style>
