<template>
  <div
    v-if="state.user != null"
    class="profile-container"
  >
    <button
      class="icon-button"
      title="profile"
      @click="openProfile"
    >
      <img src="@/assets/icons/profile.svg">
    </button>
    <b-modal
      id="profile"
      centered
      hide-header
      body-class="phbDialogDialogOverride"
      content-class="phbDialogContentOverride"
      dialog-class="profile-dialog-override"
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
            <h2 class="title">
              Profile
            </h2>
            <button
              class="icon-button"
              title="logout"
              @click="cancel"
            >
              <img src="@/assets/icons/clear.svg">
            </button>
          </div>
          <div class="profile-inputs">
            <div class="phbTextInputContainer">
              <div class="phbInputLabel">
                Email:
              </div>
              <div
                v-if="!editEmail"
                class="profile-ph-container"
              >
                <div class="profile-input-ph">
                  {{ state.user.email }}
                </div>
                <button
                  class="icon-button"
                  title="edit email"
                  @click="toggleEditEmail"
                >
                  <img src="@/assets/icons/edit.svg">
                </button>
              </div>
              <div
                v-else
                class="profile-ph-container"
              >
                <input
                  id="editEmailInput"
                  v-model="email"
                  class="phbTextInput"
                >
                <button
                  class="icon-button"
                  title="cancel"
                  @click="cancelEditEmail"
                >
                  <img src="@/assets/icons/clear.svg">
                </button>
                <button
                  class="icon-button"
                  title="submit"
                  @click="saveEmail"
                >
                  <img src="@/assets/icons/check.svg">
                </button>
              </div>
            </div>
            <div class="phbTextInputContainer">
              <div class="phbInputLabel">
                Password:
              </div>
              <div
                v-if="!editPassword"
                class="profile-ph-container"
              >
                <div class="profile-input-ph">
                  •••••••••••
                </div>
                <button
                  class="icon-button"
                  title="edit password"
                  @click="toggleEditPassword"
                >
                  <img src="@/assets/icons/edit.svg">
                </button>
              </div>
              <div
                v-else
                class="profile-ph-container"
              >
                <input
                  id="editPasswordInput"
                  v-model="password"
                  type="password"
                  class="phbTextInput"
                  :placeholder="passwordPh"
                  @focus="clearPh"
                  @blur="setPh"
                >
                <button
                  class="icon-button"
                  title="cancel"
                  @click="cancelEditPw"
                >
                  <img src="@/assets/icons/clear.svg">
                </button>
                <button
                  class="icon-button"
                  title="submit"
                  @click="savePassword"
                >
                  <img src="@/assets/icons/check.svg">
                </button>
              </div>
            </div>
          </div>
          <div class="profile-buttons">
            <button
              class="phbButton"
              @click="logout"
            >
              log out
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
//import axios from "axios"
import * as firebase from "firebase/app"
import "firebase/auth"

const editSvg = "@/assets/icons/edit.svg";
const clearSvg = "@/assets/icons/clear.svg";
const editEmailTitle = "Edit email";
const clearTitle = "Clear";
const editPasswordTitle = "Edit password";
const defaultPh = '•••••••••••';

export default {
  name: 'Profile',
  data() {
    return {
      state: this.$root.$data.state,
      email: '',
      password: '',
      editEmail: false,
      editPassword: false,
      toggleEmailTitle: '',
      toggleEmailSvg: '',
      togglePasswordTitle: '',
      togglePasswordSvg: '',
      passwordPh: '',
    }
  },
  methods: {
    cancel: function() {
      this.$bvModal.hide('profile');
    },
    openProfile: function() {
      this.email = this.state.user.email;
      this.password = '';
      this.editEmail = false;
      this.editPassword = false;
      this.toggleEmailTitle = editEmailTitle;
      this.toggleEmailSvg = editSvg;
      this.togglePasswordTitle = editPasswordTitle;
      this.togglePasswordSvg = editSvg;
      this.passwordPh = defaultPh;
      this.$bvModal.show('profile');
    },
    saveEmail: function() {
      if (this.email != this.state.user.email) {
        this.state.user.updateEmail(this.email)
        .then(() => {
          this.$bvToast.toast(`Email updated`,
            {
              toaster: 'b-toaster-bottom-center',
              variant: 'success',
              noCloseButton: true,
            });
        })
        .catch(err => {
          this.$bvToast.toast(`Error updating email: ${err.message}`,
            {
              toaster: 'b-toaster-bottom-center',
              variant: 'warning',
              noCloseButton: true,
            });
          window.console.error(err);
        });
      }
      this.editEmail = false;
    },
    savePassword: function() {
      if (this.password != '') {
        this.state.user.updatePassword(this.password)
        .then(() => {
          this.$bvToast.toast(`Password updated`,
            {
              toaster: 'b-toaster-bottom-center',
              variant: 'success',
              noCloseButton: true,
            });
        })
        .catch(err => {
          this.$bvToast.toast(`Error updating password: ${err.message}`,
            {
              toaster: 'b-toaster-bottom-center',
              variant: 'warning',
              noCloseButton: true,
            });
          window.console.error(err);
        });
      }
      this.editPassword = false;
    },
    logout: function() {
      firebase.auth().signOut();
      this.$bvModal.hide('profile');
    },
    toggleEditEmail: function() {
      this.editEmail = !this.editEmail;
      if (this.editEmail) {
        this.toggleEmailTitle = clearTitle;
        this.toggleEmailSvg = clearSvg;
        setTimeout(() => {document.getElementById('editEmailInput').focus();}, 10);
      } else {
        this.toggleEmailTitle = editEmailTitle;
        this.toggleEmailSvg = editSvg;
        this.email = this.state.user.email;
      }
    },
    toggleEditPassword: function() {
      this.editPassword = !this.editPassword;
      if (this.editPassword) {
        this.togglePasswordTitle = clearTitle;
        this.togglePasswordSvg = clearSvg;
        setTimeout(() => {document.getElementById('editPasswordInput').focus();}, 10);
      } else {
        this.togglePasswordTitle = editPasswordTitle;
        this.togglePasswordSvg = editSvg;
        this.password = '';
      }
    },
    clearPh: function() {
      if (this.passwordPh != '') {
        this.passwordPh = '';
      }
    },
    setPh: function(event) {
      if (this.passwordPh != defaultPh && event.target.value == '') {
        this.passwordPh = defaultPh;
      }
    },
    cancelEditPw: function() {
      this.editPassword = false;
    },
    cancelEditEmail: function() {
      this.editEmail = false;
    }
  },
}
</script>

<style>
.profile-dialog-override {
  max-width: 540px !important;
}
</style>

<style scoped>
h2 {
  font-family: 'MrEavesSmallCaps';
  font-size: 24px !important;
}

.profile-container {
  display: inline-block;
}

.profile-inputs {
  margin: 24px 0;
}

.profile-buttons {
  display: flex;
  justify-content: center;
}

.profile-buttons button {
  margin-left: 8px;
}

.profile-input-ph {
  display: inline-block;
  padding-left: 2px;
  width: 100%;
  min-height: 29px;
  text-overflow: ellipsis;
}

.profile-ph-container {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  width: 100%;
}
</style>
