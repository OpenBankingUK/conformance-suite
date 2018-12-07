import * as types from './mutation-types';
import router from '../../../router';
import DiscoveryExample from './discovery-example.json';
import discovery from '../../../api/discovery';

export default {
  setDiscoveryModel({ commit }, editorString) {
    try {
      const discoveryModel = JSON.parse(editorString);
      commit(types.SET_DISCOVERY_MODEL, discoveryModel);
      commit(types.DISCOVERY_MODEL_PROBLEMS, null);
    } catch (e) {
      commit(types.DISCOVERY_MODEL_PROBLEMS, [e.message]);
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
  updateDiscoveryModel({ commit }, discoveryModel) {
    commit(types.UPDATE_DISCOVERY_MODEL, discoveryModel);
  },
  deleteDiscoveryModel({ commit }, discoveryModel) {
    commit(types.DELETE_DISCOVERY_MODEL, discoveryModel);
  },
  resetDiscoveryConfig({ commit }) {
    // TODO: Maybe validate the default example ... not sure.
    commit(types.DISCOVERY_MODEL_RESET, DiscoveryExample);
    commit(types.DISCOVERY_MODEL_PROBLEMS, null);
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
      commit(types.DISCOVERY_MODEL_PROBLEMS, [e.message]);
    }
    return null;
  },
};
