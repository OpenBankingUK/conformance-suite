import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue/dist/bootstrap-vue.esm.min';

/* global fetch */
import 'whatwg-fetch';

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-vue/dist/bootstrap-vue.min.css';

import App from './App.vue';
import router from './router';
import store from './store';

Vue.use(BootstrapVue);

// Don't warn about using the dev version of Vue in development.
Vue.config.productionTip = process.env.NODE_ENV === 'production';
// Debug when not running in production
Vue.config.debug = process.env.NODE_ENV !== 'production';

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');
