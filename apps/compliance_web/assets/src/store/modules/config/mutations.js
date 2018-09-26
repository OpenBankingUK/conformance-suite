import * as types from './mutation-types';

export default {
  [types.SET_CONFIG](state, config) {
    state.main = config;
  },
  [types.SET_PAYLOAD](state, payload) {
    state.payload = payload;
  },
  [types.UPDATE_PAYLOAD](state, payload) {
    state.payload = [
      ...state.payload,
      payload,
    ];
  },
  [types.DELETE_PAYLOAD](state, payload) {
    state.payload = state.payload.filter(item => JSON.stringify(item) !== JSON.stringify(payload));
  },
};
