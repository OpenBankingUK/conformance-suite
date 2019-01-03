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

  [types.SET_CONFIGURATION_SIGNING_PRIVATE](state, signingPrivate) {
    state.configuration.signing_private = signingPrivate;
  },
  [types.SET_CONFIGURATION_SIGNING_PUBLIC](state, signingPublic) {
    state.configuration.signing_public = signingPublic;
  },
  [types.SET_CONFIGURATION_TRANSPORT_PRIVATE](state, transportPrivate) {
    state.configuration.transport_private = transportPrivate;
  },
  [types.SET_CONFIGURATION_TRANSPORT_PUBLIC](state, transportPublic) {
    state.configuration.transport_public = transportPublic;
  },
  [types.SET_CONFIGURATION_ERRORS](state, errors) {
    state.errors.configuration = errors;
  },
  [types.ADD_CONFIGURATION_ERRORS](state, errors) {
    state.errors.configuration.push(errors);
  },
  [types.CLEAR_CONFIGURATION_ERRORS](state) {
    state.errors.configuration = [];
  },

  [types.SET_TEST_CASES](state, testCases) {
    state.testCases = testCases;
  },
  [types.SET_TEST_CASES_ERROR](state, errors) {
    state.errors.testCases = errors;
  },
};
