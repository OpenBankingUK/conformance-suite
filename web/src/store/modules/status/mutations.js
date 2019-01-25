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
};
