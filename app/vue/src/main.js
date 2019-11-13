import Vue from 'vue'
import BootstrapVue from 'bootstrap-vue'
import App from './App.vue'
import * as firebase from "firebase/app"
import "firebase/auth";
import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import './assets/style/style.css'

// OMITTED FOR PRIVACY
const firebaseConfig = {
  apiKey: "", 
  authDomain: "", 
  databaseURL: "",
  projectId: "",
  storageBucket: "",
  messagingSenderId: "",
  appId: ""
};

firebase.initializeApp(firebaseConfig);

firebase.auth().onAuthStateChanged(function(user) {
  if (user) {
    store.state.user = user;
    firebase.auth().currentUser.getIdToken(true).then(function(idToken) {
      store.state.jwt = idToken;
    }).catch(function(error) {
      console.error("couldn't get jwt: "+ error)
    });
  } else {
    store.state.user = null;
    store.state.jwt = null;
    store.state.characterID = null;
  }
});

// `store` is the global state store for our app.
// Add properties to the `state` object to make them accessible across components.
let store = {
  debug: true,
  state: {
    message: 'I\'m global!',
    user: Object,
    jwt: null,
    skipToCharPage: false,
    characters: null,
    characterID: null,
  },
  setMessageAction (newValue) {
    if (this.debug) console.log('setMessageAction triggered with', newValue)
    this.state.message = newValue
  },
  clearMessageAction () {
    if (this.debug) console.log('clearMessageAction triggered')
    this.state.message = ''
  }
}

/**
 * Check for #character, used to bypass login and load a default character
 * while character select screen is under development.
 */

if (location.hash == "#character") {
  store.state.skipToCharPage = true
  store.state.character = {
    name: null,
    concentrating: {
      name: null,
    },
    spells: [
      {
        name: null,
        description: null,
      },
    ]
  };
}

Vue.config.productionTip = false

Vue.use(BootstrapVue);

new Vue({
  data: {
    state: store.state,
  },
  render: h => h(App),
}).$mount('#app')
