import * as types from './mutation-types';

export default {
  [types.SET_ERRORS](state, errors) {
    state.errors = errors;
  },
  [types.SET_NOTIFICATIONS](state, notifications) {
    state.notifications = notifications;
  },
  [types.PUSH_NOTIFICATION](state, notification) {
    state.notifications.push(notification);
  },
  [types.SET_SHOW_LOADING](state, showLoading) {
    state.showLoading = showLoading;
  },
  [types.SET_SUITE_VERSION](state, version) {
    state.suiteVersion = version;
  },
};
