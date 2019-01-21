import * as types from './mutation-types';

export default {
  setErrors({ commit }, errors) {
    commit(types.SET_ERRORS, errors);
  },
};
