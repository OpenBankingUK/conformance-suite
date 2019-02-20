import actions from './actions';
import * as types from './mutation-types';

const stateWithErrors = {
  errors: [
    new Error('Error message'),
    'text message',
  ],
};

const stateNoErrors = {
  errors: [],
};

describe('setErrors', () => {
  let commit;

  beforeEach(() => {
    commit = jest.fn();
  });

  it('commits errors array', () => {
    actions.setErrors({ commit, state: stateNoErrors }, stateWithErrors.errors);
    expect(commit).toHaveBeenCalledWith(types.SET_ERRORS, stateWithErrors.errors);
  });

  it('commits empty array', () => {
    actions.setErrors({ commit, state: stateWithErrors }, []);
    expect(commit).toHaveBeenCalledWith(types.SET_ERRORS, []);
  });

  it('does not commit null', () => {
    actions.setErrors({ commit, state: stateWithErrors }, null);
    expect(commit).not.toHaveBeenCalledWith(types.SET_ERRORS, null);
  });
});

describe('clearErrors', () => {
  let commit;

  beforeEach(() => {
    commit = jest.fn();
  });

  it('does not commit when errors not present', () => {
    actions.clearErrors({ commit, state: stateNoErrors });
    expect(commit).not.toHaveBeenCalledWith(types.SET_ERRORS, []);
  });

  it('commits empty array when errors present', () => {
    actions.clearErrors({ commit, state: stateWithErrors });
    expect(commit).toHaveBeenCalledWith(types.SET_ERRORS, []);
  });
});

const stateWithNotifications = {
  notifications: [{
    message: 'sample-message',
    extURL: 'https://www.example',
  }],
};

const stateNoNotifications = {
  notifications: [],
};

describe('pushNotification', () => {
  let commit;

  beforeEach(() => {
    commit = jest.fn();
  });

  it('notification count increments by 1 upon pushing a notification', () => {
    const n = stateWithNotifications.notifications[0];
    actions.pushNotification({ commit, state: stateNoNotifications }, n);
    expect(commit).toHaveBeenCalledWith(types.PUSH_NOTIFICATION, n);
  });

  it('does not commit if notification is null, when pushing a notification', () => {
    actions.pushNotification({ commit, state: stateNoNotifications }, null);
    expect(commit).not.toHaveBeenCalledWith(types.PUSH_NOTIFICATION, null);
  });
});

describe('clearNotifications', () => {
  let commit;

  beforeEach(() => {
    commit = jest.fn();
  });

  it('clearing notifications results in 0 notifications', () => {
    actions.clearNotifications({ commit, state: stateWithNotifications });
    expect(commit).toHaveBeenCalledWith(types.SET_NOTIFICATIONS, []);
  });
  it('does not commit when no notifications present', () => {
    actions.clearNotifications({ commit, state: stateNoNotifications });
    expect(commit).not.toHaveBeenCalledWith(types.SET_NOTIFICATIONS, []);
  });

  it('commits empty array when errors present', () => {
    actions.clearNotifications({ commit, state: stateWithNotifications });
    expect(commit).toHaveBeenCalledWith(types.SET_NOTIFICATIONS, []);
  });
});

describe('setShowLoading', () => {
  let commit;

  beforeEach(() => {
    commit = jest.fn();
  });

  it('commits showLoading flag', () => {
    const state = { showLoading: false };
    actions.setShowLoading({ commit, state }, true);
    expect(commit).toHaveBeenCalledWith(types.SET_SHOW_LOADING, true);
  });
});
