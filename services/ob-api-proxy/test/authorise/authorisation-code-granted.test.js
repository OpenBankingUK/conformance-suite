const assert = require('assert');
const proxyquire = require('proxyquire');
const httpMocks = require('node-mocks-http');
const sinon = require('sinon');

const redirectionUrl = 'http://localhost:9999/tpp/authorized';
const authorisationServerId = '123';
const accessToken = 'access-token';
const authorisationCode = '12345_67xxx';
const sessionId = 'testSession';
const validationRunId = 'testValidationRunId';
const username = 'testUser';
const scope = 'accounts';
const accountRequestId = 'testAccountRequestId';
const { mungeToken } = require('../../app/authorise/obtain-access-token');
const { base64EncodeJSON } = require('../../app/ob-util');

const exampleConfig = {
  client_id: 'testClientId',
  client_secret: 'testClientSecret',
  redirect_uri: redirectionUrl,
};

const tokenResponsePayload = {
  access_token: accessToken,
  expires_in: 3600,
};

const tokenResponse = mungeToken(tokenResponsePayload);

describe('Authorized Code Granted', () => {
  let redirection;
  let obtainAuthorisationCodeAccessTokenStub;
  let setConsentStub;
  let getUsernameStub;
  let request;
  let response;

  beforeEach(() => {
    setConsentStub = sinon.stub();
    obtainAuthorisationCodeAccessTokenStub = sinon.stub().returns(tokenResponse);
    getUsernameStub = sinon.stub().returns(username);
    redirection = proxyquire('../../app/authorise/authorisation-code-granted.js', {
      './obtain-access-token': { obtainAuthorisationCodeAccessToken: obtainAuthorisationCodeAccessTokenStub },
      './consents': { setConsent: setConsentStub },
      '../session': {
        session: {
          getUsername: getUsernameStub,
        },
      },
    });

    request = httpMocks.createRequest({
      method: 'POST',
      url: '/tpp/authorized',
      headers: {
        'authorization': sessionId,
        'x-validation-run-id': validationRunId,
        'x-config': base64EncodeJSON(exampleConfig),
      },
      body: {
        authorisationCode,
        authorisationServerId,
        scope,
        accountRequestId,
      },
    });
    response = httpMocks.createResponse();
  });

  afterEach(() => {
    obtainAuthorisationCodeAccessTokenStub.reset();
  });

  describe('redirect url configured', () => {
    it('handles the redirection route', async () => {
      await redirection.authorisationCodeGrantedHandler(request, response);
      assert.equal(response.statusCode, 204);
    });

    it('calls obtainAuthorisationCodeAccessToken to obtain an access token', async () => {
      await redirection.authorisationCodeGrantedHandler(request, response);
      assert(obtainAuthorisationCodeAccessTokenStub.calledWithExactly(
        redirectionUrl,
        authorisationCode,
        exampleConfig,
      ));
    });

    it('calls setConsent to store obtained consent', async () => {
      await redirection.authorisationCodeGrantedHandler(request, response);
      const { args } = setConsentStub.getCalls()[0];
      assert.deepEqual({
        username, authorisationServerId, scope, validationRunId,
      }, args[0]);
      assert.deepEqual(
        {
          username,
          authorisationServerId,
          scope,
          accountRequestId,
          expirationDateTime: null,
          authorisationCode,
          token: tokenResponse,
        },
        args[1],
      );
    });

    it('returns 400 status for invalid request body payload', async () => {
      request.body = {};
      await redirection.authorisationCodeGrantedHandler(request, response);
      assert.equal(response.statusCode, 400);
      // eslint-disable-next-line no-underscore-dangle
      assert.deepEqual(response._getData(), {
        message:
        'validatePayload: missingKeys=authorisationServerId, authorisationCode, scope, accountRequestId missing from request payload={}',
      });
    });

    describe('error handling', () => {
      const status = 403;
      const message = 'message';
      const error = new Error(message);
      error.status = status;

      beforeEach(() => {
        obtainAuthorisationCodeAccessTokenStub = sinon.stub().throws(error);
        redirection = proxyquire('../../app/authorise/authorisation-code-granted.js', {
          './obtain-access-token': { obtainAuthorisationCodeAccessToken: obtainAuthorisationCodeAccessTokenStub },
        });
      });

      it('relays errors including any upstream status', async () => {
        await redirection.authorisationCodeGrantedHandler(request, response);
        assert.equal(response.statusCode, status);
        // eslint-disable-next-line no-underscore-dangle
        assert.deepEqual(response._getData(), { message });
      });
    });
  });
});
