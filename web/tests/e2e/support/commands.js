// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This is will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })

// Loads config fixture JSON and replaces ENV variables with ENV variable values.
Cypress.Commands.add('configFixture', (file) => {
  cy.fixture(file).then((template) => {
    cy.replaceEnvVarConfig(template).then((config) => {
      const indentedConfig = JSON.stringify(config, null, 2);
      return indentedConfig;
    });
  });
});

// Replace ENV variables in config template, with ENV variable values.
//
// There are several ways to set ENV variables, including:
//   1. Create a web/cypress.env.json file, i.e.
//      > cp web/cypress.sample.env.json web/cypress.env.json
//      Edit to add config values. E.g.
//   {
//     "OZONE_CLIENT_ID": "example_client_id",
//     "OZONE_CLIENT_SECRET": "example_client_secret",
//     "OZONE_SIGNING_PRIVATE": "-----BEGIN PRIVATE KEY-----\nexample\n-----END PRIVATE KEY-----\n",
//     "OZONE_SIGNING_PUBLIC": "-----BEGIN CERTIFICATE-----\nexample\n-----END CERTIFICATE-----\n",
//     "OZONE_TRANSPORT_PRIVATE": "-----BEGIN PRIVATE KEY-----\nexample\n-----END PRIVATE KEY-----\n",
//     "OZONE_TRANSPORT_PUBLIC": "-----BEGIN CERTIFICATE-----\nexample\n-----END CERTIFICATE-----\n"
//     ...
//   2. Or export ENV prefixed as `CYPRESS_*`
//
// For more options see: https://docs.cypress.io/guides/guides/environment-variables.html#Setting
Cypress.Commands.add('replaceEnvVarConfig', (config) => {
  const replaceFields = [
    'signing_private',
    'signing_public',
    'transport_private',
    'transport_public',
    'client_id',
    'client_secret',
  ];
  replaceFields.forEach((field) => {
    const envVar = config[field];
    const value = Cypress.env(envVar);
    config[field] = value; // eslint-disable-line
  });
  return config;
});

const nextButton = '#next';

Cypress.Commands.add('clickNext', () => {
  cy.get(nextButton).click();
});

Cypress.Commands.add('nextButtonContains', (text, opts) => {
  cy.contains(nextButton, text, opts);
});

Cypress.Commands.add('selectDiscoveryTemplate', (templateSelectorId) => {
  cy.visit('https://localhost:8443', { timeout: 8000 });
  cy.get(templateSelectorId).click();
  cy.clickNext();
});

Cypress.Commands.add('enterConfiguration', (configTemplateFile) => {
  const configJsonTab = '#json-view___BV_tab_button__';
  cy.get(configJsonTab).click();

  // Note: We can't use async/await with Cypress then() func, as it does not
  // return a Promise.
  // See: https://docs.cypress.io/guides/core-concepts/variables-and-aliases.html#Closures
  cy.configFixture(configTemplateFile).then((config) => {
    cy.window().then((win) => {
      const editor = win.ace.edit('configuration-editor');
      editor.getSession().setValue(config);
    });
  });
  cy.clickNext();
  cy.contains('a', 'Account and Transaction API Specification');
});

Cypress.Commands.add('runTestCases', () => {
  cy.nextButtonContains('Run');
  cy.clickNext();

  cy.contains('h6', 'PASSED', { timeout: 8000 });
  cy.nextButtonContains('Next Export', { timeout: 8000 });
});

Cypress.Commands.add('exportConformanceReport', () => {
  cy.clickNext();
  cy.get('#implementer').type('implementer_example');
  cy.get('#authorised_by').type('authorised_by_example');
  cy.get('#job_title').type('job_title_example');
  cy.get('#has_agreed').click({ force: true });
  cy.nextButtonContains('Export Conformance Report');
  cy.clickNext();
  cy.contains('h5', 'Exported Results');
});
