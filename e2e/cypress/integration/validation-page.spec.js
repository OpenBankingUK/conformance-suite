import { gapi } from '../utils';

const referenceMockServer = process.env.REFERENCE_MOCK_SERVER === 'localhost'
  ? 'localhost'
  : 'reference-mock-server';

const exampleConfig = {
  authorization_endpoint: `http://${referenceMockServer}:8001/aaaj4NmBD8lQxmLh2O/authorize`,
  client_id: 'spoofClientId',
  client_secret: 'spoofClientSecret',
  fapi_financial_id: 'aaax5nTR33811Qy',
  issuer: 'http://aspsp.example.com',
  redirect_uri: 'http://localhost:8080/tpp/authorized',
  resource_endpoint: `http://${referenceMockServer}:8001`,
  signing_key: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
  signing_kid: 'XXXXXX-XXXXxxxXxXXXxxx_xxxx',
  token_endpoint: `http://${referenceMockServer}:8001/aaaj4NmBD8lQxmLh2O/token`,
  token_endpoint_auth_method: 'client_secret_basic',
  transport_cert: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
  transport_key: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----'
};

describe('Config page', () => {
  describe('Start validation from JSON config', () => {
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
      cy
        .get('.go-to-config')
        .click();
    });

    it('should add JSON config and start validation', () => {
      cy
        .window().then(win => {
          const editor = win.ace.edit('editor');
          editor.getSession().setValue(JSON.stringify(exampleConfig, null, 2))
        })
        .get('.start_validation')
        .click({ force: true });

      cy
        .get('.endpoint')
        .contains('/open-banking/v1.1/payments')
        .parent()
        .next()
        .should('not.contain', 'Failed calls');
      cy
        .get('.endpoint')
        .contains('/open-banking/v1.1/payment-submissions')
        .parent()
        .next()
        .should('not.contain', 'Failed calls');
    });
  });
});
