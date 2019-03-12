// https://docs.cypress.io/api/introduction/api.html

describe('PSU consent model bank test case run', () => {
  const discoveryTemplateId = '#ob-v3-1-ozone';
  const configTemplate = 'ozone-psu-config.json';

  it('Gets a PASSED result', () => {
    cy.selectDiscoveryTemplate(discoveryTemplateId);
    cy.enterConfiguration(configTemplate);
    cy.runTestCases();
    cy.exportConformanceReport();
  });
});
