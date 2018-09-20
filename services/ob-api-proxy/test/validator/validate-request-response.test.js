const assert = require('assert');

const { runSwaggerValidation, logFormat } = require('../../app/validator/validate-request-response');
const { validate } = require('../../app/validator');
const Balances = require('./fixtures/accounts/balances').default;
const Transactions = require('./fixtures/accounts/transactions').default;

const invalidResponse = () => {
  const bodyInvalid = JSON.parse(Balances.response().body);
  delete bodyInvalid.Data.Balance[0].Amount;

  const responseNew = Balances.response();
  responseNew.body = JSON.stringify(bodyInvalid);

  return responseNew;
};

const details = {
  interactionId: '590bcc25-517c-4caf-a140-077b41ffe095',
  sessionId: '2789f200-4960-11e8-b019-35d9f0621d63',
  authorisationServerId: 'testAuthServerId',
  validationRunId: 'testValidationRunId',
  scope: 'accounts',
  swaggerUris: ['https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-basic.json'],
};

describe('validate', () => {
  before(() => {
    process.env.VALIDATE_RESPONSE = 'true';
  });

  after(() => {
    delete process.env.VALIDATE_RESPONSE;
  });

  describe('without response provided', () => {
    it('returns failedValidation true with message', async () => {
      const response = await validate(Balances.request(), { body: '{}' }, details);

      assert.equal(response.failedValidation, true);
      assert.equal(response.message, 'Response validation failed: failed schema validation');
      assert.deepEqual(response.results.errors, [{
        code: 'OBJECT_MISSING_REQUIRED_PROPERTY',
        message: 'Missing required property: Meta',
        path: [],
      },
      {
        code: 'OBJECT_MISSING_REQUIRED_PROPERTY',
        message: 'Missing required property: Links',
        path: [],
      },
      {
        code: 'OBJECT_MISSING_REQUIRED_PROPERTY',
        message: 'Missing required property: Data',
        path: [],
      }]);
    });
  });

  describe('with valid request and response', () => {
    it('returns failedValidation false', async () => {
      const response = await validate(Balances.request(), Balances.response(), details);
      assert.equal(response.failedValidation, false);
    });
  });

  describe('with valid request and invalid response', () => {
    const validationResults = {
      errors: [{
        code: 'OBJECT_MISSING_REQUIRED_PROPERTY',
        message: 'Missing required property: Amount',
        path: ['Data', 'Balance', '0'],
        description: 'Set of elements used to define the balance details.',
      }],
      warnings: [],
    };

    it('returns failedValidation true with message and errors', async () => {
      const response = await validate(Balances.request(), invalidResponse(), details);
      const expected = {
        failedValidation: true,
        message: 'Response validation failed: failed schema validation',
        results: validationResults,
      };
      assert.deepEqual(response, expected);
    });

    it('and logFormat returns output object for logging', async () => {
      const response = await runSwaggerValidation(Balances.request(), invalidResponse(), details);
      const { failedValidation, message, results } = response.body;
      const output = logFormat(Balances.request(), invalidResponse(), details, response);

      assert.deepEqual(output.request, Balances.request());
      assert.deepEqual(output.response, invalidResponse());
      assert.deepEqual(output.details, details);
      assert.ok(output.report, 'expect output object to contain report property');
      const expectedReport = {
        failedValidation,
        message,
        results,
      };
      assert.deepEqual(output.report, expectedReport);
    });
  });

  describe('without swaggerUris set', () => {
    it('throws error', async () => {
      try {
        const detailsWithoutswaggerUris = Object.assign({}, details, {
          swaggerUris: [],
        });
        await validate(
          Balances.request(), Balances.response(),
          detailsWithoutswaggerUris,
        );
        assert.ok(false);
      } catch (err) {
        assert.equal(err.message, 'checkDetails: swaggerUris missing from validate call');
      }
    });
  });

  describe('with x-swagger-uris header', () => {
    it('returns failedValidation false on basic with basic swaggerUris', async () => {
      const detailsWithSwaggerUris = Object.assign({}, details, {
        swaggerUris: [
          'https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-basic.json',
        ],
      });

      const response = await validate(
        Transactions.valid().request,
        Transactions.valid().basic,
        detailsWithSwaggerUris,
      );

      // console.error('response:', JSON.stringify(response));
      assert.deepEqual(
        response,
        { failedValidation: false },
        JSON.stringify(response),
      );
    }).timeout(1000 * 5);

    it('returns failedValidation true on detail with basic swaggerUris', async () => {
      const detailsWithSwaggerUris = Object.assign({}, details, {
        swaggerUris: [
          'https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-basic.json',
        ],
      });

      const response = await validate(
        Transactions.valid().request,
        Transactions.valid().detail,
        detailsWithSwaggerUris,
      );

      // console.error('response:', JSON.stringify(response));
      assert.deepEqual(response, {
        failedValidation: true,
        results: {
          errors: [{
            code: 'OBJECT_ADDITIONAL_PROPERTIES',
            message: 'Additional properties not allowed: MerchantDetails,Balance,TransactionInformation',
            path: ['Data', 'Transaction', '0'],
            description: 'Provides further details on an entry in the report.',
          }],
          warnings: [],
        },
        message: 'Response validation failed: failed schema validation',
      }, JSON.stringify(response));
    }).timeout(1000 * 5);

    it('returns failedValidation false on detail with detail swaggerUris', async () => {
      const detailsWithSwaggerUris = Object.assign({}, details, {
        swaggerUris: [
          'https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-detail.json',
        ],
      });

      const response = await validate(
        Transactions.valid().request,
        Transactions.valid().detail,
        detailsWithSwaggerUris,
      );

      // console.error('response:', JSON.stringify(response));
      assert.deepEqual(
        response,
        { failedValidation: false },
        JSON.stringify(response),
      );
    }).timeout(1000 * 5);
  });
});
