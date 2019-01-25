import getters from './getters';

const stateWithErrors = {
  errors: [
    new Error('Error message'),
    'text message',
  ],
};

const stateNoErrors = {
  errors: [],
};

const stateNotifications = {
  notifications: [{
    message: 'sample-message',
    extURL: 'https://www.example',
  }],
};

const stateNoNotifications = { notifications: [] };

describe('errorMessages', () => {
  it('returns array of error messages when errors', () => {
    const list = getters.errorMessages(stateWithErrors);
    expect(list[0]).toEqual('Error message');
    expect(list[1]).toEqual('text message');
  });

  it('returns empty array when no errors', () => {
    const list = getters.errorMessages(stateNoErrors);
    expect(list).toEqual([]);
  });
});

describe('hasErrors', () => {
  it('true when errors', () => {
    const flag = getters.hasErrors(stateWithErrors);
    expect(flag).toBe(true);
  });

  it('false when no errors', () => {
    const flag = getters.hasErrors(stateNoErrors);
    expect(flag).toBe(false);
  });
});

describe('hasNotifications', () => {
  it('true when notifications present', () => {
    const flag = getters.hasNotifications(stateNotifications);
    expect(flag).toBe(true);
  });

  it('false when notifications present', () => {
    const flag = getters.hasNotifications(stateNoNotifications);
    expect(flag).toBe(false);
  });
});

describe('notifications', () => {
  it('returns array of notifications when notifications present', () => {
    const list = getters.notifications(stateNoNotifications);
    expect(list[0]).toEqual(stateNoNotifications.notifications[0]);
  });

  it('returns empty array of notifications when no notifications present', () => {
    const list = getters.notifications(stateNoNotifications);
    expect(list).toEqual([]);
  });
});
