<template>
  <div>
    <div v-if="loading" class="loading-message">
      Traversing planes<span class="dot1">.</span><span class="dot2">.</span><span class="dot3">.</span>
    </div>
    <transition name="fade">
      <div v-if="!loading" id="charinfo" class="char-info">
      <!--- CHARACTER INFORMATION  -->
      <div v-if="charInfo" class="container-border char-main top-spacer">
        <div class="portrait-container">
          <img id="portrait" class="char-thumb" :src="charInfo.portraitPath"/>
          <img class="char-thumb-overlay" src="@/assets/icons/characterPortraitFrame.svg"/>
        </div>
        <div class="char-details">
          <div>
            <div class="char-name-container">
              <h2>{{ charInfo.name }}</h2>
            </div>
            <div class="char-level-container">
              <h3>Level {{ charInfo.level }} {{ charInfo.class.name }}<span v-if="charInfo.subclass">, <span class="no-break">{{ charInfo.subclass.name }}</span></span></h3>
            </div>
          </div>
          <div class="char-stat-container">
            <div>
              <div>Spell Save DC:</div>
              <div>{{ spellSaveDC(charInfo.level, charInfo.abilityScore) }}</div>
            </div>
            <div>
              <div>Spell Attack Mod:</div>
              <div>{{ spellAttackMod(charInfo.level, charInfo.abilityScore) }}</div>
            </div>
          </div>
          <!--- CHARACTER BUTTONS-->
          <div class="char-buttons">
            <b-button class="btn btn-dark" @click="longRest">
              Long Rest
            </b-button>
            <b-button v-if="charInfo.class.name == 'Wizard'" class="btn btn-dark" @click="arcaneRecovery">
              Arcane Recovery
            </b-button>
            <b-button id="show-btn" class="btn btn-dark" @click="edit">
              Edit
            </b-button>
            </div>
          </div>
        </div>
        <!-- CHARACTER EDIT MODAL -->
        <b-modal
          id="edit_char"
          centered
          hide-header
          body-class="phbDialogDialogOverride"
          content-class="phbDialogContentOverride"
          dialog-class="dialogClass"
          hide-footer>
          <div class="container-border edit-dialog-container">
            <div style="display: flex;">
              <div class="edit-thumb-container">
              <img
                class="edit-thumb"
                :src="editInfo.portraitPath"
              >
              <img
                class="edit-thumb-overlay"
                src="@/assets/icons/portraitFrame.svg"
              >
            </div>
              <div class="edit-left-space">
                <div class="phbTextInputContainer">
                  <div class="phbInputLabel min-label">
                    Name:
                  </div>
                  <input v-model="charInfo.name"
                        class="phbTextInput">
                </div>
                <div class="phbTextInputContainer">
                  <div class="phbInputLabel min-label">
                    Level:
                  </div>
                  <input v-model.number="charInfo.level"
                        class="phbTextInput">
                </div>
                <div class="phbTextInputContainer">
                  <div class="phbInputLabel min-label">
                    Score:
                  </div>
                  <input v-model.number="charInfo.abilityScore"
                        class="phbTextInput">
                </div>
              </div>
            </div>
            <div class="edit-buttons">
              <b-button
                class="btn btn-dark"
                block
                @click="selectFile"
              >Edit Portrait</b-button>
              <b-button
                class="btn btn-dark"
                block
                @click="cancel"
              >Cancel</b-button>
              <b-button
                class="btn btn-dark"
                block
                @click="save"
              >Save</b-button>
            </div>
          </div>
        </b-modal>
      <!-- SPELL SLOTS -->
      <div class="container-border spell-slots-container top-spacer">
        <h2>Spell Slots</h2>
        <div class="col">
          <div class="spell_slots">
            <div class="slot row no-negging">
              <div
                id="L1Slots"
                class="col"
              >
                Level 1: {{ charInfo.Level1SlotsRemaining }}
              </div>
              <div
                id="L4Slots"
                class="col"
              >
                Level 4: {{ charInfo.Level4SlotsRemaining }}
              </div>
              <div
                id="L7Slots"
                class="col"
              >
                Level 7: {{ charInfo.Level7SlotsRemaining }}
              </div>
            </div>
            <div class="slot row no-negging">
              <div
                id="L2Slots"
                class="col"
              >
                Level 2: {{ charInfo.Level2SlotsRemaining }}
              </div>
              <div
                id="L5Slots"
                class="col"
              >
                Level 5: {{ charInfo.Level5SlotsRemaining }}
              </div>
              <div
                id="L8Slots"
                class="col"
              >
                Level 8: {{ charInfo.Level8SlotsRemaining }}
              </div>
            </div>
            <div class="slot row no-negging">
              <div
                id="L3Slots"
                class="col"
              >
                Level 3: {{ charInfo.Level3SlotsRemaining }}
              </div>
              <div
                id="L6Slots"
                class="col"
              >
                Level 6: {{ charInfo.Level6SlotsRemaining }}
              </div>
              <div
                id="L9Slots"
                class="col"
              >
                Level 9: {{ charInfo.Level9SlotsRemaining }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- CONCENTRATION-->
      <div v-if="concSpell" class="container-border top-spacer">
        <h2>Concentrating</h2>
        <Spell
          :key="concSpell.SpellKey"
          :spell="concSpell"
          @cast-me="do_cast(null)"
          :displayPrepare="0"
        />
      </div>

      <!--- CONC WARN -->
      <span>
        <b-modal
          id="conc_warn"
          hide-footer
          hide-header
          centered
          body-class="phbDialogDialogOverride"
          content-class="phbDialogContentOverride"
          dialog-class="dialogClass"
          class="transparent"
        >
          <div class="container-border conc-dialog">
            <h2 v-if="concSpell">
              You are already concentrating on {{ concSpell.SpellName }}
            </h2>
            <div class="conc-dialog-buttons">
              <b-button
                class="m-1"
                block
                @click="conc_cancel"
              >Cancel</b-button>
              <b-button
                class="m-1"
                block
                @click="conc_continue"
              >Continue</b-button>
            </div>
          </div>
        </b-modal>
      </span>

      <!--- SPELL TABS -->
      <div class="container-border top-spacer">
        <ul
          class="nav nav-tabs"
          role="tablist"
        >
          <li class="nav-item">
            <a
              id="prepared-spells-tab"
              class="rounded-0 nav-link active"
              data-toggle="tab"
              href="#prepared-spells"
              role="tab"
              aria-controls="prepared-spells"
              aria-selected="true"
            >Prepared Spells</a>
          </li>
          <li class="nav-item">
            <a
              id="known-spells-tab"
              class="rounded-0 nav-link"
              data-toggle="tab"
              href="#known-spells"
              role="tab"
              aria-controls="known-spells"
              aria-selected="false"
            >Known Spells</a>
          </li>
          <li class="nav-item">
            <a
              id="add-spell-tab"
              class="rounded-0 nav-link"
              data-toggle="tab"
              href="#add-spell"
              role="tab"
              aria-controls="add-spell"
              aria-selected="false"
            >Add Spell</a>
          </li>
        </ul>

        <div class="row border-bottom align-items-center h5 no-negging">
          <div class="col" />
          <div class="col title3">
            Name
          </div>
          <div class="col title3">
            Time
          </div>
          <div class="col title3">
            Range
          </div>
          <div class="col" />
        </div>

        <div class="tab-content pt-3">
          <!-- PREPARED SPELLS TAB -->
          <div
            id="prepared-spells"
            class="tab-pane fade show active text-left"
            role="tabpanel"
            aria-labelledby="prepared-spells-tab"
          >
            <!--- CANTRIPS -->
            <Spells
              :key="'prepCantrips' + preparedkey"
              :list-type="listTypePrep"
              list-level="0"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />
            <!--- 1st LEVEL -->
            <Spells
              :key="'prep1' + preparedkey"
              :list-type="listTypePrep"
              list-level="1"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />

            <!--- 2nd LEVEL -->
            <Spells
              :key="'prep2' + preparedkey"
              :list-type="listTypePrep"
              list-level="2"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />

            <!--- 3rd LEVEL -->
            <Spells
              :key="'prep3' + preparedkey"
              :list-type="listTypePrep"
              list-level="3"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />

            <!--- 4th LEVEL -->
            <Spells
              :key="'prep4' + preparedkey"
              :list-type="listTypePrep"
              list-level="4"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />

            <!--- 5th LEVEL -->
            <Spells
              :key="'prep5' + preparedkey"
              :list-type="listTypePrep"
              list-level="5"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />

            <!--- 6th LEVEL -->
            <Spells
              :key="'prep6' + preparedkey"
              :list-type="listTypePrep"
              list-level="6"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />

            <!--- 7th LEVEL -->
            <Spells
              :key="'prep7' + preparedkey"
              :list-type="listTypePrep"
              list-level="7"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />

            <!--- 8th LEVEL -->
            <Spells
              :key="'prep8' + preparedkey"
              :list-type="listTypePrep"
              list-level="8"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />

            <!--- 9th LEVEL -->
            <Spells
              :key="'prep9' + preparedkey"
              :list-type="listTypePrep"
              list-level="9"
              @cast-me="do_cast"
              @refresh-ui="refreshUI"
            />
          </div>

          <!-- KNOWN SPELLS TAB -->
          <div
            id="known-spells"
            class="tab-pane fade show"
            role="tabpanel"
            aria-labelledby="known-spells-tab"
          >
            <!--- CANTRIPS -->
            <Spells
              :key="'knownCantrips' + compkey"
              :list-type="listTypeAllKnown"
              list-level="0"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />

            <!--- 1st Level -->
            <Spells
              :key="'known1' + compkey"
              :list-type="listTypeAllKnown"
              list-level="1"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
            <!--- 2nd Level -->
            <Spells
              :key="'known2' + compkey"
              :list-type="listTypeAllKnown"
              list-level="2"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
            <!--- 3rd Level -->
            <Spells
              :key="'known3' + compkey"
              :list-type="listTypeAllKnown"
              list-level="3"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
            <!--- 4th Level -->
            <Spells
              :key="'known4' + compkey"
              :list-type="listTypeAllKnown"
              list-level="4"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
            <!--- 5th Level -->
            <Spells
              :key="'known5' + compkey"
              :list-type="listTypeAllKnown"
              list-level="5"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
            <!--- 6th Level -->
            <Spells
              :key="'known6' + compkey"
              :list-type="listTypeAllKnown"
              list-level="6"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
            <!--- 7th Level -->
            <Spells
              :key="'known7' + compkey"
              :list-type="listTypeAllKnown"
              list-level="7"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
            <!--- 8th Level -->
            <Spells
              :key="'known8' + compkey"
              :list-type="listTypeAllKnown"
              list-level="8"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
            <!--- 9th Level -->
            <Spells
              :key="'known9' + compkey"
              :list-type="listTypeAllKnown"
              list-level="9"
              @refresh-ui="refreshUI"
              @all-spells="knownSpellsPop"
              @prep-spell="prepareSpell"
              @unprep-spell="unprepareSpell"
            />
          </div>

          <!-- ADD SPELLS TAB -->
          <div
            id="add-spell"
            class="tab-pane fade"
            role="tabpanel"
            aria-labelledby="add-spell-tab"
          >
            <Spells
              :key="'allCantrips' + compkey"
              :all-class="allClass"
              list-level="0"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all1' + compkey"
              :all-class="allClass"
              list-level="1"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all2' + compkey"
              :all-class="allClass"
              list-level="2"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all3' + compkey"
              :all-class="allClass"
              list-level="3"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all4' + compkey"
              :all-class="allClass"
              list-level="4"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all5' + compkey"
              :all-class="allClass"
              list-level="5"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all6' + compkey"
              :all-class="allClass"
              list-level="6"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all7' + compkey"
              :all-class="allClass"
              list-level="7"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all8' + compkey"
              :all-class="allClass"
              list-level="8"
              @refresh-ui="refreshUI"
            />
            <Spells
              :key="'all9' + compkey"
              :all-class="allClass"
              list-level="9"
              @refresh-ui="refreshUI"
            />
          </div>

        <!-- END OF ALL SPELLS TABS -->
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
import Spells from "./Spells.vue";
import Spell from "./Spell.vue";
import axios from "axios";

export default {
  name: "Character",
  components: {
    Spells,
    Spell
  },
  data() {
    return {
      loading: true,
      compkey: 0,
      preparedkey: 0,
      charInfo: {
        // make null when ready to use axios instead
        name: null,
        level: null,
        class: {
          name: null
        },
        subclass: {
          name: null
        },
        portraitPath: null,
        abilityScore: null,
        concentrating: null
      },
      editInfo: {
        name: null,
        level: null,
        abilityScore: null,
        concentrating: null, // to handle cancelling of concentration
        portraitPath: null
      },
      portrait: null,
      concSpell: null,
      nextConcSpell: null,
      listTypePrep: "prep",
      listTypeAllKnown: "allknown",
      allClass: "allclass",
      componentKey: 0,
      spells: null,
      level1SpellsRemaining: 0,
      level2SpellsRemaining: 0,
      level3SpellsRemaining: 0,
      level4SpellsRemaining: 0,
      level5SpellsRemaining: 0,
      level6SpellsRemaining: 0,
      level7SpellsRemaining: 0,
      level8SpellsRemaining: 0,
      level9SpellsRemaining: 0,
      spellLevel: ""
    };
  },
  mounted() {
    axios
      .get("/api/characters/" + this.$root.$data.state.characterID, {
        headers: {
          Authorization: this.$root.$data.state.jwt,
        },
      })
      .then(response => {
        this.loading = false;
        this.reloadCharacter(response);
      })
      .catch(err => {
        window.console.log(err);
      });
  },
  methods: {
    setSSDC() {
      return 5;
    },
    do_cast(e) {
      if (e === null) {
        // remove concentration
        this.charInfo.concentrating = 0;
        this.concSpell = null; // remove the concentrated spell component
        this.save();
      } else {
        if (e.RequiresConcentration && this.concSpell !== null) {
          this.editInfo.concentrating = e.SpellKey;
          this.nextConcSpell = e;
          this.$bvModal.show("conc_warn");
        } else {
          this.continue_cast(e);
        }
      }
    },
    continue_cast(e) {
      var concSpellKey = 0;
      if (this.concSpell != null) {
        concSpellKey = this.concSpell.SpellKey;
      }
      if (e.RequiresConcentration) {
        concSpellKey = e.SpellKey;
      }
      axios
        .patch(
          `/api/characters/` + this.$root.$data.state.characterID,
          {
            character: {
              id: this.$root.$data.state.characterID,
              concentrating: concSpellKey,
              Level1SlotsRemaining: -3, // Cast a spell
              Level2SlotsRemaining: e.Level, // Spell Level
            },
          },
          {
            headers: {
              Authorization: this.$root.$data.state.jwt,
              "Content-Type": "application/json"
            }
          }
        )
        .then( response => {
          return axios
                  .get("/api/characters/" + this.$root.$data.state.characterID,
                    {
                      headers: {
                        Authorization: this.$root.$data.state.jwt,
                      }
                    })
        })
        .then( response => {
          if (!response.data.CastSuccess) {
            this.$bvToast.toast(`Cast failed: no slots left!`,
              {
                toaster: 'b-toaster-top-center',
                variant: 'warning',
                noCloseButton: true,
                toastClass: 'lower-toast',
                bodyClass: 'center-toast-text',
                autoHideDelay: 2000,
              });
            }
          this.reloadCharacter(response);
        })
        .catch(err => {
          // TODO: handle creation errors.
          window.console.log(err);
        });
    },
    conc_continue(e) {
      // continue to apply concentration
      this.$bvModal.hide("conc_warn");
      this.continue_cast(this.nextConcSpell); // focus on new spell
    },
    conc_cancel() {
      // cancel applying concentration
      this.$bvModal.hide("conc_warn");
    },
    edit() {
      this.$bvModal.show("edit_char");
      // store old vals because v-model will overwrite even if cancel
      this.old_name = this.charInfo.name;
      this.editInfo.level = this.charInfo.level;
      this.editInfo.abilityScore = this.charInfo.abilityScore;
      this.editInfo.portraitPath = this.charInfo.portraitPath;
    },
    save() {
      axios
        .patch(
          `/api/characters/` + this.$root.$data.state.characterID,
          {
            character: {
              id: this.$root.$data.state.characterID,
              name: this.charInfo.name, // editable from modal
              class: this.charInfo.class,
              subclass: this.charInfo.subclass,
              level: this.charInfo.level, // editable from modal
              abilityScore: this.charInfo.abilityScore, // editable from modal
              portraitPath: this.charInfo.portraitPath,
              concentrating: this.charInfo.concentrating
            },
            portrait: this.portrait // the new thing
          },
          {
            headers: {
              Authorization: this.$root.$data.state.jwt,
              "Content-Type": "application/json"
            }
          }
        )
        .then( response => {
          return axios
                .get("/api/characters/" + this.$root.$data.state.characterID,
                    {
                        headers: {
                            Authorization: this.$root.$data.state.jwt,
                        }
                    })
        })
        .then( response => {
          this.reloadCharacter(response);
        })
        .catch(err => {
          // TODO: handle creation errors.
          window.console.log(err);
        });

      this.$bvModal.hide("edit_char");
      // refresh the portrait if it was changed
      if (this.portrait != null) {
        this.charInfo.portraitPath = this.portrait;
      }
    },
    cancel() {
      // reset input to old values
      this.charInfo.name = this.old_name;
      this.charInfo.level = this.editInfo.level;
      this.charInfo.abilityScore = this.editInfo.abilityScore;

      this.$bvModal.hide("edit_char");
    },
    selectFile: function() {
      // from s/o: https://stackoverflow.com/questions/16215771/how-open-select-file-dialog-via-js/16215950
      const input = document.createElement("input");
      input.type = "file";
      input.accept = 'image/*';

      input.onchange = e => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = readerEvent => {
          this.portrait = readerEvent.target.result;
          this.editInfo.portraitPath = this.portrait;
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
      };
      input.click();
    },
    refreshUI() {
      this.componentKey += 1;
      this.compkey += 1;
      this.refreshPrepared += 1;
      this.preparedkey += 1;
    },
    dorefreshPrepared() {
        this.preparedkey += 1;
    },
    arcaneRecovery() {
      axios
        .patch(
          `/api/characters/` + this.$root.$data.state.characterID,
          {
            character: {
              id: this.$root.$data.state.characterID,
              Level1SlotsRemaining: -2, // Have the backend do arcane recovery
            },
          },
          {
            headers: {
              Authorization: this.$root.$data.state.jwt,
              "Content-Type": "application/json"
            }
          }
        )
        .then( response => {
          return axios
                  .get("/api/characters/" + this.$root.$data.state.characterID,
                    {
                      headers: {
                        Authorization: this.$root.$data.state.jwt,
                      }
                    })
        })
        .then( response => {
          this.reloadCharacter(response);
        })
        .catch(err => {
          // TODO: handle creation errors.
          window.console.log(err);
        });
    },
    longRest() {
      axios
        .patch(
          `/api/characters/` + this.$root.$data.state.characterID,
          {
            character: {
              id: this.$root.$data.state.characterID,
              Level1SlotsRemaining: -1, // Reset the slots
            },
          },
          {
            headers: {
              Authorization: this.$root.$data.state.jwt,
              "Content-Type": "application/json"
            }
          }
        )
        .then( response => {
          return axios
                  .get("/api/characters/" + this.$root.$data.state.characterID,
                    {
                      headers: {
                        Authorization: this.$root.$data.state.jwt,
                      }
                    })
        })
        .then( response => {
          this.reloadCharacter(response);
        })
        .catch(err => {
          // TODO: handle creation errors.
          window.console.log(err);
        });
    },
    knownSpellsPop(e) {
      this.spells = e;
    },
    prep_cancel() {
      this.$bvModal.hide("prep_warn");
    },
    prepareSpell(e) {
      axios
        .patch(
          `/api/characters/` +
            this.$root.$data.state.characterID +
            `/spells/` +
            e["spellKey"],
          {
            Prepared: 1
          },
          {
            headers: {
              Authorization: this.$root.$data.state.jwt,
              "Content-Type": "application/json"
            }
          }
        )
        .catch(err => {
          // TODO: handle creation errors.
          window.console.log(err);
        })
        .then(response => {
          this.dorefreshPrepared();
        });
    },
    unprepareSpell(e) {
      axios
        .patch(
          `/api/characters/` +
            this.$root.$data.state.characterID +
            `/spells/` +
            e["spellKey"],
          {
            Prepared: 0
          },
          {
            headers: {
              Authorization: this.$root.$data.state.jwt,
              "Content-Type": "application/json"
            }
          }
        )
        .catch(err => {
          // TODO: handle creation errors.
          window.console.log(err);
        })
        .then(response => {
          this.dorefreshPrepared();
        });
    },
    reloadCharacter(response) {
      // fill the charInfo with response
      this.charInfo = response.data;
      if (this.charInfo.concentrating) {
        axios
          .get(
            `/api/spells/` +
              this.charInfo.concentrating,
            {
              headers: {
                Authorization: this.$root.$data.state.jwt,
                "Content-Type": "application/json"
              }
            }
          )
          .then(response => {
            this.concSpell = response.data;
          })
          .catch(err => {
            // TODO: handle creation errors.
            window.console.log(err);
          });
      } 
    },
    spellSaveDC(level, abilityScore) {
      return 8 + proficiencyBonus(level) +  abilityMod(abilityScore);
    },
    spellAttackMod(level, abilityScore) {
      return proficiencyBonus(level) + abilityMod(abilityScore);
    },
  }
};

function proficiencyBonus(level) {
  return 2 + Math.floor((level - 1)/4);
}

function abilityMod(abilityScore) {
  return Math.floor((abilityScore - 10)/2);
}

</script>

<style scoped>
.portrait-container {
  display: inline-block;
  position: relative;
  padding: 14px 4px;
  box-sizing: content-box;
  margin: 4px 12px;
}

.char-thumb {
  width: 184px;
  height: 184px !important;
}

.char-thumb-overlay {
  position: absolute;
  top: 0px;
  left: 0px;
}

.char-main {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
}

.char-details {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-grow: 1;
  justify-content: space-between;
  align-items:baseline;
  margin: 12px 0 12px 12px;
  font-family: 'Baskerville';
  font-size: 16px;
  min-height: 192px;
}

.char-details h2 {
  font-family: 'MrEavesSmallCaps';
  font-size: 36px;
  margin-bottom: 0;
}

.char-details h3 {
  font-family: 'MrEavesSmallCaps';
  font-size: 22px;
  margin-bottom: 0;
}

.no-break {
  white-space: nowrap;
}

.char-stat-container {
  font-family: 'ScalySansCaps';
  margin: 8px 0;
}

.char-stat-container div div {
  display: inline-block;
}

.char-stat-container div div:first-child {
  min-width: 140px;
}

.char-buttons {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  flex-wrap: wrap;
}

.char-buttons button {
  margin: 8px 8px 2px 0;
}

.spell-slots-container {
  padding: 8px 18px;
  font-family: 'ScalySansCaps';
}

h2 {
  font-family: 'MrEavesSmallCaps';
  font-size: 22px;
}

.top-spacer {
  margin-top: 24px;
}

.edit-dialog-container {
  padding: 0 12px;
}

.min-label {
  min-width: 56px !important;
}

.edit-thumb-container {
  position: relative;
}

.edit-thumb {
  width: 128px;
  height: 128px;
  margin: 4px 0;
  border-radius: 50%;
}

.edit-thumb-overlay {
  position: absolute;
  top: -4px;
  left: -8px;
}

.edit-left-space {
  margin-left: 18px;
}

.edit-buttons {
  display:flex;
  margin-top: 16px;
}

.edit-buttons button {
  margin: 0 0 0 8px;
}

.edit-buttons button:first-child {
  margin:0;
}

.conc-dialog {
  padding: 8px 12px;
  text-align: center;
}

.conc-dialog h2 {
  margin-bottom: 20px;
}

.conc-dialog-buttons {
  display: flex;
  flex-direction: row;
}
</style>

<style>
.lower-toast:first-child {
  margin-top: 50vh;
}

.center-toast-text {
  text-align: center;
}

.char-info {
  display: flex;
  flex-direction: column !important;
  justify-content: center;
  align-items: center;
}
</style>