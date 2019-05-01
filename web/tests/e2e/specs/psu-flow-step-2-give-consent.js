describe('PSU gives consent', () => {
  it('gets redirect back URL', () => {
    cy.readFile('consentUrls.json').then((consentUrls) => {
      const urlsFile = 'redirectBackUrls.json';
      cy.writeFile(urlsFile, '[');
      let firstLine = true;

      consentUrls.forEach((consentUrl) => {
        cy.visit(consentUrl, { timeout: 8000 });
        cy.get('#loginName').type('mits');
        cy.get('#password').type('mits');
        cy.get('button[type="submit"]').click();

        cy.get('input[name="accounts"]').each((input) => { input.click(); });
        cy.get('button[type="submit"]').click();

        cy.contains('a[role="button"]', 'Yes').invoke('attr', 'href').then((href) => {
          const host = 'https://modelobankauth2018.o3bank.co.uk:4101';
          cy.request({
            url: host + href,
            followRedirect: false,
          }).then((response) => {
            const { location } = response.headers;
            const line = firstLine ? `\n  "${location}"` : `,\n  "${location}"`;
            firstLine = false;
            cy.writeFile(urlsFile, line, { flag: 'a' });
          });
        });
      });

      cy.writeFile(urlsFile, `\n]\n`, { flag: 'a' }); // eslint-disable-line
    });
  });
});
