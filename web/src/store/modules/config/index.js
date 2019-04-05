
import Vue from 'vue';
import * as _ from 'lodash';
import actions from './actions';
import constants from './constants';

export const mutationTypes = {
  SET_DISCOVERY_MODEL: 'SET_DISCOVERY_MODEL',
  DISCOVERY_MODEL_PROBLEMS: 'DISCOVERY_MODEL_PROBLEMS',
  SET_CONFIGURATION: 'SET_CONFIGURATION',
  SET_CONFIGURATION_SIGNING_PRIVATE: 'SET_CONFIGURATION_SIGNING_PRIVATE',
  SET_CONFIGURATION_SIGNING_PUBLIC: 'SET_CONFIGURATION_SIGNING_PUBLIC',
  SET_CONFIGURATION_TRANSPORT_PRIVATE: 'SET_CONFIGURATION_TRANSPORT_PRIVATE',
  SET_CONFIGURATION_TRANSPORT_PUBLIC: 'SET_CONFIGURATION_TRANSPORT_PUBLIC',
  SET_DISCOVERY_TEMPLATES: 'SET_DISCOVERY_TEMPLATES',
  SET_WIZARD_STEP: 'SET_WIZARD_STEP',
  SET_CLIENT_ID: 'SET_CLIENT_ID',
  SET_CLIENT_SECRET: 'SET_CLIENT_SECRET',
  SET_TOKEN_ENDPOINT: 'SET_TOKEN_ENDPOINT',
  SET_TOKEN_ENDPOINT_AUTH_METHOD: 'SET_TOKEN_ENDPOINT_AUTH_METHOD',
  SET_TOKEN_ENDPOINT_AUTH_METHODS: 'SET_TOKEN_ENDPOINT_AUTH_METHODS',
  SET_REQUEST_OBJECT_SIGNING_ALG_VALUES_SUPPORTED: 'SET_REQUEST_OBJECT_SIGNING_ALG_VALUES_SUPPORTED',
  SET_REQUEST_OBJECT_SIGNING_ALG: 'SET_REQUEST_OBJECT_SIGNING_ALG',
  SET_AUTHORIZATION_ENDPOINT: 'SET_AUTHORIZATION_ENDPOINT',
  SET_RESOURCE_BASE_URL: 'SET_RESOURCE_BASE_URL',
  SET_X_FAPI_FINANCIAL_ID: 'SET_X_FAPI_FINANCIAL_ID',
  SET_ISSUER: 'SET_ISSUER',
  SET_REDIRECT_URL: 'SET_REDIRECT_URL',
  SET_RESOURCE_ACCOUNT_ID: 'SET_RESOURCE_ACCOUNT_ID',
  SET_RESOURCE_STATEMENT_ID: 'SET_RESOURCE_STATEMENT_ID',
  SET_RESOURCE_ACCOUNT_IDS: 'SET_RESOURCE_ACCOUNT_IDS',
  SET_RESOURCE_STATEMENT_IDS: 'SET_RESOURCE_STATEMENT_IDS',
  ADD_RESOURCE_ACCOUNT_ID: 'ADD_RESOURCE_ACCOUNT_ID',
  REMOVE_RESOURCE_ACCOUNT_ID: 'REMOVE_RESOURCE_ACCOUNT_ID',
  ADD_RESOURCE_STATEMENT_ID: 'ADD_RESOURCE_STATEMENT_ID',
  REMOVE_RESOURCE_STATEMENT_ID: 'REMOVE_RESOURCE_STATEMENT_ID',
  SET_CREDITOR_ACCOUNT_NAME_SCHEME_NAME: 'SET_CREDITOR_ACCOUNT_NAME_SCHEME_NAME',
  SET_CREDITOR_ACCOUNT_IDENTIFICATION: 'SET_CREDITOR_ACCOUNT_IDENTIFICATION',
  SET_CREDITOR_ACCOUNT_NAME: 'SET_CREDITOR_ACCOUNT_NAME',
};

export const mutations = {
  [mutationTypes.SET_DISCOVERY_MODEL](state, discoveryModel) {
    Vue.set(state, 'discoveryModel', discoveryModel);
  },
  [mutationTypes.DISCOVERY_MODEL_PROBLEMS](state, problems) {
    state.problems = problems;
  },

  [mutationTypes.SET_CONFIGURATION](state, configuration) {
    state.configuration = configuration;
  },
  [mutationTypes.SET_CONFIGURATION_SIGNING_PRIVATE](state, signingPrivate) {
    state.configuration.signing_private = signingPrivate;
  },
  [mutationTypes.SET_CONFIGURATION_SIGNING_PUBLIC](state, signingPublic) {
    state.configuration.signing_public = signingPublic;
  },
  [mutationTypes.SET_CONFIGURATION_TRANSPORT_PRIVATE](state, transportPrivate) {
    state.configuration.transport_private = transportPrivate;
  },
  [mutationTypes.SET_CONFIGURATION_TRANSPORT_PUBLIC](state, transportPublic) {
    state.configuration.transport_public = transportPublic;
  },
  [mutationTypes.SET_DISCOVERY_TEMPLATES](state, templates) {
    state.discoveryTemplates = templates;
  },

  [mutationTypes.SET_WIZARD_STEP](state, step) {
    state.wizard.step = step;
  },

  [mutationTypes.SET_CLIENT_ID](state, value) {
    state.configuration.client_id = value;
  },
  [mutationTypes.SET_CLIENT_SECRET](state, value) {
    state.configuration.client_secret = value;
  },
  [mutationTypes.SET_TOKEN_ENDPOINT](state, value) {
    state.configuration.token_endpoint = value;
  },
  [mutationTypes.SET_TOKEN_ENDPOINT_AUTH_METHOD](state, value) {
    state.configuration.token_endpoint_auth_method = value;
  },
  [mutationTypes.SET_TOKEN_ENDPOINT_AUTH_METHODS](state, list) {
    state.token_endpoint_auth_methods = list;
  },
  [mutationTypes.SET_REQUEST_OBJECT_SIGNING_ALG_VALUES_SUPPORTED](state, list) {
    state.request_object_signing_alg_values_supported = list;
  },
  [mutationTypes.SET_REQUEST_OBJECT_SIGNING_ALG](state, value) {
    state.configuration.request_object_signing_alg = value;
  },
  [mutationTypes.SET_AUTHORIZATION_ENDPOINT](state, value) {
    state.configuration.authorization_endpoint = value;
  },
  [mutationTypes.SET_RESOURCE_BASE_URL](state, value) {
    state.configuration.resource_base_url = value;
  },
  [mutationTypes.SET_X_FAPI_FINANCIAL_ID](state, value) {
    state.configuration.x_fapi_financial_id = value;
  },
  [mutationTypes.SET_ISSUER](state, value) {
    state.configuration.issuer = value;
  },
  [mutationTypes.SET_RESOURCE_ACCOUNT_ID](state, { index, value }) {
    // Without the use of Vue.set the JSON editor tab view does not update on form input change.
    // https://vuejs.org/v2/api/#Vue-set
    const id = { account_id: value };
    Vue.set(state.configuration.resource_ids.account_ids, index, id);
  },
  [mutationTypes.SET_RESOURCE_STATEMENT_ID](state, { index, value }) {
    // Without the use of Vue.set the JSON editor tab view does not update on form input change.
    // https://vuejs.org/v2/api/#Vue-set
    const id = { statement_id: value };
    Vue.set(state.configuration.resource_ids.statement_ids, index, id);
  },
  [mutationTypes.SET_RESOURCE_ACCOUNT_IDS](state, value) {
    state.configuration.resource_ids.account_ids = value;
  },
  [mutationTypes.SET_RESOURCE_STATEMENT_IDS](state, value) {
    state.configuration.resource_ids.statement_ids = value;
  },
  [mutationTypes.ADD_RESOURCE_ACCOUNT_ID](state, value) {
    state.configuration.resource_ids.account_ids.push(value);
  },
  [mutationTypes.REMOVE_RESOURCE_ACCOUNT_ID](state, index) {
    state.configuration.resource_ids.account_ids.splice(index, 1);
  },
  [mutationTypes.ADD_RESOURCE_STATEMENT_ID](state, value) {
    state.configuration.resource_ids.statement_ids.push(value);
  },
  [mutationTypes.REMOVE_RESOURCE_STATEMENT_ID](state, index) {
    state.configuration.resource_ids.statement_ids.splice(index, 1);
  },

  [mutationTypes.SET_CREDITOR_ACCOUNT_NAME_SCHEME_NAME](state, value) {
    state.configuration.creditor_account.scheme_name = value;
  },
  [mutationTypes.SET_CREDITOR_ACCOUNT_IDENTIFICATION](state, value) {
    state.configuration.creditor_account.identification = value;
  },
  [mutationTypes.SET_CREDITOR_ACCOUNT_NAME](state, value) {
    state.configuration.creditor_account.name = value;
  },
};

// Converts problem key to discovery model JSON path.
const parseProblem = ({ key, error }) => {
  if (key && error) {
    const parts = key
      .replace('API', 'Api')
      .replace('URL', 'Url')
      .split('.')
      .map(w => _.lowerFirst(w));

    const path = parts.join('.');
    const parent = parts.slice(0, -1).join('.');

    return {
      path,
      parent,
      error,
    };
  }
  return {
    path: null,
    error,
  };
};

export const getters = {
  discoveryModel: state => state.discoveryModel,
  discoveryModelString: state => JSON.stringify(state.discoveryModel, null, 2),
  discoveryTemplates: state => state.discoveryTemplates,
  tokenAcquisition: state => (state.discoveryModel ? state.discoveryModel.discoveryModel.tokenAcquisition : null),
  problems: state => state.problems,
  discoveryProblems: state => (state.problems ? state.problems.map(p => parseProblem(p)) : null),
  configuration: state => state.configuration,
  configurationString: state => JSON.stringify(state.configuration, null, 2),
  resourceAccountIds: state => state.configuration.resource_ids.account_ids,
  resourceStatementIds: state => state.configuration.resource_ids.statement_ids,
  /**
   * Computes what the user can navigate to based on the current step they are on.
   */
  navigation: (state) => {
    const { step } = state.wizard;
    const navigation = {
      '/wizard/continue-or-start': step > 0,
      '/wizard/import/review': step > 0,
      '/wizard/import/rerun': step > 0,
      '/wizard/discovery-config': step > constants.WIZARD.STEP_ONE,
      '/wizard/configuration': step > constants.WIZARD.STEP_TWO,
      '/wizard/overview-run': step > constants.WIZARD.STEP_THREE,
      '/wizard/export': step > constants.WIZARD.STEP_FOUR,
    };
    return navigation;
  },
};

export const state = {
  discoveryTemplates: [],
  discoveryModel: null,
  problems: null,
  token_endpoint_auth_methods: [],
  request_object_signing_alg_values_supported: [],

  configuration: {
    signing_private: '',
    signing_public: '',
    transport_private: '',
    transport_public: '',
    client_id: '',
    client_secret: '',
    token_endpoint: '',
    token_endpoint_auth_method: 'client_secret_basic',
    request_object_signing_alg: '',
    authorization_endpoint: '',
    resource_base_url: '',
    x_fapi_financial_id: '',
    issuer: '',
    redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
    resource_ids: {
      account_ids: [{ account_id: '' }],
      statement_ids: [{ statement_id: '' }],
    },
    creditor_account: {
      scheme_name: '',
      identification: '',
      name: '',
    },
  },

  wizard: {
    step: constants.WIZARD.STEP_ONE,
  },
};

export default {
  namespaced: true,
  state,
  actions,
  getters,
  mutations,
  mutationTypes,
};
