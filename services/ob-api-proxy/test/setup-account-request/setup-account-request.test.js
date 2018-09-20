const assert = require('assert');
const proxyquire = require('proxyquire');
const sinon = require('sinon');
const { checkErrorThrown } = require('../utils');
const { setupAccountRequest } = require('../../app/setup-account-request'); // eslint-disable-line

const authorisationServerId = 'testAuthorisationServerId';
const fapiFinancialId = 'testFinancialId';
const interactionId = 'testInteractionId';
const sessionId = 'testSessionId';
const username = 'testUsername';
const testPermissions = ['ReadAccountsDetail'];
const resourcePath = 'http://example.com';
const config = { api_version: '1.1', resource_endpoint: resourcePath };

const headers = {
  fapiFinancialId,
  interactionId,
  sessionId,
  username,
  permissions: testPermissions,
  authorisationServerId,
  config,
};

describe('setupAccountRequest called with authorisationServerId and fapiFinancialId', () => {
  const accessToken = 'access-token';
  const testAccountRequest = '88379';
  let setupAccountRequestProxy;
  let accountRequestsStub;
  const accountRequestsResponse = status => ({
    Data: {
      AccountRequestId: testAccountRequest,
      Status: status,
      Permissions: testPermissions,
    },
  });

  const setup = status => () => {
    if (status) {
      accountRequestsStub = sinon.stub().returns(accountRequestsResponse(status));
    } else {
      accountRequestsStub = sinon.stub().returns({});
    }
    setupAccountRequestProxy = proxyquire('../../app/setup-account-request/setup-account-request', {
      '../authorise': { obtainClientCredentialsAccessToken: () => accessToken },
      './account-requests': { postAccountRequests: accountRequestsStub },
    }).setupAccountRequest;
  };

  describe('when AwaitingAuthorisation', () => {
    before(setup('AwaitingAuthorisation'));

    it('returns { accountRequestId, permissions } from postAccountRequests call', async () => {
      const { accountRequestId, permissions } = await setupAccountRequestProxy(headers); // eslint-disable-line
      assert.equal(accountRequestId, testAccountRequest);
      assert.equal(permissions, testPermissions);
      const headersWithToken = Object.assign({ accessToken }, headers);
      assert(accountRequestsStub.calledWithExactly(resourcePath, headersWithToken));
    });
  });

  describe('when Authorised', () => {
    before(setup('Authorised'));

    it('returns { accountRequestId, permissions } from postAccountRequests call', async () => {
      const { accountRequestId, permissions } = await setupAccountRequestProxy(headers); // eslint-disable-line
      assert.equal(accountRequestId, testAccountRequest);
      assert.equal(permissions, testPermissions);
    });
  });

  describe('when Rejected', () => {
    before(setup('Rejected'));

    it('throws error for now', async () => {
      await checkErrorThrown(
        async () => setupAccountRequestProxy(headers),
        500, 'Account request response status: "Rejected"',
      );
    });
  });

  describe('when Revoked', () => {
    before(setup('Revoked'));

    it('throws error for now', async () => {
      await checkErrorThrown(
        async () => setupAccountRequestProxy(headers),
        500, 'Account request response status: "Revoked"',
      );
    });
  });

  describe('when AccountRequestId not found in payload', () => {
    before(setup(null));

    it('throws error', async () => {
      await checkErrorThrown(
        async () => setupAccountRequestProxy(headers),
        500, 'Account request response missing payload',
      );
    });
  });
});
