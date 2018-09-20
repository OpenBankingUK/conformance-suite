import Vue from 'vue';
import actions from './actions';
import router from '../../../router';

describe('User', () => {
  describe('actions', () => {
    let setCurrentUSer;
    let dispatch;
    let commit;
    let routerSpy;

    beforeEach(() => {
      dispatch = jest.fn();
      commit = jest.fn();
      routerSpy = jest.spyOn(router, 'push');

      // Mock window.gapi
      global.gapi = {
        load: jest.fn(),
        auth2: {
          init: jest.fn(),
        },
      };

      // Mock this.googleAuth
      actions.googleAuth = {
        currentUser: { get: jest.fn() },
        signOut: jest.fn(),
        signIn: jest.fn(),
      };

      // Set current user, by default is signedIn and has a googleId
      setCurrentUSer = (param = { isSignedIn: true, googleId: 'GOOGLE_ID' }) =>
        actions
          .googleAuth
          .currentUser
          .get
          .mockImplementation(() => ({
            isSignedIn: () => param.isSignedIn,
            getId: () => param.googleId,
          }));
    });

    describe('initGapi', () => {
      it('should init and call gapi.load', async () => {
        global.gapi.load.mockImplementation((auth, cb) => cb());
        await actions.initGapi();
        expect(global.gapi.load).toHaveBeenCalled();
      });
    });

    describe('isSignedIn', () => {
      it('should dispatch signOut if googleId is undefined', async () => {
        setCurrentUSer({ googleId: undefined });
        await actions.isSignedIn({ dispatch, state: {} });
        expect(dispatch).toHaveBeenNthCalledWith(1, 'initGapi');
        expect(dispatch).toHaveBeenNthCalledWith(2, 'signOut');
      });

      it('should dispatch signOut if state.profile is undefined', async () => {
        setCurrentUSer();
        await actions.isSignedIn({ dispatch, state: {} });
        expect(dispatch).toHaveBeenNthCalledWith(1, 'initGapi');
        expect(dispatch).toHaveBeenNthCalledWith(2, 'signOut');
      });

      it('should dispatch verifyToken if state.profile and googleId are set', async () => {
        const state = {
          profile: { googleId: 'GOOGLE_ID', access_token: 'some token' },
        };

        setCurrentUSer();
        await actions.isSignedIn({ dispatch, state });
        expect(dispatch).toHaveBeenNthCalledWith(1, 'initGapi');
        expect(dispatch).toHaveBeenNthCalledWith(2, 'verifyToken');
      });
    });

    describe('signIn', () => {
      beforeEach(() => {
        setCurrentUSer({ isSignedIn: false, googleId: undefined });
        // Mock googleAuth signIn methods
        actions.googleAuth.signIn.mockImplementation(() => ({
          getAuthResponse: () => 'GOOGLE_TOKEN',
          getId: () => 'GOOGLE_ID',
          getBasicProfile: () => ({
            getImageUrl: () => 'GOOGLE_IMAGE_URL',
          }),
        }));
      });

      it('should signin user in the backend and redirect to / if no errors', async () => {
        // Mock POST /user response
        Vue.axios.post.mockResolvedValue({
          data: { profile: { first_name: 'Firstname' } },
        });

        await actions.signIn({ commit, dispatch });
        expect(dispatch).toHaveBeenNthCalledWith(1, 'initGapi');
        expect(actions.googleAuth.signIn).toHaveBeenCalled();
        expect(commit).toHaveBeenCalledWith('USER_SIGNIN', {
          googleId: 'GOOGLE_ID',
          avatar: 'GOOGLE_IMAGE_URL',
          first_name: 'Firstname',
        });
        expect(dispatch).toHaveBeenNthCalledWith(2, 'setAuthorizationHeader');
        expect(routerSpy).toHaveBeenCalledWith('/');
      });

      it('should redirect to /login if errors', async () => {
        Vue.axios.post.mockRejectedValue({ error: 'some error' });

        try {
          await actions.signIn({ commit, dispatch });
        } catch (e) {
          expect(e).toEqual({ error: 'some error' });
          expect(dispatch).toHaveBeenLastCalledWith('signOut');
        }
      });
    });

    describe('signOut', () => {
      it('should commit USER_SIGNOUT and redirect to /login', () => {
        setCurrentUSer({ googleId: undefined });
        actions.signOut({ commit });
        expect(commit).toHaveBeenCalledWith('USER_SIGNOUT');
        expect(routerSpy).toHaveBeenCalledWith('/login');
        expect(actions.googleAuth.signOut).not.toHaveBeenCalled();
      });

      it('should call googleAuth.signOut() if googleId', () => {
        setCurrentUSer();
        actions.signOut({ commit });
        expect(commit).toHaveBeenCalledWith('USER_SIGNOUT');
        expect(routerSpy).toHaveBeenCalledWith('/login');
        expect(actions.googleAuth.signOut).toHaveBeenCalled();
      });
    });

    describe('verifyToken', () => {
      it('should dispatch setAuthorizationHeader and signIn the user', async () => {
        Vue.axios.get.mockResolvedValue({
          data: { user: { first_name: 'Firstname' } },
        });

        await actions.verifyToken({ commit, dispatch, state: { profile: { access_token: 'TOKEN' } } });
        expect(dispatch).toHaveBeenCalledWith('setAuthorizationHeader');
        expect(commit).toHaveBeenCalledWith('USER_SIGNIN', { first_name: 'Firstname' });
      });

      it('should signOut if GET /user returns an error', async () => {
        Vue.axios.get.mockRejectedValue({ error: 'some error' });

        try {
          await actions.verifyToken({ commit, dispatch, state: { profile: { access_token: 'TOKEN' } } });
        } catch (e) {
          expect(e).toEqual({ error: 'some error' });
          expect(dispatch).toHaveBeenCalledWith('signOut');
        }
      });

      it('should signOut if no state.profile', async () => {
        try {
          await actions.verifyToken({ commit, dispatch, state: {} });
        } catch (e) {
          expect(dispatch).toHaveBeenCalledWith('signOut');
        }
      });
    });

    describe('setAuthorizationHeader', () => {
      it('should set the global Authorization Headers', () => {
        actions.setAuthorizationHeader({ state: { profile: { access_token: 'TOKEN' } } });
        expect(Vue.axios.defaults.headers.common.Authorization).toBe('Bearer TOKEN');
      });
    });
  });
});
