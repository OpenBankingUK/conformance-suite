import URI from 'urijs';

describe('PSU consent granted model bank test case run', () => {
  const discoveryTemplateId = '#ob-v3-1-ozone';
  const configTemplate = 'ozone-psu-config.json';

  it('sets consent URL', () => {
    cy.selectDiscoveryTemplate(discoveryTemplateId);
    cy.enterConfiguration(configTemplate);

    cy.nextButtonContains('Pending PSU Consent');

    cy.get('#ws-connected', { timeout: 16000 });

    cy.readFile('redirectBackUrl.txt').then((redirectBackUrl) => {
      const url = redirectBackUrl.replace('0.0.0.0', 'localhost').replace('127.0.0.1', 'localhost');
      const uri = new URI(url);

      let callbackUrl;
      let params;
      if (uri.fragment().length > 0) {
        callbackUrl = '/api/redirect/fragment/ok';
        params = URI.parseQuery(uri.fragment());
      }
      if (uri.query().length > 0) {
        callbackUrl = '/api/redirect/query/ok';
        params = URI.parseQuery(uri.query());
      }
      cy.request({
        url: callbackUrl,
        method: 'POST',
        body: params,
      }).then((response) => {
        console.log(response.status); // eslint-disable-line
        cy.runTestCases();
        cy.exportConformanceReport();
      });
    });
  });
});
