import URI from 'urijs';

describe('PSU consent granted model bank test case run', () => {
  const discoveryTemplateId = '#ob-v3-1-ozone';
  const configTemplate = 'ozone-psu-config.json';

  it('sets consent URL', () => {
    cy.selectDiscoveryTemplate(discoveryTemplateId);
    cy.enterConfiguration(configTemplate);

    cy.nextButtonContains('Pending PSU Consent');

    // wait for Web socket to be connected:
    cy.get('#ws-connected', { timeout: 16000 });

    cy.readFile('redirectBackUrl.txt').then((redirectBackUrl) => {
      // Use localhost domain to avoid security warnings in browser:
      const url = redirectBackUrl.replace('0.0.0.0', 'localhost').replace('127.0.0.1', 'localhost');
      const uri = new URI(url);
      const isFragment = uri.fragment().length > 0;
      const isQuery = uri.query().length > 0;

      const params = { method: 'POST' };

      if (isFragment) {
        Object.assign(params, {
          url: '/api/redirect/fragment/ok',
          body: URI.parseQuery(uri.fragment()),
        });
      }
      if (isQuery) {
        Object.assign(params, {
          url: '/api/redirect/query/ok',
          body: URI.parseQuery(uri.query()),
        });
      }

      cy.request(params).then((response) => {
        console.log(response.status); // eslint-disable-line
        cy.runTestCases();
        cy.exportConformanceReport();
      });
    });
  });
});
