<template>
  <!--- row of spell info for tab -->
  <div class="row spell no-negging">
    <div class="col">
      <b-button
        v-if="listType === 'prep'"
        class="m-2 btn btn-dark"
        @click="castSpell()"
      >
        Cast
      </b-button>
      <b-button
        v-else-if="allClass === 'allclass'"
        class="m-2 btn btn-dark"
        @click="addSpell(spell.SpellKey)"
      >
        Add
      </b-button>
      <b-button
        v-if="allClass !== 'allclass' && listType !== 'prep' && isPrepared === 0 && spell.Level !== 0 && displayPrepare"
        class="m-2 btn btn-dark"
        @click="prepareSpell(spell.SpellKey, spell.Level)"
      >
        Prepare
      </b-button>
      <b-button
        v-if="allClass !== 'allclass' && listType !== 'prep' && isPrepared === 1 && spell.Level !== 0"
        class="btn btn-secondary"
        @click="unprepareSpell(spell.SpellKey, spell.Level)"
      >
        Prepared
      </b-button>
      <b-button
        v-if="(listType !== 'prep') && (listType !== 'allknown') && (allClass != 'allclass')"
        @click="castSpell()"
      >
        End Concentration
      </b-button>
    </div>
    <div class="col spell-name">
      {{ spell.SpellName }}
    </div>
    <div class="col spell-stat-col">
      {{ spell.CastingTime }}
    </div>
    <div class="col spell-stat-col">
      {{ spell.SpellRange }}
    </div>

    <!-- info modal for spell -->
    <div class="col">
      <b-button
        id="show-btn"
        class="m-2 btn btn-dark"
        @click="$bvModal.show(modalID)"
      >
        Info
      </b-button>
      <b-modal 
        :id="modalID"
        centered
        hide-header
        body-class="phbDialogDialogOverride"
        content-class="phbDialogContentOverride"
        dialog-class="dialogClass"
        hide-footer>
        <div class="container-border spell-dialog-container">
          <h3 class="title2">
            {{ spell.SpellName }}
          </h3>
          <p class="long_text">
            {{ spell.Description }}
          </p>
          <div class="spell-stat">
            Magic School: {{ spell.MagicSchool }}
          </div>
          <div class="spell-stat">
            Range: {{ spell.SpellRange }}
          </div>
          <div class="spell-stat">
            Duration: {{ spell.Duration }}
          </div>
          <div class="spell-stat">
            Casting Time: {{ spell.CastingTime }}
          </div>
          <div
            v-if="spell.VerbalComponent == 1 ? true: false"
            class="spell-stat"
          >
            Verbal Component
          </div>
          <div
            v-if="spell.SomaticComponent === 1 ? true: false"
            class="spell-stat"
          >
            Somatic Component
          </div>
          <div
            v-if="spell.MaterialComponent === ''? false: true"
            class="spell-stat"
          >
            Material Component: {{ spell.MaterialComponent }}
          </div>
          <div
            v-if="spell.RequiresConcentration == 1"
            class="spell-stat"
          >
            Requires Concentration
          </div>
          <b-button block class="btn btn-dark space-top" @click="$bvModal.hide(modalID)">
          Close
        </b-button>
        </div>
      </b-modal>
      <!-- END OF MODAL -->
    </div>
  </div>
</template>

<script>
import axios from "axios";

let modalID = 0;

export default {
  props: {
    spell: {
      type: Object,
      required: true
    },
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
    isPrepared: {
      type: Number,
      default: 0
    },
    displayPrepare: {
      type: Number,
      default: 1
    }
  },
  beforeCreate() {
    this.modalID = modalID.toString();
    modalID += 1;
  },
  methods: {
    castSpell() {
      this.$emit("cast-me", this.spell);
    },
    addSpell(spellKey) {
      axios.put("/api/characters/" + this.$root.$data.state.characterID + "/spells/" + spellKey);
      this.$emit("refresh-ui");
    },
    prepareSpell(spellKey, level) {
      this.$emit("prep-spell", {'spellKey':spellKey, 'level':level});
    },
    unprepareSpell(spellKey, level) {
      this.$emit("unprep-spell", {'spellKey':spellKey, 'level':level});
    }
  }
};
</script>

<style scoped>
.spell-dialog-container {
  padding: 0 18px;
}

.space-top {
  margin-top: 12px;
}
</style>
