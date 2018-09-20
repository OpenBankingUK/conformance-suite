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
  submitConfig({ commit }) {
    commit(types.SUBMIT_CONFIG);
  },
  startValidation({ getters, dispatch }) {
    dispatch('resetValidationsRun');
    dispatch('submitConfig');
    dispatch('validations/validate', {
      payload: getters.getPayload,
      config: getters.getConfig,
    }, { root: true });
    router.push('/reports');
  },
};
