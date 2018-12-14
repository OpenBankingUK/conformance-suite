import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue';
// import axios from 'axios';
// import VueAxios from 'vue-axios';

/* global fetch */
import 'whatwg-fetch';

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

import App from './App.vue';
import router from './router';
import store from './store/';
import './registerServiceWorker';

import './assets/css/app.css';

Vue.use(BootstrapVue);
// Vue.use(VueAxios, axios);
// Vue.use(Antd);

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');
