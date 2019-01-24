import * as types from './mutation-types';

export default {
  [types.SET_ERRORS](state, errors) {
    state.errors = errors;
  },
};
