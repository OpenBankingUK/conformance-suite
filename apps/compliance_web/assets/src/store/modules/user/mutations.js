import * as types from './mutation-types';

export default {
  [types.USER_SIGNIN](state, profile) {
    state.signedIn = true;
    state.loading = false;
    if (profile) state.profile = { ...state.profile, ...profile };
  },
  [types.USER_SIGNOUT](state) {
    state.signedIn = false;
    state.loading = false;
    state.profile = null;
  },
};
