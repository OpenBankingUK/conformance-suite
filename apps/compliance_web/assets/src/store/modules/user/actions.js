import Vue from 'vue';
import router from '../../../router';
import * as types from './mutation-types';

// https://developers.google.com/identity/sign-in/web/reference
export default {
  initGapi() {
    return new Promise((resolve) => {
      window.gapi.load('auth2', () => {
        this.googleAuth = window.gapi.auth2.init();
        resolve();
      });
    });
  },
  async isSignedIn({ dispatch, state }) {
    await dispatch('initGapi');
    await this.googleAuth;
    const currentUser = this.googleAuth.currentUser.get();
    const googleId = currentUser.getId();

    if (!googleId) return dispatch('signOut');
    if (
      state.profile &&
      state.profile.googleId &&
      state.profile.googleId === googleId &&
      state.profile.access_token
    ) {
      return dispatch('verifyToken');
    }

    return dispatch('signOut');
  },
  async signIn({ dispatch, commit }) {
    try {
      await dispatch('initGapi');
      const resp = await this.googleAuth.signIn();
      const { data: { profile } } = await Vue.axios.post('/auth', { id_token: resp.getAuthResponse().id_token });
      commit(types.USER_SIGNIN, {
        googleId: resp.getId(),
        avatar: resp.getBasicProfile().getImageUrl(),
        ...profile,
      });
      dispatch('setAuthorizationHeader');
      return router.push('/');
    } catch (e) {
      return dispatch('signOut');
    }
  },
  signOut({ commit }) {
    if (this.googleAuth.currentUser.get().isSignedIn()) this.googleAuth.signOut();
    commit(types.USER_SIGNOUT);
    return router.push('/login');
  },
  async verifyToken({ commit, state, dispatch }) {
    try {
      const { data: { user } } = await Vue.axios.get('/user', {
        headers: { Authorization: `Bearer ${state.profile.access_token}` },
      });
      dispatch('setAuthorizationHeader');
      return commit(types.USER_SIGNIN, user);
    } catch (e) {
      return dispatch('signOut');
    }
  },
  setAuthorizationHeader({ state }) {
    Vue.axios.defaults.headers.common.Authorization = `Bearer ${state.profile.access_token}`;
  },
};
