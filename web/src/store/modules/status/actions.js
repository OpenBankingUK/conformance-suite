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
  setShowLoading({ commit }, showLoading) {
    commit(types.SET_SHOW_LOADING, showLoading);
  },
  pushNotification({ commit, state }, notification) {
    if (state.notifications && notification) {
      commit(types.PUSH_NOTIFICATION, notification);
    }
  },
  clearNotifications({ commit, state }) {
    if (state.notifications && state.notifications.length > 0) {
      commit(types.SET_NOTIFICATIONS, []);
    }
  },
};
