import getters from './getters';

describe('User', () => {
  let state;

  beforeEach(() => {
    state = {
      signedIn: false,
      loading: true,
      profile: null,
    };
  });

  describe('getters', () => {
    it('isSignedIn', () => {
      expect(getters.isSignedIn(state)).toEqual(false);
    });

    it('isLoading', () => {
      expect(getters.isLoading(state)).toEqual(true);
    });

    it('getAccessToken', () => {
      expect(getters.getAccessToken(state)).toEqual(null);
      const accessToken = 'some_token';
      state = {
        profile: {
          access_token: accessToken,
        },
      };
      expect(getters.getAccessToken(state)).toEqual(accessToken);
    });
  });
});
