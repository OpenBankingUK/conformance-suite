import Vue from 'vue';

import Vuex from 'vuex';
import createPersistedState from 'vuex-persistedstate';
import validations from './modules/validations';
import reporter from './modules/reporter';
import user from './modules/user';
import config from './modules/config';

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';
const plugins = [createPersistedState({ paths: ['user.profile', 'user.signedIn'] })];

// to debug the store install the Vue.js chrome/firefox extension
export default new Vuex.Store({
  modules: {
    validations,
    reporter,
    user,
    config,
  },
  strict: debug,
  plugins,
});
