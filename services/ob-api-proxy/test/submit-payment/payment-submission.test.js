const request = require('supertest');
const assert = require('assert');
const proxyquire = require('proxyquire');
const sinon = require('sinon');
const express = require('express');
const bodyParser = require('body-parser');
const { base64EncodeJSON } = require('../../app/ob-util');

const fapiFinancialId = 'testFapiFinancialId';
const authorisationServerId = 'testAuthServerId';
const interactionId = 'testInteractionId';
const sessionId = 'testSessionId';
const username = 'testUsername';

const config = {
  api_version: '1.1',
  client_id: 'clientId',
  client_secret: 'clientSecret',
  resource_endpoint: 'http://example.com',
};

const setupApp = (submitPaymentStub, consentAccessTokenStub) => {
  const { paymentSubmission } = proxyquire(
    '../../app/submit-payment',
    {
      './submit-payment': {
        submitPayment: submitPaymentStub,
      },
      '../authorise': {
        consentAccessToken: consentAccessTokenStub,
      },
      '../session': {
        extractHeaders: () => ({
          fapiFinancialId, interactionId, sessionId, username,
        }),
      },
    },
  );
  const app = express();
  app.use(bodyParser.json());
  app.post('/payment-submissions', paymentSubmission);
  return app;
};

const PAYMENT_SUBMISSION_ID = 'PS456';
const accessToken = 'testAccessToken';

const doPost = app => request(app)
  .post('/payment-submissions')
  .set('x-authorization-server-id', authorisationServerId)
  .set('x-fapi-interaction-id', interactionId)
  .set('x-config', base64EncodeJSON(config))
  .send();

describe('/payment-submission with successful submitPayment', () => {
  const submitPaymentStub = sinon.stub().returns(PAYMENT_SUBMISSION_ID);
  const consentAccessTokenStub = sinon.stub().returns(accessToken);
  const app = setupApp(submitPaymentStub, consentAccessTokenStub);

  it('make payment submission and returns paymentSubmissionId', (done) => {
    doPost(app)
      .end((e, r) => {
        assert.equal(r.status, 201);

        const header = r.headers['access-control-allow-origin'];
        assert.equal(header, '*');
        done();
      });
  });
});

describe('/payment-submit with error thrown by submitPayment', () => {
  const status = 403;
  const message = 'message';
  const error = new Error(message);
  error.status = status;
  const submitPaymentStub = sinon.stub().throws(error);
  const consentAccessTokenStub = sinon.stub().returns(accessToken);
  const app = setupApp(submitPaymentStub, consentAccessTokenStub);

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
