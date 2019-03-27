import URI from 'urijs';
import api from '../../../src/api/consentCallback';

describe('PSU consent granted model bank test case run', () => {
  const discoveryTemplateId = '#ob-v3-1-ozone';
  const configTemplate = 'ozone-psu-config.json';

  it('sets consent URL', () => {
    cy.selectDiscoveryTemplate(discoveryTemplateId);
    cy.enterConfiguration(configTemplate);

    cy.nextButtonContains('Pending PSU Consent');

    // wait for Web socket to be connected:
    cy.get('#ws-connected', { timeout: 16000 });

    cy.readFile('redirectBackUrls.json').then((urls) => {
      const calls = urls.map((redirectBackUrl) => {
        // Use localhost domain to avoid security warnings in browser:
        const url = redirectBackUrl.replace('0.0.0.0', 'localhost').replace('127.0.0.1', 'localhost');
        const uri = new URI(url);
        return {
          method: 'POST',
          url: api.consentCallbackEndpoint(uri),
          body: api.consentParams(uri),
        };
      });

      const req = (list) => {
        const params = list.pop();
        cy
          .request(params)
          .then(() => {
            if (calls.length === 0) {
              cy.runTestCases();
              cy.exportConformanceReport();
              return;
            }
            // else recurse
            req(calls);
          });
      };

      req(calls);
    });
  });
});
