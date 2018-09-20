const assert = require('assert');
const proxyquire = require('proxyquire');
const sinon = require('sinon');
const { base64EncodeJSON } = require('../../app/ob-util');

const authorisationServerId = 'testAuthorisationServerId';
const sessionId = 'testSessionId';
const username = 'testUsername';
const generatedInteractionId = 'testInteractionId';
const fapiFinancialId = 'aaax5nTR33811Qy';
const validationRunId = 'testValidationRunId';
const permissions = 'ReadAccountsDetail ReadTransactionsDebits';
const permissionsList = permissions.split(' ');

const { extractHeaders } = proxyquire(
  '../../app/session/request-headers.js',
  {
    './session': {
      getUsername: async () => username,
    },
    'uuid/v4': sinon.stub().returns(generatedInteractionId),
  },
);

const ACCOUNT_SWAGGER = 'https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/ee715e094a59b37aeec46aef278f528f5d89eb03/dist/v1.1/account-info-swagger.json';

const exampleConfig = {
  authorization_endpoint: 'https://aspsp.example.com/auth',
  authorization_server_id: 'aaaj4NmBD8lQxmLh2O',
  client_id: 'testClientId',
  client_secret: 'testClientSecret',
  fapi_financial_id: fapiFinancialId,
  id_token_signed_response_alg: 'RS256',
  issuer: 'http://aspsp.example.com', // openid config issuer
  openid_config_uri: 'https://aspsp.example.com/.well-known/openid-configuration',
  redirect_uri: 'https://tpp.example.com/oauth2/callback',
  request_object_signing_alg: 'RS256',
  resource_endpoint: 'https://aspsp.example.com', // without /open-banking/v*.*
  response_type: 'code id_token',
  signing_key: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
  signing_kid: 'XXXXXX-XXXXxxxXxXXXxxx_xxxx',
  software_statement_id: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx',
  token_endpoint: 'http: //example.com/token',
  token_endpoint_auth_method: 'private_key_jwt',
  token_endpoint_auth_signing_alg: 'RS256',
  transport_cert: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
  transport_key: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
};

const requestHeaders = {
  'authorization': sessionId,
  'x-authorization-server-id': authorisationServerId,
  'x-validation-run-id': validationRunId,
  'x-permissions': permissions,
  'x-swagger-uris': ACCOUNT_SWAGGER,
  'x-config': base64EncodeJSON(exampleConfig),
};

const expectedHeaders = extra => Object.assign({}, {
  fapiFinancialId,
  sessionId,
  username,
  authorisationServerId,
  validationRunId,
  swaggerUris: [ACCOUNT_SWAGGER],
  permissions: permissionsList,
  config: exampleConfig,
}, extra);

describe('extractHeaders from request headers', () => {
  it('returns headers object', async () => {
    const interactionId = generatedInteractionId;
    const headers = await extractHeaders(requestHeaders);
    assert.deepEqual(headers, expectedHeaders({ interactionId }));
  });

  describe('when x-fapi-interaction-id in headers', () => {
    it('returns headers with same interactionId', async () => {
      const interactionId = 'existingId';
      const headers = await extractHeaders(Object.assign({ 'x-fapi-interaction-id': interactionId }, requestHeaders));
      assert.deepEqual(headers, expectedHeaders({ interactionId }));
    });
  });

  describe('when x-swagger-uris in headers', () => {
    it('returns headers with same x-swagger-uris', async () => {
      const interactionId = generatedInteractionId;

      const BASIC_ACCOUNT_SWAGGER = 'https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v2.0.0/account-info-swagger.json';
      const DETAIL_ACCOUNT_SWAGGER = 'https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v2.0.0/account-info-swagger-detail.json';
      const xswaggerUris = `${BASIC_ACCOUNT_SWAGGER} ${DETAIL_ACCOUNT_SWAGGER}`;
      const xswaggerUrisExpected = xswaggerUris.split(' ');

      const headers = Object.assign({}, requestHeaders, { 'x-swagger-uris': xswaggerUris });
      const extracted = await extractHeaders(headers);

      assert.deepEqual(
        extracted,
        expectedHeaders({ swaggerUris: xswaggerUrisExpected, interactionId }),
      );
    });
  });
});
