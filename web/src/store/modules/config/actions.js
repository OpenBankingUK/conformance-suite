import * as _ from 'lodash';
import * as types from './mutation-types';
import constants from './constants';

import discovery from '../../../api/discovery';
import api from '../../../api';

export default {
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
      const { success, problems } = await discovery.validateDiscoveryConfig(state.discoveryModel);
      if (success) {
        commit(types.DISCOVERY_MODEL_PROBLEMS, null);
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
    if (!_.isEmpty(errors)) {
      dispatch('status/setErrors', errors, { root: true });
      return false;
    }

    try {
      // NB: We do not care what value this method call returns as long
      // as it does not throw, we know the configuration is valid.
      const { configuration } = state;
      await api.validateConfiguration(configuration);
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_FOUR);

      return true;
    } catch (err) {
      dispatch('status/setErrors', [err], { root: true });
      commit(types.SET_WIZARD_STEP, constants.WIZARD.STEP_THREE);

      return false;
    }
  },
  setTestCaseErrors({ commit }, errors) {
    commit(types.SET_TEST_CASES_ERROR, errors);
  },
  setExecutionErrors({ commit }, errors) {
    commit(types.SET_EXECUTION_ERROR, errors);
  },
  setTestCaseResultsErrors({ commit }, errors) {
    commit(types.SET_TEST_CASE_RESULTS_ERROR, errors);
  },
  setWizardStep({ commit }, step) {
    commit(types.SET_WIZARD_STEP, step);
  },
};
