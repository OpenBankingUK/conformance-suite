import Vue from 'vue';
import axios from 'axios';
import VueAxios from 'vue-axios';
// https://vuecomponent.github.io/ant-design-vue/docs/vue/introduce/
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';

import App from './App';
import router from './router';
import store from './store';
import socketPlugin from './plugins/socket';
import Default from './layouts/Default';
import Clean from './layouts/Clean';
import '../css/app.css';

Vue.component('default-layout', Default);
Vue.component('clean-layout', Clean);

Vue.use(VueAxios, axios);
Vue.use(Antd);
Vue.use(socketPlugin);
Vue.config.productionTip = false;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  components: { App },
  router,
  store,
  template: '<App/>',
});
