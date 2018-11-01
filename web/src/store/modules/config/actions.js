import * as types from './mutation-types';
import router from '../../../router';

export default {
  setDiscoveryModel({ commit }, discoveryModel) {
    commit(types.SET_DISCOVERY_MODEL, discoveryModel);
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
};
