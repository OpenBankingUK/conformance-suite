
describe('PSU consent enter config', () => {
  const discoveryTemplateId = '#ob-v3-1-ozone';
  const configTemplate = 'ozone-psu-invalid-client-id-config.json';

  it('sets consent URL', () => {
    cy.selectDiscoveryTemplate(discoveryTemplateId);
    cy.enterConfiguration(configTemplate);

    cy.nextButtonContains('Pending PSU Consent');

    cy.get('.psu-consent-link:first').invoke('attr', 'title').then((consentUrl) => {
      cy.writeFile('consentUrl.txt', consentUrl);
    });
  });
});
