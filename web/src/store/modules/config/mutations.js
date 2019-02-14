import Vue from 'vue';
import * as types from './mutation-types';

export default {
  [types.SET_DISCOVERY_MODEL](state, discoveryModel) {
    Vue.set(state, 'discoveryModel', discoveryModel);
  },
  [types.DISCOVERY_MODEL_PROBLEMS](state, problems) {
    state.problems = problems;
  },

  [types.SET_CONFIGURATION](state, configuration) {
    state.configuration = configuration;
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
  [types.SET_DISCOVERY_TEMPLATES](state, templates) {
    state.discoveryTemplates = templates;
  },

  [types.SET_WIZARD_STEP](state, step) {
    state.wizard.step = step;
  },

  [types.SET_CLIENT_ID](state, value) {
    state.configuration.client_id = value;
  },
  [types.SET_CLIENT_SECRET](state, value) {
    state.configuration.client_secret = value;
  },
  [types.SET_TOKEN_ENDPOINT](state, value) {
    state.configuration.token_endpoint = value;
  },
  [types.SET_AUTHORIZATION_ENDPOINT](state, value) {
    state.configuration.authorization_endpoint = value;
  },
  [types.SET_X_FAPI_FINANCIAL_ID](state, value) {
    state.configuration.x_fapi_financial_id = value;
  },
  [types.SET_REDIRECT_URL](state, value) {
    state.configuration.redirect_url = value;
  },
};
