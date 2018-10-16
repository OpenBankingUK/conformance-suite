import Vue from 'vue';
import Vuex from 'vuex';

import validations from './modules/validations';
import reporter from './modules/reporter';
import config from './modules/config';

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';

// to debug the store install the Vue.js chrome/firefox extension
export default new Vuex.Store({
  modules: {
    validations,
    reporter,
    config,
  },
  strict: debug,
});
