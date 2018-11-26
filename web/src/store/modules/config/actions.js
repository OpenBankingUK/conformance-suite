import * as types from './mutation-types';
import router from '../../../router';
import DiscoveryExample from './discovery-example.json';

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
  resetDiscoveryConfig({ commit }) {
    // TODO: Maybe validate the default example ... not sure.
    commit(types.DISCOVERY_MODEL_RESET, DiscoveryExample);
    commit(types.DISCOVERY_MODEL_PROBLEMS, null);
  },
  validateDiscoveryConfig({ commit }) {
    // TODO: Remove hardcoded `problems` and call backend instead.
    const problems = `discoveryItemIndex=0, missing mandatory endpoint Method=POST, Path=/account-access-consents
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}/transactions
discoveryItemIndex=0, missing mandatory endpoint Method=POST, Path=/account-access-consents
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}/transactions
discoveryItemIndex=0, missing mandatory endpoint Method=POST, Path=/account-access-consents
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}/transactions
discoveryItemIndex=0, missing mandatory endpoint Method=POST, Path=/account-access-consents
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}/transactions
discoveryItemIndex=0, missing mandatory endpoint Method=POST, Path=/account-access-consents
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}/transactions
discoveryItemIndex=0, missing mandatory endpoint Method=POST, Path=/account-access-consents
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}/transactions
discoveryItemIndex=0, missing mandatory endpoint Method=POST, Path=/account-access-consents
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}/transactions`;
    commit(types.DISCOVERY_MODEL_PROBLEMS, problems);

    return Promise.resolve({});
  },
};
