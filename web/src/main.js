import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue';
import axios from 'axios';
import VueAxios from 'vue-axios';

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

// https://vuecomponent.github.io/ant-design-vue/docs/vue/introduce/
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';

import App from './App.vue';
import router from './router';
import store from './store/';
import './registerServiceWorker';

import Default from './layouts/Default.vue';
import Clean from './layouts/Clean.vue';
import './assets/css/app.css';

Vue.component('default-layout', Default);
Vue.component('clean-layout', Clean);

Vue.use(BootstrapVue);
Vue.use(VueAxios, axios);
Vue.use(Antd);

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');
