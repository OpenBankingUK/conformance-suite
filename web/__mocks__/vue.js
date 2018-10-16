import Vue from 'vue';

// Mock vue-axios for http calls
Vue.axios = {
  post: jest.fn(),
  get: jest.fn(),
  defaults: { headers: { common: {} } },
};

export default Vue;
