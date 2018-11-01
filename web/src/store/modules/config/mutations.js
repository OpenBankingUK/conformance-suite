import * as types from './mutation-types';

export default {
  [types.SET_CONFIG](state, config) {
    state.main = config;
  },
  [types.SET_DISCOVERY_MODEL](state, discoveryModel) {
    state.discoveryModel = discoveryModel;
  },
  [types.UPDATE_DISCOVERY_MODEL](state, discoveryModel) {
    state.discoveryModel = [
      ...state.discoveryModel,
      discoveryModel,
    ];
  },
  [types.DELETE_DISCOVERY_MODEL](state, discoveryModel) {
    state.discoveryModel =
      state.discoveryModel.filter(item => JSON.stringify(item) !== JSON.stringify(discoveryModel));
  },
};
