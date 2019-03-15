import * as _ from 'lodash';
import * as types from './mutation-types';
import constants from './constants';

// import discovery from '../../../api/discovery';
import api from '../../../api';

const findImageData = (model, images) => {
  const { name } = model.discoveryModel;
  const customImage = `./${name}.png`;
  return images[customImage] || images['./no-image-discovery-icon.png'];
};

export default {
  setDiscoveryTemplates({ commit }, { discoveryTemplates, discoveryImages }) {
    const templates = discoveryTemplates.map(template => ({
      model: template,
      image: findImageData(template, discoveryImages),
    }));
    commit(types.SET_DISCOVERY_TEMPLATES, templates);
  },
  setDiscoveryModel({ commit, dispatch, state }, editorString) {
    const value = JSON.stringify(state.discoveryModel);
    if (_.isEqual(value, editorString)) {
      return;
    }

    try {
      const discoveryModel = JSON.parse(editorString);
      commit(types.SET_DISCOVERY_MODEL, discoveryModel);
      commit(types.DISCOVERY_MODEL_PROBLEMS, null);
      dispatch('status/clearErrors', null, { root: true });
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_TWO);
    } catch (e) {
      const problems = [{
        key: null,
        error: e.message,
      }];
      commit(types.DISCOVERY_MODEL_PROBLEMS, problems);
      dispatch('status/setErrors', [e.message], { root: true });
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_TWO);
    }
  },
  /**
   * Step 2: validate the Discovery Config.
   * Route: `/wizard/discovery-config`.
   */
  async validateDiscoveryConfig({ commit, dispatch, state }) {
    try {
      const setShowLoading = flag => dispatch('status/setShowLoading', flag, { root: true });
      const { success, problems, response } = await api.validateDiscoveryConfig(state.discoveryModel, setShowLoading);
      if (success) {
        commit(types.DISCOVERY_MODEL_PROBLEMS, null);
        const tokenEndpoint = _.first(_.values(response.token_endpoints));
        commit(types.SET_TOKEN_ENDPOINT, tokenEndpoint);

        const defaultAuthMethod = _.first(_.values(response.default_token_endpoint_auth_method));
        commit(types.SET_TOKEN_ENDPOINT_AUTH_METHOD, defaultAuthMethod);

        const authMethods = _.first(_.values(response.token_endpoint_auth_methods));
        commit(types.SET_TOKEN_ENDPOINT_AUTH_METHODS, authMethods);

        const authorizationEndpoint = _.first(_.values(response.authorization_endpoints));
        commit(types.SET_AUTHORIZATION_ENDPOINT, authorizationEndpoint);

        const issuer = _.first(_.values(response.issuers));
        commit(types.SET_ISSUER, issuer);

        dispatch('status/clearErrors', null, { root: true });
        commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);
      } else {
        commit(types.DISCOVERY_MODEL_PROBLEMS, problems);
        dispatch('status/setErrors', problems.map(p => p.error), { root: true });
        commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_TWO);
      }
    } catch (e) {
      commit(types.DISCOVERY_MODEL_PROBLEMS, [{
        key: null,
        error: e.message,
      }]);
      dispatch('status/setErrors', [e.message], { root: true });
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_TWO);
    }
    return null;
  },

  setConfigurationJSON({ commit, dispatch, state }, editorString) {
    const value = JSON.stringify(state.configuration);
    if (_.isEqual(value, editorString)) {
      return;
    }

    try {
      const config = JSON.parse(editorString);
      const merged = _.merge(_.clone(state.configuration), config);
      const validKeys = [
        'signing_private',
        'signing_public',
        'transport_private',
        'transport_public',
        'client_id',
        'client_secret',
        'token_endpoint',
        'token_endpoint_auth_method',
        'authorization_endpoint',
        'resource_base_url',
        'x_fapi_financial_id',
        'issuer',
        'redirect_url',
        'resource_ids',
      ];
      const newConfig = _.pick(merged, validKeys);
      commit(types.SET_CONFIGURATION, newConfig);
      dispatch('status/clearErrors', null, { root: true });
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);
    } catch (e) {
      dispatch('status/setErrors', [e.message], { root: true });
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);
    }
  },
  setConfigurationSigningPrivate({ commit, state }, signingPrivate) {
    if (_.isEqual(state.configuration.signing_private, signingPrivate)) {
      return;
    }

    commit(types.SET_CONFIGURATION_SIGNING_PRIVATE, signingPrivate);
    commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);
  },
  setConfigurationSigningPublic({ commit, state }, signingPublic) {
    if (_.isEqual(state.configuration.signing_public, signingPublic)) {
      return;
    }

    commit(types.SET_CONFIGURATION_SIGNING_PUBLIC, signingPublic);
    commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);
  },
  setConfigurationTransportPrivate({ commit, state }, transportPrivate) {
    if (_.isEqual(state.configuration.transport_private, transportPrivate)) {
      return;
    }

    commit(types.SET_CONFIGURATION_TRANSPORT_PRIVATE, transportPrivate);
    commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);
  },
  setConfigurationTransportPublic({ commit, state }, transportPublic) {
    if (_.isEqual(state.configuration.transport_public, transportPublic)) {
      return;
    }

    commit(types.SET_CONFIGURATION_TRANSPORT_PUBLIC, transportPublic);
    commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);
  },
  addResourceAccountID({ commit }, accountID) {
    commit(types.ADD_RESOURCE_ACCOUNT_ID, accountID);
  },
  removeResourceAccountID({ commit, state }, index) {
    if (index < 0 || index > state.configuration.resource_ids.account_ids - 1) {
      return;
    }

    commit(types.REMOVE_RESOURCE_ACCOUNT_ID, index);
  },
  addResourceStatementID({ commit }, statementID) {
    commit(types.ADD_RESOURCE_STATEMENT_ID, statementID);
  },
  removeResourceStatementID({ commit, state }, index) {
    if (index < 0 || index > state.configuration.resource_ids.statement_ids - 1) {
      return;
    }

    commit(types.REMOVE_RESOURCE_STATEMENT_ID, index);
  },
  /**
   * Step 3: Validate the configuration.
   * Route: `/wizard/configuration`.
   */
  async validateConfiguration({ commit, dispatch, state }) {
    dispatch('status/clearErrors', null, { root: true });

    const errors = [];
    if (_.isEmpty(state.configuration.signing_private)) {
      errors.push('Signing Private Certificate (.key) empty');
    }
    if (_.isEmpty(state.configuration.signing_public)) {
      errors.push('Signing Public Certificate (.pem) empty');
    }
    if (_.isEmpty(state.configuration.transport_private)) {
      errors.push('Transport Private Certificate (.key) empty');
    }
    if (_.isEmpty(state.configuration.transport_public)) {
      errors.push('Transport Public Certificate (.pem) empty');
    }
    if (_.isEmpty(state.configuration.resource_ids.account_ids)) {
      errors.push('Account IDs empty');
    }
    if (_.isEmpty(state.configuration.resource_ids.statement_ids)) {
      errors.push('Statement IDs empty');
    }

    if (_.isEmpty(state.configuration.client_id)) {
      errors.push('Client ID empty');
    }
    if (_.isEmpty(state.configuration.client_secret)) {
      errors.push('Client Secret empty');
    }
    if (_.isEmpty(state.configuration.token_endpoint)) {
      errors.push('Token Endpoint empty');
    }
    if (_.isEmpty(state.configuration.token_endpoint_auth_method)) {
      errors.push('Token Endpoint Auth Method empty');
    }
    if (_.isEmpty(state.configuration.authorization_endpoint)) {
      errors.push('Authorization Endpoint empty');
    }
    if (_.isEmpty(state.configuration.resource_base_url)) {
      errors.push('Resource Base URL empty');
    }
    if (_.isEmpty(state.configuration.x_fapi_financial_id)) {
      errors.push('x-fapi-financial-id empty');
    }
    if (_.isEmpty(state.configuration.issuer)) {
      errors.push('issuer empty');
    }
    if (_.isEmpty(state.configuration.redirect_url)) {
      errors.push('Redirect URL empty');
    }

    if (!_.isEmpty(errors)) {
      dispatch('status/setErrors', errors, { root: true });
      return false;
    }

    try {
      // NB: We do not care what value this method call returns as long
      // as it does not throw, we know the configuration is valid.
      const { configuration } = state;
      const setShowLoading = flag => dispatch('status/setShowLoading', flag, { root: true });

      await api.validateConfiguration(configuration, setShowLoading);
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_FOUR);

      return true;
    } catch (err) {
      dispatch('status/setErrors', [err], { root: true });
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);

      return false;
    }
  },
  setWizardStep({ commit }, step) {
    commit(types.SET_WIZARD_STEP, step);
  },
};
