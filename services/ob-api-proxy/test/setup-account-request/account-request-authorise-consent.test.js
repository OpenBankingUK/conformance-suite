const request = require('supertest');
const assert = require('assert');
const proxyquire = require('proxyquire');
const sinon = require('sinon');
const express = require('express');
const bodyParser = require('body-parser');
const qs = require('qs');
const { DefaultPermissions } = require('../../app/setup-account-request/account-request-authorise-consent.js');
const { statePayload } = require('../../app/authorise/authorise-uri.js');
const { base64DecodeJSON } = require('../../app/ob-util');

const authorisationServerId = '123';
const sessionId = 'testSession';
const username = 'testUsername';
const accountRequestId = 'account-request-id';
const validationRunId = 'testValidationRunId';
const clientId = 'testClientId';
const clientSecret = 'testClientSecret';
const redirectUrl = 'http://example.com/redirect';
const jsonWebSignature = 'testSignedPayload';
const interactionId = 'testInteractionId';
const interactionId2 = 'testInteractionId2';
const fapiFinancialId = 'testFapiFinancialId';
const authEndpoint = 'http://example.com/auth';
const config = {
  authorization_endpoint: authEndpoint,
  client_id: clientId,
  client_secret: clientSecret,
  redirect_uri: redirectUrl,
  resource_endpoint: 'http://example.com',
};

const setupApp = (setupAccountRequestStub, setConsentStub) => {
  const createJsonWebSignatureStub = sinon.stub().returns(jsonWebSignature);
  const { generateRedirectUri } = proxyquire(
    '../../app/authorise/authorise-uri.js',
    {
      './request-jws': {
        createJsonWebSignature: createJsonWebSignatureStub,
      },
    },
  );
  const { accountRequestAuthoriseConsent } = proxyquire(
    '../../app/setup-account-request/account-request-authorise-consent.js',
    {
      './setup-account-request': {
        setupAccountRequest: setupAccountRequestStub,
      },
      '../authorise': {
        generateRedirectUri,
        setConsent: setConsentStub,
      },
      '../session': {
        extractHeaders: () => ({
          fapiFinancialId,
          interactionId,
          sessionId,
          username,
          authorisationServerId,
          validationRunId,
          config,
        }),
      },
      'uuid/v4': sinon.stub().returns(interactionId2),
    },
  );
  const app = express();
  app.use(bodyParser.json());
  app.post('/account-request-authorise-consent', accountRequestAuthoriseConsent);
  return app;
};

const doPost = app => request(app)
  .post('/account-request-authorise-consent')
  .set('authorization', sessionId)
  .set('x-authorization-server-id', authorisationServerId)
  .set('x-validation-run-id', validationRunId)
  .send();

const parseState = state => base64DecodeJSON(state);

describe('/account-request-authorise-consent with successful setupAccountRequest', () => {
  const permissions = DefaultPermissions;
  const setupAccountRequestStub = sinon.stub().returns({ accountRequestId, permissions });
  const setConsentStub = sinon.stub();
  const app = setupApp(setupAccountRequestStub, setConsentStub);

  const scope = 'openid accounts';
  const expectedStateBase64 = statePayload(
    authorisationServerId,
    sessionId,
    scope,
    interactionId2,
    accountRequestId,
  );
  const expectedParams = {
    client_id: clientId,
    redirect_uri: redirectUrl,
    request: jsonWebSignature,
    response_type: 'code',
    scope,
    state: expectedStateBase64,
  };
  const expectedState = {
    authorisationServerId,
    interactionId: interactionId2,
    scope,
    sessionId,
    accountRequestId,
  };

  it('creates a redirect URI with a 200 code via the to /authorize endpoint', (done) => {
    doPost(app)
      .end((e, r) => {
        assert.equal(r.status, 200);
        const location = r.body.uri;
        const parts = location.split('?');
        const host = parts[0];
        const params = qs.parse(parts[1]);
        assert.equal(host, authEndpoint);

        const state = parseState(params.state);
        assert.deepEqual(state, expectedState);

        assert.deepEqual(params, expectedParams);
        const header = r.headers['access-control-allow-origin'];
        assert.equal(header, '*');
        assert(setupAccountRequestStub.calledWithExactly({
          fapiFinancialId,
          interactionId,
          sessionId,
          username,
          permissions: DefaultPermissions,
          authorisationServerId,
          validationRunId,
          config,
        }));
        const keys = {
          username, authorisationServerId, scope: 'accounts', validationRunId,
        };
        const payload = { accountRequestId, permissions };
        assert(setConsentStub.calledWithExactly(keys, payload));
        done();
      });
  });
});

describe('/account-request-authorise-consent with error thrown by setupAccountRequest', () => {
  const status = 403;
  const message = 'message';
  const error = new Error(message);
  error.status = status;
  const setupAccountRequestStub = sinon.stub().throws(error);
  const setConsentStub = sinon.stub();
  const app = setupApp(setupAccountRequestStub, setConsentStub);

  it('returns status from error', (done) => {
    doPost(app)
      .end((e, r) => {
        assert.equal(r.status, status);
        assert.deepEqual(r.body, { message });
        const header = r.headers['access-control-allow-origin'];
        assert.equal(header, '*');
        assert(setupAccountRequestStub.calledWithExactly({
          fapiFinancialId,
          interactionId,
          sessionId,
          username,
          permissions: DefaultPermissions,
          authorisationServerId,
          validationRunId,
          config,
        }));
        assert.equal(setConsentStub.called, false);
        done();
      });
  });
});
