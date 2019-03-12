// https://docs.cypress.io/api/introduction/api.html

describe('Headless consent model bank test case run', () => {
  const discoveryTemplate = '#ob-v3-1-ozone-headless';
  const configJsonTab = '#json-view___BV_tab_button__';

  // Note: We can't use async/await with Cypress then() func, as it does not
  // return a Promise.
  // See: https://docs.cypress.io/guides/core-concepts/variables-and-aliases.html#Closures
  it('Gets a PASSED result', () => {
    cy.visit('https://localhost:8443', { timeout: 8000 });
    cy.get(discoveryTemplate).click();
    cy.clickNext();
    cy.get(configJsonTab).click();
    cy.configFixture('ozone-headless-config.json').then((config) => {
      cy.window().then((win) => {
        const editor = win.ace.edit('configuration-editor');
        editor.getSession().setValue(config);
      });
    });

    cy.clickNext();
    cy.contains('a', 'Account and Transaction API Specification');
    cy.nextButtonContains('Run');
    cy.clickNext();
    cy.contains('h6', 'PASSED', { timeout: 8000 });
    cy.nextButtonContains('Next Export', { timeout: 8000 });
    cy.clickNext();
    cy.get('#implementer').type('implementer_example');
    cy.get('#authorised_by').type('authorised_by_example');
    cy.get('#job_title').type('job_title_example');
    cy.get('#has_agreed').click({ force: true });
    cy.nextButtonContains('Export Conformance Report');
    cy.clickNext();
    cy.contains('h5', 'Exported Results');
  });
});
