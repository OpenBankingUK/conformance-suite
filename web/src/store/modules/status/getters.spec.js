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
