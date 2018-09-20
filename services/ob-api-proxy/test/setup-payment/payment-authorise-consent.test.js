const request = require('supertest');
const assert = require('assert');
const proxyquire = require('proxyquire');
const sinon = require('sinon');
const express = require('express');
const bodyParser = require('body-parser');
const qs = require('qs');

const { base64DecodeJSON } = require('../../app/ob-util');
const { statePayload } = require('../../app/authorise/authorise-uri.js');
const { base64EncodeJSON } = require('../../app/ob-util');

const authorisationServerId = '123';
const sessionId = 'testSessionId';
const username = 'testUsername';
const clientId = 'testClientId';
const clientSecret = 'testClientSecret';
const redirectUrl = 'http://example.com/redirect';
const jsonWebSignature = 'testSignedPayload';
const key = 'testKey';
const interactionId = key;
const fapiFinancialId = 'testFapiFinancialId';
const authEndpoint = 'http://example.com/auth';

const config = {
  api_version: '1.1',
  authorization_endpoint: authEndpoint,
  client_id: clientId,
  client_secret: clientSecret,
  redirect_uri: redirectUrl,
  resource_endpoint: 'http://example.com',
};

const setupApp = (setupPaymentStub) => {
  const createJsonWebSignatureStub = sinon.stub().returns(jsonWebSignature);
  const keyStub = sinon.stub().returns(key);
  const { generateRedirectUri } = proxyquire(
    '../../app/authorise/authorise-uri.js',
    {
      './request-jws': {
        createJsonWebSignature: createJsonWebSignatureStub,
      },
    },
  );
  const { paymentAuthoriseConsent } = proxyquire(
    '../../app/setup-payment',
    {
      './setup-payment': {
        setupPayment: setupPaymentStub,
      },
      '../authorise': {
        generateRedirectUri,
      },
      '../session': {
        extractHeaders: () => ({
          fapiFinancialId, interactionId, sessionId, username, authorisationServerId, config,
        }),
      },
      'uuid/v4': keyStub,
    },
  );
  const app = express();
  app.use(bodyParser.json());
  app.post('/payment-authorise-consent', paymentAuthoriseConsent);
  return app;
};

const doPost = app => request(app)
  .post('/payment-authorise-consent')
  .set('authorization', sessionId)
  .set('x-authorization-server-id', authorisationServerId)
  .set('x-config', base64EncodeJSON(config))
  .send();

const parseState = state => base64DecodeJSON(state);

describe('/payment-authorise-consent with successful setupPayment', () => {
  const setupPaymentStub = sinon.stub();
  const app = setupApp(setupPaymentStub);

  const scope = 'openid payments';
  const expectedStateBase64 = statePayload(authorisationServerId, sessionId, scope, interactionId);
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
    interactionId,
    scope,
    sessionId,
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
        done();
      });
  });
});

describe('/payment-authorise-consent with error thrown by setupPayment', () => {
  const status = 403;
  const message = 'message';
  const error = new Error(message);
  error.status = status;
  const setupPaymentStub = sinon.stub().throws(error);
  const app = setupApp(setupPaymentStub);

  it('returns status from error', (done) => {
    doPost(app)
      .end((e, r) => {
        assert.equal(r.status, status);
        assert.deepEqual(r.body, { message });
        const header = r.headers['access-control-allow-origin'];
        assert.equal(header, '*');
        done();
      });
  });
});
