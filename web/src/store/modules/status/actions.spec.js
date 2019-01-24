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
