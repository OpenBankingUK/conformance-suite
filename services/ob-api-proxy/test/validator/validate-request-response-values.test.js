const assert = require('assert');

const { validateResponse } = require('../../app/validator/validator-response-request-values');
const { validate } = require('../../app/validator');
const AccountRequests = require('./fixtures/accounts/account-requests').default;

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

  describe('validateResponse', () => {
    describe('payments', () => {
      describe('/open-banking/v1.1/payments POST', () => {
        it('isValid is true, when request.Data.Initiation == response.Data.Initiation && request.Risk == response.Risk', async () => {
          const request = {
            method: 'POST',
            path: '/open-banking/v1.1/payments',
            body: {
              Data: {
                Initiation: {
                  InstructionIdentification: '15abd6c0-18f1-4257-a765-c15ef8bf1c',
                  EndToEndIdentification: '8a30c4fe-a779-436f-b231-f21c05bd22',
                  InstructedAmount: {
                    Currency: 'GBP',
                    Amount: '10.00',
                  },
                  CreditorAccount: {
                    SchemeName: 'SortCodeAccountNumber',
                    Name: 'Sam Morse',
                    Identification: '11111112345678',
                  },
                },
              },
              Risk: {},
            },
          };
          const response = {
            body: '{"Data":{"PaymentId":"f5f92b5e-5fcd-4476-86cd-9f37ea015ffe","Status":"AcceptedTechnicalValidation","CreationDateTime":"2018-08-21T07:57:47+00:00","Initiation":{"InstructionIdentification":"15abd6c0-18f1-4257-a765-c15ef8bf1c","EndToEndIdentification":"8a30c4fe-a779-436f-b231-f21c05bd22","InstructedAmount":{"Currency":"GBP","Amount":"10.00"},"CreditorAccount":{"SchemeName":"SortCodeAccountNumber","Name":"Sam Morse","Identification":"11111112345678"}}},"Risk":{},"Links":{"Self":"/open-banking/v1.1/payments/f5f92b5e-5fcd-4476-86cd-9f37ea015ffe"},"Meta":{}}',
          };

          const result = validateResponse(request, response);

          assert.equal(result.isValid, true);
          assert.deepEqual(result.errors, []);
        });

        it('isValid is false, when request.Data.Initiation != response.Data.Initiation, field values', async () => {
          // change the value of EndToEndIdentification, InstructedAmount.Currency,
          // CreditorAccount.Name and CreditorAccount.Identification
          const request = {
            method: 'POST',
            path: '/open-banking/v1.1/payments',
            body: {
              Data: {
                Initiation: {
                  InstructionIdentification: '15abd6c0-18f1-4257-a765-c15ef8bf1c',
                  EndToEndIdentification: '00000000-0000-0000-0000-0000000000',
                  InstructedAmount: {
                    Currency: 'AED',
                    Amount: '10.00',
                  },
                  CreditorAccount: {
                    SchemeName: 'SortCodeAccountNumber',
                    Name: 'James Bond',
                    Identification: '007',
                  },
                },
              },
              Risk: {},
            },
          };
          const response = {
            body: '{"Data":{"PaymentId":"f5f92b5e-5fcd-4476-86cd-9f37ea015ffe","Status":"AcceptedTechnicalValidation","CreationDateTime":"2018-08-21T07:57:47+00:00","Initiation":{"InstructionIdentification":"15abd6c0-18f1-4257-a765-c15ef8bf1c","EndToEndIdentification":"8a30c4fe-a779-436f-b231-f21c05bd22","InstructedAmount":{"Currency":"GBP","Amount":"10.00"},"CreditorAccount":{"SchemeName":"SortCodeAccountNumber","Name":"Sam Morse","Identification":"11111112345678"}}},"Risk":{},"Links":{"Self":"/open-banking/v1.1/payments/f5f92b5e-5fcd-4476-86cd-9f37ea015ffe"},"Meta":{}}',
          };

          const result = validateResponse(request, response);

          assert.equal(result.isValid, false);
          assert.deepEqual(result.errors, [{
            path: ['Data', 'Initiation', 'EndToEndIdentification'],
            message: 'request.Data.Initiation.EndToEndIdentification="00000000-0000-0000-0000-0000000000" != response.Data.Initiation.EndToEndIdentification="8a30c4fe-a779-436f-b231-f21c05bd22"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Initiation', 'InstructedAmount', 'Currency'],
            message: 'request.Data.Initiation.InstructedAmount.Currency="AED" != response.Data.Initiation.InstructedAmount.Currency="GBP"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Initiation', 'CreditorAccount', 'Name'],
            message: 'request.Data.Initiation.CreditorAccount.Name="James Bond" != response.Data.Initiation.CreditorAccount.Name="Sam Morse"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Initiation', 'CreditorAccount', 'Identification'],
            message: 'request.Data.Initiation.CreditorAccount.Identification="007" != response.Data.Initiation.CreditorAccount.Identification="11111112345678"',
            code: 'OBJECTS_NOT_EQUAL',
          }]);
        });

        it('isValid is false, when request.Data.Initiation != response.Data.Initiation', async () => {
          // remove all the fields in the request.body.Data.Initiation, we should get an error
          // because the response includes it.
          const request = {
            method: 'POST',
            path: '/open-banking/v1.1/payments',
            body: {
              Data: {
                Initiation: {
                },
              },
              Risk: {},
            },
          };
          const response = {
            body: '{"Data":{"PaymentId":"f5f92b5e-5fcd-4476-86cd-9f37ea015ffe","Status":"AcceptedTechnicalValidation","CreationDateTime":"2018-08-21T07:57:47+00:00","Initiation":{"InstructionIdentification":"15abd6c0-18f1-4257-a765-c15ef8bf1c","EndToEndIdentification":"8a30c4fe-a779-436f-b231-f21c05bd22","InstructedAmount":{"Currency":"GBP","Amount":"10.00"},"CreditorAccount":{"SchemeName":"SortCodeAccountNumber","Name":"Sam Morse","Identification":"11111112345678"}}},"Risk":{},"Links":{"Self":"/open-banking/v1.1/payments/f5f92b5e-5fcd-4476-86cd-9f37ea015ffe"},"Meta":{}}',
          };

          const result = validateResponse(request, response);

          assert.equal(result.isValid, false);
          assert.deepEqual(result.errors, [{
            path: ['Data', 'Initiation', 'InstructionIdentification'],
            message: 'request.Data.Initiation.InstructionIdentification="" != response.Data.Initiation.InstructionIdentification="15abd6c0-18f1-4257-a765-c15ef8bf1c"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Initiation', 'EndToEndIdentification'],
            message: 'request.Data.Initiation.EndToEndIdentification="" != response.Data.Initiation.EndToEndIdentification="8a30c4fe-a779-436f-b231-f21c05bd22"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Initiation', 'InstructedAmount'],
            message: 'request.Data.Initiation.InstructedAmount="" != response.Data.Initiation.InstructedAmount={"Currency":"GBP","Amount":"10.00"}',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Initiation', 'CreditorAccount'],
            message: 'request.Data.Initiation.CreditorAccount="" != response.Data.Initiation.CreditorAccount={"SchemeName":"SortCodeAccountNumber","Name":"Sam Morse","Identification":"11111112345678"}',
            code: 'OBJECTS_NOT_EQUAL',
          }]);
        });

        it('isValid is false, when request.Risk != response.Risk', async () => {
          // add fields to request.body.Data.Risk, we should get an error
          // because the request includes additional fields.
          const request = {
            method: 'POST',
            path: '/open-banking/v1.1/payments',
            body: {
              Data: {
                Initiation: {
                  InstructionIdentification: '15abd6c0-18f1-4257-a765-c15ef8bf1c',
                  EndToEndIdentification: '8a30c4fe-a779-436f-b231-f21c05bd22',
                  InstructedAmount: {
                    Currency: 'GBP',
                    Amount: '10.00',
                  },
                  CreditorAccount: {
                    SchemeName: 'SortCodeAccountNumber',
                    Name: 'Sam Morse',
                    Identification: '11111112345678',
                  },
                },
              },
              Risk: {
                FakeRiskField1: 'FakeRiskField1_Value',
              },
            },
          };
          const response = {
            body: '{"Data":{"PaymentId":"f5f92b5e-5fcd-4476-86cd-9f37ea015ffe","Status":"AcceptedTechnicalValidation","CreationDateTime":"2018-08-21T07:57:47+00:00","Initiation":{"InstructionIdentification":"15abd6c0-18f1-4257-a765-c15ef8bf1c","EndToEndIdentification":"8a30c4fe-a779-436f-b231-f21c05bd22","InstructedAmount":{"Currency":"GBP","Amount":"10.00"},"CreditorAccount":{"SchemeName":"SortCodeAccountNumber","Name":"Sam Morse","Identification":"11111112345678"}}},"Risk":{},"Links":{"Self":"/open-banking/v1.1/payments/f5f92b5e-5fcd-4476-86cd-9f37ea015ffe"},"Meta":{}}',
          };

          const result = validateResponse(request, response);

          assert.equal(result.isValid, false);
          assert.deepEqual(result.errors, [{
            path: ['Risk', 'FakeRiskField1'],
            message: 'request.Risk.FakeRiskField1="FakeRiskField1_Value" != response.Risk.FakeRiskField1=""',
            code: 'OBJECTS_NOT_EQUAL',
          }]);
        });
      });
    });

    describe('accounts', () => {
      describe('/open-banking/v1.1/account-requests POST', () => {
        it('isValid is true, when request.Risk == response.Risk', async () => {
          const request = AccountRequests.post.request();
          const response = AccountRequests.post.response();

          const report = await validate(request, response, details);

          assert.deepEqual(report, { failedValidation: false });
        });

        it('isValid is true, when request and response have equal Data.Permissions, Data.ExpirationDateTime, Data.TransactionFromDateTime, Data.TransactionToDateTime, Risk field values in body', async () => {
          const request = AccountRequests.post.request();
          const response = AccountRequests.post.response();

          const report = await validate(request, response, details);

          assert.deepEqual(report, { failedValidation: false });
        });

        it('isValid is false, when request.Data.Permissions != response.Data.Permissions', async () => {
          // change request.Data.Permissions
          const request = AccountRequests.post.request();
          request.body.Data.Permissions = [
            'ReadBalances',
            'ReadDirectDebits',
            'ReadProducts',
            'ReadTransactionsCredits',
            'ReadTransactionsDebits',
          ];
          const response = AccountRequests.post.response();

          const report = await validate(request, response, details);

          assert.equal(report.failedValidation, true);
          assert.deepEqual(report.results.errors, [{
            path: ['Data', 'Permissions', 8],
            message: 'request.Data.Permissions.8="" != response.Data.Permissions.8="ReadTransactionsDetail"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Permissions', 7],
            message: 'request.Data.Permissions.7="" != response.Data.Permissions.7="ReadTransactionsDebits"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Permissions', 6],
            message: 'request.Data.Permissions.6="" != response.Data.Permissions.6="ReadTransactionsCredits"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Permissions', 5],
            message: 'request.Data.Permissions.5="" != response.Data.Permissions.5="ReadStandingOrdersDetail"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Permissions', 4],
            message: 'request.Data.Permissions.4="ReadTransactionsDebits" != response.Data.Permissions.4="ReadProducts"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Permissions', 3],
            message: 'request.Data.Permissions.3="ReadTransactionsCredits" != response.Data.Permissions.3="ReadDirectDebits"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Permissions', 2],
            message: 'request.Data.Permissions.2="ReadProducts" != response.Data.Permissions.2="ReadBeneficiariesDetail"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Permissions', 1],
            message: 'request.Data.Permissions.1="ReadDirectDebits" != response.Data.Permissions.1="ReadBalances"',
            code: 'OBJECTS_NOT_EQUAL',
          },
          {
            path: ['Data', 'Permissions', 0],
            message: 'request.Data.Permissions.0="ReadBalances" != response.Data.Permissions.0="ReadAccountsDetail"',
            code: 'OBJECTS_NOT_EQUAL',
          }]);
        });

        it('isValid is false, when request.Data.ExpirationDateTime != response.Data.ExpirationDateTime', async () => {
          // increment request.Data.ExpirationDateTime by one year
          const request = AccountRequests.post.request();
          const newDateTime = new Date(request.body.Data.ExpirationDateTime);
          newDateTime.setMonth(newDateTime.getMonth() + 12);
          request.body.Data.ExpirationDateTime = newDateTime.toISOString();

          const response = AccountRequests.post.response();

          const report = await validate(request, response, details);

          assert.equal(report.failedValidation, true);
          assert.deepEqual(report.results.errors, [{
            path: ['Data', 'ExpirationDateTime'],
            message: 'request.Data.ExpirationDateTime="2018-05-02T00:00:00.000Z" != response.Data.ExpirationDateTime="2017-05-02T00:00:00+00:00"',
            code: 'OBJECTS_NOT_EQUAL',
          }]);
        });

        it('isValid is false, when request.Data.TransactionFromDateTime != response.Data.TransactionFromDateTime', async () => {
          // increment request.Data.TransactionFromDateTime by one year
          const request = AccountRequests.post.request();
          const newDateTime = new Date(request.body.Data.TransactionFromDateTime);
          newDateTime.setMonth(newDateTime.getMonth() + 12);
          request.body.Data.TransactionFromDateTime = newDateTime.toISOString();

          const response = AccountRequests.post.response();

          const report = await validate(request, response, details);

          assert.equal(report.failedValidation, true);
          assert.deepEqual(report.results.errors, [{
            path: ['Data', 'TransactionFromDateTime'],
            message: 'request.Data.TransactionFromDateTime="2018-05-03T00:00:00.000Z" != response.Data.TransactionFromDateTime="2017-05-03T00:00:00+00:00"',
            code: 'OBJECTS_NOT_EQUAL',
          }]);
        });

        it('isValid is false, when request.Data.TransactionToDateTime != response.Data.TransactionToDateTime', async () => {
          // increment request.Data.TransactionToDateTime by one year
          const request = AccountRequests.post.request();
          const newDateTime = new Date(request.body.Data.TransactionToDateTime);
          newDateTime.setMonth(newDateTime.getMonth() + 12);
          request.body.Data.TransactionToDateTime = newDateTime.toISOString();

          const response = AccountRequests.post.response();

          const report = await validate(request, response, details);

          assert.equal(report.failedValidation, true);
          assert.deepEqual(report.results.errors, [{
            path: ['Data', 'TransactionToDateTime'],
            message: 'request.Data.TransactionToDateTime="2018-12-03T00:00:00.000Z" != response.Data.TransactionToDateTime="2017-12-03T00:00:00+00:00"',
            code: 'OBJECTS_NOT_EQUAL',
          }]);
        });

        it('isValid is false, when request.Risk != response.Risk', async () => {
          // add stuff to OBReadRequest1/Risk (OBRisk2) and OBReadResponse1/Risk (OBRisk2)
          const request = AccountRequests.post.request();
          request.body.Risk.FakeRiskField1 = 'FakeRiskField1';
          const response = AccountRequests.post.response();

          const report = await validate(request, response, details);

          assert.equal(report.failedValidation, true);
          assert.deepEqual(report.results.errors, [{
            code: 'OBJECT_ADDITIONAL_PROPERTIES',
            message: 'Additional properties not allowed: FakeRiskField1',
            path: ['Risk'],
            description: 'The Risk section is sent by the initiating party to the ASPSP. It is used to specify additional details for risk scoring for Account Info.',
          },
          {
            path: ['Risk', 'FakeRiskField1'],
            message: 'request.Risk.FakeRiskField1="FakeRiskField1" != response.Risk.FakeRiskField1=""',
            code: 'OBJECTS_NOT_EQUAL',
          }]);
        });
      });
    });

    it('if path is not "/open-banking/v1.1/account-requests" POST is it ignored', async () => {
      // add stuff to OBReadRequest1/Risk (OBRisk2) and OBReadResponse1/Risk (OBRisk2)
      // which could normally cause `failedValidation` to be true,
      // then change request from POST to GET which causes the validation not to run.
      const request = AccountRequests.post.request();
      request.method = 'GET';
      request.body.Risk.FakeRiskField1 = 'FakeRiskField1';
      const response = AccountRequests.post.response();

      const report = await validate(request, response, details);

      assert.deepEqual(report, { failedValidation: false });
    });

    it('if path is not "/open-banking/v1.1/payments" POST is it ignored', async () => {
      // change method to `GET` and remove body.Data.Initiation.CreditorAccount
      // which would normally lead to an error but since the path is
      // ignored this won't throw an error.
      const request = {
        method: 'GET',
        path: '/open-banking/v1.1/payments',
        body: {
          Data: {
            Initiation: {
              InstructionIdentification: '15abd6c0-18f1-4257-a765-c15ef8bf1c',
              EndToEndIdentification: '8a30c4fe-a779-436f-b231-f21c05bd22',
              InstructedAmount: {
                Currency: 'GBP',
                Amount: '10.00',
              },
            },
          },
          Risk: {},
        },
      };
      const response = {
        body: '{"Data":{"PaymentId":"f5f92b5e-5fcd-4476-86cd-9f37ea015ffe","Status":"AcceptedTechnicalValidation","CreationDateTime":"2018-08-21T07:57:47+00:00","Initiation":{"InstructionIdentification":"15abd6c0-18f1-4257-a765-c15ef8bf1c","EndToEndIdentification":"8a30c4fe-a779-436f-b231-f21c05bd22","InstructedAmount":{"Currency":"GBP","Amount":"10.00"},"CreditorAccount":{"SchemeName":"SortCodeAccountNumber","Name":"Sam Morse","Identification":"11111112345678"}}},"Risk":{},"Links":{"Self":"/open-banking/v1.1/payments/f5f92b5e-5fcd-4476-86cd-9f37ea015ffe"},"Meta":{}}',
      };

      const result = validateResponse(request, response);

      assert.equal(result.isValid, true);
      assert.deepEqual(result.errors, []);
    });
  });
});
