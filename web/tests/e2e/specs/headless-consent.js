// https://docs.cypress.io/api/introduction/api.html

describe('Headless consent model bank test case run', () => {
  const discoveryTemplateId = '#ob-v3-1-ozone-headless';
  const configTemplate = 'ozone-headless-config.json';

  it('Gets a PASSED result', () => {
    cy.selectDiscoveryTemplate(discoveryTemplateId);
    cy.enterConfiguration(configTemplate);
    cy.runTestCases();
    cy.exportConformanceReport();
  });
});
