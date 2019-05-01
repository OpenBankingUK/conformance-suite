import Vue from 'vue';
import Vuex from 'vuex';
import createLogger from 'vuex/dist/logger';

import loadDiscoveryTemplates from './modules/config/loadDiscoveryTemplates';

import config from './modules/config';
import testcases from './modules/testcases';
import status from './modules/status';
import exporter from './modules/exporter';
import importer from './modules/importer';

Vue.use(Vuex);

const strict = process.env.NODE_ENV !== 'production';
const plugins = process.env.NODE_ENV !== 'production' && process.env.NODE_ENV !== 'test' ? [createLogger()] : [];
const store = new Vuex.Store({
  modules: {
    config,
    testcases,
    status,
    exporter,
    importer,
  },
  strict,
  plugins,
});

const { discoveryTemplates, discoveryImages } = loadDiscoveryTemplates();
// Store templates and images.
store.dispatch('config/setDiscoveryTemplates', { discoveryTemplates, discoveryImages }, { root: true });

// to debug the store install the Vue.js chrome/firefox extension
export default store;
