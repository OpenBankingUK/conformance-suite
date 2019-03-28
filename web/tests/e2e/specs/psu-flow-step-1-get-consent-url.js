
describe('PSU consent enter config', () => {
  const discoveryTemplateId = '#ob-v3-1-ozone';
  const configTemplate = 'ozone-psu-config.json';

  it('sets consent URL', () => {
    const urlsFile = 'consentUrls.json';
    cy.writeFile(urlsFile, '[');
    cy.selectDiscoveryTemplate(discoveryTemplateId);
    cy.enterConfiguration(configTemplate);

    cy.nextButtonContains('Pending PSU Consent');
    let firstLine = true;

    cy.get('.psu-consent-link').each((link) => {
      cy.wrap(link).invoke('attr', 'title').then((consentUrl) => {
        const line = firstLine ? `\n  "${consentUrl}"` : `,\n  "${consentUrl}"`;
        firstLine = false;
        cy.writeFile(urlsFile, line, { flag: 'a' });
      });
    });
    cy.writeFile(urlsFile, `\n]\n`, { flag: 'a' }); // eslint-disable-line
  });
});
