import * as types from './mutation-types';

export default {
  [types.SET_CONFIG](state, config) {
    state.raw = config;
  },
  [types.SET_PAYLOAD](state, payload) {
    state.payload.raw = payload;
  },
  [types.SUBMIT_CONFIG](state) {
    state.parsed = JSON.parse(state.raw);
    state.payload.parsed = JSON.parse(state.payload.raw);
  },
};
