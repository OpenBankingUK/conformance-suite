import mutations from './mutations';
import * as types from './mutation-types';

describe('User', () => {
  describe('mutations', () => {
    let state;
    const profile = { first_name: 'A' };

    beforeEach(() => {
      state = {};
    });

    it(`${types.USER_SIGNIN} with profile`, () => {
      const expectedState = {
        loading: false,
        profile: { first_name: 'A' },
        signedIn: true,
      };
      mutations[types.USER_SIGNIN](state, profile);
      expect(state).toEqual(expectedState);
    });

    it(`${types.USER_SIGNIN} without profile`, () => {
      const expectedState = {
        loading: false,
        signedIn: true,
      };
      mutations[types.USER_SIGNIN](state);
      expect(state).toEqual(expectedState);
    });

    it(types.USER_SIGNOUT, () => {
      const expectedState = {
        loading: false,
        profile: null,
        signedIn: false,
      };
      mutations[types.USER_SIGNIN](state, profile);
      mutations[types.USER_SIGNOUT](state);
      expect(state).toEqual(expectedState);
    });
  });
});
