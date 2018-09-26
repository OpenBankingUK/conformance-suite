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

const fillForm = () => {
  cy
    .get('input[name=authorization_endpoint]')
    .clear()
    .type(exampleConfig.authorization_endpoint, { delay: 0 })
    .get('input[name=resource_endpoint]')
    .clear()
    .type(exampleConfig.resource_endpoint, { delay: 0 })
    .get('input[name=token_endpoint]')
    .clear()
    .type(exampleConfig.token_endpoint, { delay: 0 });
};

const goToPayload = () => {
  cy
    .get('.validation-button-container .next')
    .click()
    .get('.validation-button-container .next')
    .click()
    .get('.validation-button-container .next')
    .click();
};

const expectNoFailedCalls = () => {
  const endpoints = [
    '/open-banking/v1.1/payments',
    '/open-banking/v1.1/payment-submissions',
  ];

  endpoints.forEach((endpoint) => {
    cy
      .get(`[data-endpoint="${endpoint}"]`)
      .should('contain', 'Total calls')
      .should('not.contain', 'Failed calls');
  });
};

describe('Config page', () => {
  describe('Start validation from JSON/WebForm config', () => {
    beforeEach(() => {
      const googleId = 'GOOGLEID';
      cy
        .createUser(googleId)
        .visit('/', {
          onBeforeLoad(win) {
            // mock gapi with a signed in user
            win.gapi = gapi(googleId);
          }
        })
        .get('.go-to-config')
        .click();
    });

    it('should switch to JSON tab, change config and start validation', () => {
      cy
        .get('.config-header button')
        .click()
        .window().then(win => {
          const editor = win.ace.edit('editor');
          editor.getSession().setValue(JSON.stringify(exampleConfig, null, 2))
        })
        .get('.start_validation')
        .click({ force: true });
        expectNoFailedCalls();
    });

    it('should switch between json / webform view', () => {
      cy
        .get('.config-header button')
        .click()
        .get('.editor')
        .should('be.visible')
        .get('.config-header button')
        .click()
        .get('.config .ant-steps-item-title')
        .should('contain', 'ASPSP');
    });

    it('should fill the webform view and start validation', () => {
      fillForm();
      goToPayload();
      cy.get('.validation-button-container .start').click();
      expectNoFailedCalls();
    });

    it('should fill the webform view, add new payload and start validation', () => {
      fillForm();
      goToPayload();
      cy
        .get('.add-payment .ant-select')
        .click({ force: true })
        .get('.ant-select-dropdown-menu-item')
        .contains('1.1')
        .click()
        .get('input[name=name]')
        .type('Firstname Lastname', { delay: 0 })
        .get('input[name=sort_code]')
        .type('123456', { delay: 0 })
        .get('input[name=account_number]')
        .type('12345678', { delay: 0 })
        .get('input[name=amount]')
        .type('100', { delay: 0 })
        .get('.add-payment button')
        .click()
        .get('.payload-item')
        .eq(4)
        .should('contain', 'Firstname')
        .get('.validation-button-container .start')
        .click();
      expectNoFailedCalls();
    });

    it('should remove one payload and start validation', () => {
      fillForm();
      goToPayload();
      cy
        .get('.payload-item')
        .should('have.length', 5)
        .eq(2)
        .should('contain', 'Sam Morse')
        .get('.payload-item')
        .eq(3)
        .should('contain', 'Michael Burnham')
        .get('.ant-btn-danger[data-item=payment-1]')
        .click()
        .get('.payload-item')
        .should('have.length', 4)
        .get('.payload-item')
        .eq(3)
        .should('not.contain', 'Michael Burnham')
        .get('.validation-button-container .start')
        .click();
      expectNoFailedCalls();
    });
  });
});
