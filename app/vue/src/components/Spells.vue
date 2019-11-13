<template>
  <!-- conc-on-me passes on the spellkey emitted by the Spell component onwards to the Character component -->
  <div id="spells">
    <div v-if="spells && spells.length">
      <h4 v-if="listLevel == 0"
          class="border-bottom align-items-center spell-level-header">
        Cantrips
      </h4>
      <h4 v-else
          class="border-bottom align-items-center spell-level-header">
        Level {{ listLevel }}
      </h4>
      <Spell
        v-for="spell in spells"
        :key="spell.SpellKey"
        :list-type="listType"
        list-level="listLevel"
        :all-class="allClass"
        :spell="spell"
        :is-prepared="spell.Prepared"
        class="row border-bottom no-negging"
        @cast-me="$emit('cast-me', $event)"
        @refresh-ui="$emit('refresh-ui')"
        @prep-spell="prepareSpell($event)"
        @unprep-spell="unprepareSpell($event)"
      />
    </div>
  </div>
</template>

<script>
import Spell from "./Spell.vue";
import axios from "axios";

export default {
  name: "Spells",
  components: {
    Spell
  },
  props: {
    listLevel: {
      type: String,
      default: null
    },
    listType: {
      type: String,
      default: null
    },
    allClass: {
      type: String,
      default: null
    },
    componentKey: {
      type: Number,
      default: 0
    },
  },
  data() {
    return {
      spells: null
    };
  },
  mounted() {
    switch (this.allClass) {
      case "allclass": // add spells tab
        axios
          .get("/api/spells?characterid=" + this.$root.$data.state.characterID)
          .then(response => {
             switch (this.listLevel) {
                  case "0":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 0;
                    });
                    break;
                  case "1":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 1;
                    });
                    break;
                  case "2":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 2;
                    });
                    break;
                  case "3":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 3;
                    });
                    break;
                  case "4":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 4;
                    });
                    break;
                  case "5":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 5;
                    });
                    break;
                  case "6":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 6;
                    });
                    break;
                  case "7":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 7;
                    });
                    break;
                  case "8":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 8;
                    });
                    break;
                  case "9":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 9;
                    });
                    break;
                }
          });
        break;
      default:
        axios
          .get("/api/characters/" + this.$root.$data.state.characterID + "/spells", {
            headers: {
              Authorization: this.$root.$data.state.jwt,
            },
          })
          .then(response => {
            switch (this.listType) {
              case "prep": // prepared spells tab
                switch (this.listLevel) {
                  case "0":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 0;
                    });
                    break;
                  case "1":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 1;
                    });
                    break;
                  case "2":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 2;
                    });
                    break;
                  case "3":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 3;
                    });
                    break;
                  case "4":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 4;
                    });
                    break;
                  case "5":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 5;
                    });
                    break;
                  case "6":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 6;
                    });
                    break;
                  case "7":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 7;
                    });
                    break;
                  case "8":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 8;
                    });
                    break;
                  case "9":
                    this.spells = response.data.filter(function(s) {
                      return s.Prepared === 1 && s.Level === 9;
                    });
                    break;
                }

                break;

              case "allknown": // known spells tab
                switch (this.listLevel) {
                  case "0":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 0;
                    });
                    break;
                  case "1":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 1;
                    });
                    break;
                  case "2":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 2;
                    });
                    break;
                  case "3":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 3;
                    });
                    break;
                  case "4":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 4;
                    });
                    break;
                  case "5":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 5;
                    });
                    break;
                  case "6":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 6;
                    });
                    break;
                  case "7":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 7;
                    });
                    break;
                  case "8":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 8;
                    });
                    break;
                  case "9":
                    this.spells = response.data.filter(function(s) {
                      return s.Level === 9;
                    });
                    break;
                }
                break;

              default:
                window.console.log("default spell: shouldn't hit this");
                break;
            }
          });
    }
  },
  methods: {
    prepareSpell(e) {
      for (var i = 0; i < this.spells.length; i++) {
        if (this.spells[i].SpellKey === e.spellKey) {
          this.spells[i].Prepared = 1;
        }
      }
      this.$emit("prep-spell", e);
    },
    unprepareSpell(e) {
      for (var i = 0; i < this.spells.length; i++) {
        if (this.spells[i].SpellKey === e.spellKey) {
          this.spells[i].Prepared = 0;
        }
      }
      this.$emit("unprep-spell", e);
    }
  }
};
</script>

<style scoped>
h4 {
  font-family: 'MrEavesSmallCaps';
  font-size: 20px;
}
</style>
