<template>
  <div @click.stop>
    <button
      class="icon-button hover-hide"
      title="delete character"
      @click="openDeleteCharDialog"
    >
      <img src="../assets/icons/clear.svg">
    </button>
    <b-modal
      :id="dialogId"
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
            <h2>Delete character</h2>
            <button
              class="phbIconButton"
              title="cancel"
              @click="cancel"
            >
              <img src="../assets/icons/clear.svg">
            </button>
          </div>
          <p>Are you sure?</p>
          <div class="action-buttons">
            <button
              class="phbButton"
              @click="cancel"
            >
              cancel
            </button>
            <button
              class="phbButton"
              @click="send"
            >
              confirm
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

export default {
  name: 'DeleteCharDialog',
  props: {
    charIndex: {
      type: Number,
      default: -1
    }
  },
  data() {
      return {
        dialogId: 'deleteChar'+this.charIndex,
        state: this.$root.$data.state,
      };
  },
  methods: {
    openDeleteCharDialog: function(event) {
      event.stopPropagation();
      this.$bvModal.show(this.dialogId);
    },
    cancel: function() {
      this.$bvModal.hide(this.dialogId);
    },
    send: function() {
      let chars = this.state.characters;
      let charBackup = chars.splice(this.charIndex, 1)[0];
      axios.delete(`/api/characters/${charBackup.id}`, {
        headers: {
          Authorization: this.state.jwt,
        },
      })
      .then(() => {
        this.$bvToast.toast(`Character deleted successfully`,
          {
            toaster: 'b-toaster-bottom-center',
            variant: 'success',
            noCloseButton: true,
          });
      })
      .catch(err => {
          this.state.characters.splice(this.charIndex, 0, charBackup);
          this.$bvToast.toast(`Error deleting character: ${err.message}`,
          {
            toaster: 'b-toaster-bottom-center',
            variant: 'warning',
            noCloseButton: true,
          });
        window.console.error(err);
      });
      this.$bvModal.hide(this.dialogId);
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.icon-button {
  background: none;
  border: none;
  height: 24px;
  padding: 0;
  margin-left: 4px;
  z-index: 99;
}

.icon-button img {
  height: 24px;
  width: 24px;
}

.icon-button:focus {
  outline: none;
}


.action-buttons {
  display: flex;
  justify-content: flex-end;
}

.action-buttons button {
  margin-top: 8px;
  margin-left: 12px;
}

h2 {
  font-family: 'MrEavesSmallCaps';
  font-size: 24px !important;
}

p {
  margin: 12px 0;
}
</style>
