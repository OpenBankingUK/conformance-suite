import * as types from './mutation-types';

export default {
  clearErrors({ commit, state }) {
    if (state.errors.length > 0) {
      commit(types.SET_ERRORS, []);
    }
  },
  setErrors({ commit }, errors) {
    if (errors) {
      commit(types.SET_ERRORS, errors);
    }
  },
};
