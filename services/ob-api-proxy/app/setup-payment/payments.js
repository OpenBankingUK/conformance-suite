const superagent = require('superagent');
const { createRequest, obtainResult } = require('../ob-util');
const log = require('debug')('log');
const debug = require('debug')('debug');
const error = require('debug')('error');
const assert = require('assert');

const verifyHeaders = (headers) => {
  assert.ok(headers.idempotencyKey, 'idempotencyKey missing from headers');
  assert.ok(headers.sessionId, 'sessionId missing from headers');
};

/**
 * @description Dual purpose: payments and payment-submissions
 */
const postPayments = async (resourceServerPath, paymentPathEndpoint, headers, paymentData) => {
  try {
    verifyHeaders(headers);

    const paymentsUri = `${resourceServerPath}${paymentPathEndpoint}`;
    log('services/ob-api-proxy/app/setup-payment/payments.js:postPayments -> POST to paymentsUri=%O', paymentsUri);
    const request = createRequest(paymentsUri, superagent.post(paymentsUri), headers);
    const response = await request.send(paymentData);
    debug('services/ob-api-proxy/app/setup-payment/payments.js:postPayments -> response.status=%O, paymentsUri=%O', response.status, paymentsUri);

    const result = await obtainResult(request, response, Object.assign({}, headers, { scope: 'payments' }));
    debug('services/ob-api-proxy/app/setup-payment/payments.js:postPayments -> result=%O', result);

    return result;
  } catch (err) {
    if (err.response && err.response.text) {
      error(err.response.text);
    }

    const e = new Error(err.message);
    e.status = err.response ? err.response.status : 500;
    throw e;
  }
};

module.exports = {
  postPayments,
  verifyHeaders,
};
