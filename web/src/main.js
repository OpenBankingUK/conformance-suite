import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue/dist/bootstrap-vue.esm.min';

/* global fetch */
import 'whatwg-fetch';

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-vue/dist/bootstrap-vue.min.css';

import App from './App.vue';
import router from './router';
import store from './store/';

Vue.use(BootstrapVue);

// Don't warn about using the dev version of Vue in development.
Vue.config.productionTip = process.env.NODE_ENV === 'production';

// Use webpack require.context to import templates and images.
// See: https://webpack.js.org/guides/dependency-management/#require-context
//      https://vuejs.org/v2/guide/components-registration.html#Automatic-Global-Registration-of-Base-Components
const requireTemplates = require.context('../../pkg/discovery/templates/', false, /.+\.json$/);
const discoveryTemplates = requireTemplates.keys().map(file => requireTemplates(file));

const requireImages = require.context('./assets/images/', false, /.+\.png$/);
const discoveryImages = {};
requireImages.keys().forEach(file => discoveryImages[file] = requireImages(file)); // eslint-disable-line

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');

// Store templates and images.
store.dispatch('config/setDiscoveryTemplates', { discoveryTemplates, discoveryImages }, { root: true });
