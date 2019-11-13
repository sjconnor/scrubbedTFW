<template>
  <div class="log-in-content">
    <div class="log-in-container">
      <h2>Welcome to Tinfoil Wizard</h2>
      <div class="log-in-flow">
        <div class="phbTextInputContainer">
          <div class="phbInputLabel">
            Email:
          </div>
          <input
            v-model="email"
            class="phbTextInput"
            tabindex="1"
          >
        </div>
        <div class="phbTextInputContainer">
          <div class="phbInputLabel">
            Password:
          </div>
          <input
            v-model="password"
            name="password"
            class="phbTextInput"
            type="password"
            @keyup.enter="logIn"
            tabindex="2"
          >
        </div>  
        <div class="log-in-button-row">
          <button
            class="phbButton"
            @click="logIn"
            tabindex="3"
          >
            log in
          </button>
        </div>
      </div>
    </div>
    <div class="alt-button-row">
      <div class="alt-button-container">
        <SignUpDialog />
      </div>
      |
      <div class="alt-button-container">
        <RecoverPasswordDialog />
      </div>
    </div>
  </div>
</template>

<script>
import * as firebase from "firebase/app";
import "firebase/auth";
import SignUpDialog from './SignUpDialog.vue';
import RecoverPasswordDialog from './RecoverPasswordDialog.vue'

// eslint-disable-next-line
const emailRegex = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

export default {
  name: 'Login',
  components: {
    SignUpDialog,
    RecoverPasswordDialog,
  },
  data() {
      return {
          email: '',
          password: '',
      }
  },
  methods: {
    logIn: function() {
      const email = this.email;
      const password = this.password;
      if (email.match(emailRegex) == null || password == '') {
        this.$bvToast.toast(`Log in failed: Enter a valid email and password`,
        {
          toaster: 'b-toaster-bottom-center',
          variant: 'warning',
          noCloseButton: true,
        });
      } else {
        firebase.auth().setPersistence(firebase.auth.Auth.Persistence.LOCAL)
          .then(function() {
            return firebase.auth().signInWithEmailAndPassword(email, password);
          })
          .catch(err => {
            this.$bvToast.toast(`Log in failed: ${err.message}`,
              {
                toaster: 'b-toaster-bottom-center',
                variant: 'warning',
                noCloseButton: true,
              });
            window.console.error(err);
          });
      }
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.log-in-container {
  border: solid;
  border-image: url('../assets/icons/frameBorder.svg') repeat;
  border-image-slice: 60 fill;
  border-image-width: 60px;
  padding: 24px;
}

.log-in-content {
  align-items: center;
  display: flex;
  justify-content: center;
  flex-direction: column;
  font-family: 'ScalySansCaps';
  font-size: 18px;
  height: calc(100vh - 200px);
}

.log-in-content h2 {
  font-family: MrEavesSmallCaps;
  font-size: 32px;
  text-align: center;
}

.log-in-flow {
  margin: 32px 0 0 0;
}

.log-in-button-row {
  display: flex;
  justify-content: center;
  margin-top: 28px;
  margin-bottom: 4px;
}

.alt-button-row {
  display: flex;
  flex-direction: row;
  margin-top: 24px;
}

.alt-button-container {
  display: flex;
  min-width: 200px;
  margin: 0 8px;
  justify-content: center;
}
</style>
