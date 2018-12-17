import Vue from 'vue';
import * as types from './mutation-types';

export default {
  [types.SET_CONFIG](state, config) {
    state.main = config;
  },
  [types.SET_DISCOVERY_MODEL](state, discoveryModel) {
    Vue.set(state, 'discoveryModel', discoveryModel);
  },
  [types.DISCOVERY_MODEL_PROBLEMS](state, problems) {
    state.problems = problems;
  },
};
