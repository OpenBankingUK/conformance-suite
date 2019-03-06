// https://docs.cypress.io/api/introduction/api.html

const exampleConfig = {
  signing_private: '',
  signing_public: '',
  transport_private: '',
  transport_public: '',
  client_id: '',
  client_secret: '',
  token_endpoint: 'https://modelobank2018.o3bank.co.uk:4201/token',
  token_endpoint_auth_method: 'client_secret_basic',
  authorization_endpoint: 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
  resource_base_url: 'https://modelobank2018.o3bank.co.uk:4501',
  x_fapi_financial_id: '0015800001041RHAAY',
  issuer: 'https://modelobankauth2018.o3bank.co.uk:4101',
  redirect_url: 'https://0.0.0.0:8443/conformancesuite/callback',
};

describe('Headless consent model bank test case run', () => {
  const discoveryTemplate = '#ob-v3-1-ozone-headless';
  const nextButton = '#next';
  const configJsonTab = '#json-view___BV_tab_button__';

  it('Gets results', () => {
    cy.visit('https://localhost:8443');
    cy.get(discoveryTemplate).click();
    cy.get(nextButton).click();
    cy.get(configJsonTab).click();
    cy.window().then((win) => {
      const editor = win.ace.edit('configuration-editor');
      editor.getSession().setValue(JSON.stringify(exampleConfig, null, 2));
    });
    cy.get(nextButton).click();
    cy.contains('a', 'Account and Transaction API Specification');
    cy.contains(nextButton, 'Run');
    cy.get(nextButton).click();
  });
});
