import { gapi } from '../utils';

describe('Login page', () => {
  it('Login component layout', () => {
    cy
      .visit('/', {
        onBeforeLoad(win) {
          // mock gapi with a non signed in user
          win.gapi = gapi('', false);
        }
      });
    cy
      .location().then(loc => {
        // check the redirect to /login
        expect(loc.pathname).to.eq('/login')
      });
    cy
      .get('.login')
      .should('contain', 'Sign in with Google');
  });

  describe('Successful login', () => {
    const googleId = 'GOOGLEID';
    beforeEach(() => {
      cy
        .createUser(googleId)
        .visit('/', {
          onBeforeLoad(win) {
            // mock gapi with a signed in user
            win.gapi = gapi(googleId);
          }
        });
    });

    it('As loggedIn user I should land to /', () => {
      cy
        .location().then(loc => {
          expect(loc.pathname).to.eq('/')
        });
      cy
        .get('.go-to-config')
        .contains('Start a validation run');
    });

    it('Sign out redirects to /login', () => {
      cy
        .location().then(loc => {
          expect(loc.pathname).to.eq('/')
        });
      cy
        .get('.navbar .avatar')
        .click({ force: true });
      cy
        .get('.user-menu')
        .contains('Sign out')
        .click({ force: true });
      cy
        .location().then(loc => {
          expect(loc.pathname).to.eq('/login')
        });
    });
  });
});
