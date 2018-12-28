import * as _ from 'lodash';
import * as types from './mutation-types';
import router from '../../../router';
import discovery from '../../../api/discovery';
import api from '../../../api';

export default {
  setDiscoveryModel({ commit }, editorString) {
    try {
      const discoveryModel = JSON.parse(editorString);
      commit(types.SET_DISCOVERY_MODEL, discoveryModel);
      commit(types.DISCOVERY_MODEL_PROBLEMS, null);
    } catch (e) {
      const problems = [{
        key: null,
        error: e.message,
      }];
      commit(types.DISCOVERY_MODEL_PROBLEMS, problems);
    }
  },
  setConfig({ commit }, config) {
    commit(types.SET_CONFIG, config);
  },
  resetValidationsRun({ commit }) {
    // reset validationRunId and lastUpdate for new validation
    commit('reporter/SET_WEBSOCKET_LAST_UPDATE', null, { root: true });
    commit('validations/SET_VALIDATION_DISCOVERY_MODEL', null, { root: true });
  },
  startValidation({ getters, dispatch }) {
    dispatch('resetValidationsRun');
    dispatch('validations/validate', {
      discoveryModel: getters.getDiscoveryModel,
      config: getters.getConfig,
    }, { root: true });
    router.push('/reports');
  },
  setDiscoveryModelProblems({ commit }, problems) {
    commit(types.DISCOVERY_MODEL_PROBLEMS, problems);
  },
  async validateDiscoveryConfig({ commit, state }) {
    try {
      const { success, problems } = await discovery.validateDiscoveryConfig(state.discoveryModel);
      if (success) {
        commit(types.DISCOVERY_MODEL_PROBLEMS, null);
      } else {
        commit(types.DISCOVERY_MODEL_PROBLEMS, problems);
      }
    } catch (e) {
      commit(types.DISCOVERY_MODEL_PROBLEMS, [{
        key: null,
        error: e.message,
      }]);
    }
    return null;
  },

  setConfigurationSigningPrivate({ commit }, signingPrivate) {
    commit(types.SET_CONFIGURATION_SIGNING_PRIVATE, signingPrivate);
  },
  setConfigurationSigningPublic({ commit }, signingPublic) {
    commit(types.SET_CONFIGURATION_SIGNING_PUBLIC, signingPublic);
  },
  setConfigurationTransportPrivate({ commit }, transportPrivate) {
    commit(types.SET_CONFIGURATION_TRANSPORT_PRIVATE, transportPrivate);
  },
  setConfigurationTransportPublic({ commit }, transportPublic) {
    commit(types.SET_CONFIGURATION_TRANSPORT_PUBLIC, transportPublic);
  },
  async validateConfiguration({ commit, state }) {
    commit(types.CLEAR_CONFIGURATION_ERRORS);

    if (_.isEmpty(state.configuration.signing_private)) {
      commit(types.ADD_CONFIGURATION_ERRORS, 'Signing Private Certificate (.key) empty');
    }
    if (_.isEmpty(state.configuration.signing_public)) {
      commit(types.ADD_CONFIGURATION_ERRORS, 'Signing Public Certificate (.pem) empty');
    }
    if (_.isEmpty(state.configuration.transport_private)) {
      commit(types.ADD_CONFIGURATION_ERRORS, 'Transport Private Certificate (.key) empty');
    }
    if (_.isEmpty(state.configuration.transport_public)) {
      commit(types.ADD_CONFIGURATION_ERRORS, 'Transport Public Certificate (.pem) empty');
    }

    if (!_.isEmpty(state.errors.configuration)) {
      return false;
    }

    try {
      // NB: We do not care what value this method call returns as long
      // as it does not throw, we know the configuration is valid.
      const { configuration } = state;
      await api.validateConfiguration(configuration);

      return true;
    } catch (err) {
      commit(types.SET_CONFIGURATION_ERRORS, [
        err,
      ]);

      return false;
    }
  },
};
