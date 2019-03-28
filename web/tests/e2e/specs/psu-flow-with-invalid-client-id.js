
describe('PSU consent enter config', () => {
  const discoveryTemplateId = '#ob-v3-1-ozone';
  const configTemplate = 'ozone-psu-invalid-client-id-config.json';

  it('sets consent URL', () => {
    cy.selectDiscoveryTemplate(discoveryTemplateId);
    cy.enterConfiguration(configTemplate);

    cy.contains('div.alert', 'invalid_client');
    cy.contains('div.alert', 'invalid authorization token');
    cy.contains('div.alert', 'ClientCredential Grant: HTTP Status code does not match: expected 200 got 400');
    cy.nextButtonContains('Pending PSU Consent');
  });
});
