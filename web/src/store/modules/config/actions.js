import * as types from './mutation-types';
import router from '../../../router';

export default {
  setPayload({ commit }, payload) {
    commit(types.SET_PAYLOAD, payload);
  },
  setConfig({ commit }, config) {
    commit(types.SET_CONFIG, config);
  },
  resetValidationsRun({ commit }) {
    // reset validationRunId and lastUpdate for new validation
    commit('reporter/SET_WEBSOCKET_LAST_UPDATE', null, { root: true });
    commit('validations/SET_VALIDATION_PAYLOAD', null, { root: true });
  },
  startValidation({ getters, dispatch }) {
    dispatch('resetValidationsRun');
    dispatch('validations/validate', {
      payload: getters.getPayload,
      config: getters.getConfig,
    }, { root: true });
    router.push('/reports');
  },
  updatePayload({ commit }, payload) {
    commit(types.UPDATE_PAYLOAD, payload);
  },
  deletePayload({ commit }, payload) {
    commit(types.DELETE_PAYLOAD, payload);
  },
};
