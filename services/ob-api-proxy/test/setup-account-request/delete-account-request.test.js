const assert = require('assert');
const proxyquire = require('proxyquire');
const sinon = require('sinon');
const { checkErrorThrown } = require('../utils');

const authorisationServerId = 'testAuthorisationServerId';
const fapiFinancialId = 'testFinancialId';
const interactionId = 'testInteractionId';
const sessionId = 'testSessionId';
const username = 'testUsername';
const validationRunId = 'testRunId';
const resourcePath = 'http://example.com';
const config = {
  api_version: '1.1',
  client_id: 'testClientId',
  client_secret: undefined,
  resource_endpoint: resourcePath,
};
const headers = {
  fapiFinancialId,
  interactionId,
  sessionId,
  username,
  authorisationServerId,
  validationRunId,
  config,
};

describe('deleteAccountRequest called with authorisationServerId and fapiFinancialId', () => {
  const accessToken = 'access-token';
  const accountRequestId = '88379';
  let deleteRequestProxy;
  let deleteAccountRequestStub;
  let consentAccountRequestIdStub;

  const setup = success => () => {
    deleteAccountRequestStub = sinon.stub().returns(success);
    consentAccountRequestIdStub = sinon.stub().returns(accountRequestId);
    deleteRequestProxy = proxyquire(
      '../../app/setup-account-request/delete-account-request',
      {
        '../authorise': {
          obtainClientCredentialsAccessToken: () => accessToken,
          consentAccountRequestId: consentAccountRequestIdStub,
        },
        './account-requests': { deleteAccountRequest: deleteAccountRequestStub },
      },
    ).deleteRequest;
  };

  describe('when delete successful', () => {
    before(setup(true));

    it('returns 204 from deleteRequests call', async () => {
      const status = await deleteRequestProxy(headers);
      assert.equal(status, 204);
      const headersWithToken = {
        accessToken,
        fapiFinancialId,
        interactionId,
        sessionId,
        username,
        authorisationServerId,
        validationRunId,
        config,
      };
      assert(deleteAccountRequestStub.calledWithExactly(
        accountRequestId,
        resourcePath,
        headersWithToken,
      ));
    });
  });

  describe('when delete not successful due to missing account request ID', () => {
    before(setup(false));

    it('throws error for now', async () => {
      await checkErrorThrown(
        async () => deleteRequestProxy(headers),
        400, 'Bad request - account request ID not found',
      );
    });
  });
});
