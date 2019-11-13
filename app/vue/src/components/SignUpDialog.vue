<template>
  <div>
    <button
      class="plaintextButton"
      @click="openSignUpDialog"
    >
      Sign up
    </button>
    <b-modal
      id="signUp"
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
            <h2>Create an account</h2>
            <button
              class="phbIconButton"
              title="cancel"
              @click="cancel"
            >
              <img src="@/assets/icons/clear.svg">
            </button>
          </div>
          <div class="inputs">
            <div class="phbTextInputContainer">
              <div class="phbInputLabel">
                Email:
              </div>
              <input
                v-model="email"
                class="phbTextInput"
                @input="validateEmail"
              >
              <div class="validationMark">
                <img
                  v-if="emailOk"
                  title="email ok"
                  src="@/assets/icons/check.svg"
                >
              </div>
            </div>
            <div class="phbTextInputContainer">
              <div class="phbInputLabel">
                Password:
              </div>
              <input
                v-model="password"
                type="password"
                class="phbTextInput"
                @input="validatePassword"
              >
              <div class="validationMark">
                <img
                  v-if="passwordOk"
                  title="password ok"
                  src="@/assets/icons/check.svg"
                >
              </div>
            </div>
            <div class="phbTextInputContainer">
              <div class="phbInputLabel">
                Confirm:
              </div>
              <input
                v-model="confirmPassword"
                type="password"
                class="phbTextInput"
                @input="validateConfirmation"
              >
              <div class="validationMark">
                <img
                  v-if="confirmOk"
                  title="password confirmed"
                  src="@/assets/icons/check.svg"
                >
              </div>
            </div>
          </div>
          <div class="signUpButton">
            <button
              :disabled="!emailOk || !confirmOk"
              class="phbButton"
              @click="signUp"
            >
              sign up
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
import * as firebase from "firebase/app";
import "firebase/auth";

// eslint-disable-next-line
const emailRegex = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

function baseState() {
  return {
    email: '',
    emailOk: false,
    password: '',
    passwordOk: false,
    confirmPassword: '',
    confirmOk: false,
  };
}

export default {
  name: 'SignUpDialog',
  data() {
      return baseState();
  },
  methods: {
    signUp: function() {
      firebase.auth()
        .createUserWithEmailAndPassword(this.email, this.password)
        .catch(err => {
          this.$bvToast.toast(`Couldn't create account: ${err.message}`,
            {
              toaster: 'b-toaster-bottom-center',
              variant: 'warning',
              noCloseButton: true,
            });
          window.console.error(err);
        });
      this.$bvModal.hide('signUp');
    },
    openSignUpDialog: function() {
      Object.assign(this.$data, baseState());
      this.$bvModal.show('signUp');
    },
    validateEmail: function(event) {
      if (event.target.value.match(emailRegex) != null) {
        if (!this.emailOk) {
          this.emailOk = true;
        }
      } else {
        if (this.emailOk) {
          this.emailOk = false;
        }
      }
    },
    validatePassword: function(event) {
      if (event.target.value.length > 5) {
        this.passwordOk = true;
      } else {
        if (this.passwordOk) {
          this.passwordOk = false;
        }
      }
      if (event.target.value == this.confirmPassword) {
        if (!this.confirmOk) {
          this.confirmOk = true;
        }
      } else {
        if (this.confirmOk) {
          this.confirmOk = false;
        }
      }
    },
    validateConfirmation: function(event) {
      if (this.passwordOk && event.target.value == this.password) {
        this.confirmOk = true;
      } else {
        if (this.confirmOk) {
          this.confirmOk = false;
        }
      }
    },
    cancel: function() {
      this.$bvModal.hide('signUp');
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.signUpButton {
  display: flex;
  justify-content: flex-end;
}

.signUpButton button {
  margin-top: 8px;
}

h2 {
  font-family: 'MrEavesSmallCaps';
  font-size: 24px !important;
}

.validationMark {
  min-width: 24px;
  margin: 1px 6px;
}

.inputs {
  margin: 24px 0;
}
</style>
