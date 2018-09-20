const {
  deleteAccountRequest,
  postAccountRequests,
  getAccountRequest,
  buildAccountRequestData,
} = require('../../app/setup-account-request/account-requests');
const assert = require('assert');

const nock = require('nock');

const permissions = ['ReadAccountsDetail'];
const requestBody = buildAccountRequestData(permissions);
const accountRequestId = '88379';
const response = {
  Data: {
    AccountRequestId: accountRequestId,
    Status: 'AwaitingAuthentication',
    CreationDateTime: (new Date()).toISOString(),
    Permissions: requestBody.Data.Permissions,
  },
  Risk: {},
  Links: {
    self: `/account-requests/${accountRequestId}`,
  },
  Meta: {
    'total-pages': 1,
  },
};

const accessToken = '2YotnFZFEjr1zCsicMWpAA';
const fapiFinancialId = 'abc';
const interactionId = 'xyz';
const sessionId = 'testSessionId';
const apiVersion = '2.0';
const config = {
  api_version: apiVersion,
  resource_endpoint: 'http://example.com',
  transport_cert: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
  transport_key: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
};
const headers = {
  accessToken, fapiFinancialId, interactionId, sessionId, permissions, config,
};

describe('postAccountRequests', () => {
  nock(/example\.com/)
    .post(`/prefix/open-banking/v${apiVersion}/account-requests`)
    .matchHeader('authorization', `Bearer ${accessToken}`) // required
    .matchHeader('x-fapi-financial-id', fapiFinancialId) // required
    // optional x-jws-signature
    // optional x-fapi-customer-last-logged-time
    // optional x-fapi-customer-ip-address
    // optional x-fapi-interaction-id
    .reply(201, response);

  it('returns data when 201 OK', async () => {
    const resourceServerPath = 'http://example.com/prefix';
    const result = await postAccountRequests(resourceServerPath, headers);
    result.Data.CreationDateTime = response.Data.CreationDateTime;
    assert.deepEqual(result, response);
  });
});

describe('getAccountRequest', () => {
  nock(/example\.com/)
    .get(`/prefix/open-banking/v1.1/account-requests/${accountRequestId}`)
    .matchHeader('authorization', `Bearer ${accessToken}`) // required
    .matchHeader('x-fapi-financial-id', fapiFinancialId) // required
    // optional x-jws-signature
    // optional x-fapi-customer-last-logged-time
    // optional x-fapi-customer-ip-address
    // optional x-fapi-interaction-id
    .reply(200, response);

  it('returns data when 200 OK', async () => {
    const resourceServerPath = 'http://example.com/prefix';
    const result = await getAccountRequest(accountRequestId, resourceServerPath, headers);
    result.Data.CreationDateTime = response.Data.CreationDateTime;
    assert.deepEqual(result, response);
  });
});

describe('deleteAccountRequest', () => {
  nock(/example\.com/)
    .delete(`/prefix/open-banking/v${apiVersion}/account-requests/${accountRequestId}`)
    .matchHeader('authorization', `Bearer ${accessToken}`) // required
    .matchHeader('x-fapi-financial-id', fapiFinancialId) // required
    .matchHeader('x-fapi-interaction-id', interactionId) // required
    // optional x-jws-signature
    // optional x-fapi-customer-last-logged-time
    // optional x-fapi-customer-ip-address
    .reply(204);

  it('returns true when 204 No Content', async () => {
    const resourceServerPath = 'http://example.com/prefix';
    const result = await deleteAccountRequest(accountRequestId, resourceServerPath, headers);
    assert.deepEqual(result, true);
  });
});
