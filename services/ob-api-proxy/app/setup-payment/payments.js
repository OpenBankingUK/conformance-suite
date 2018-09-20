const request = require('superagent');
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
    log(`POST to ${paymentsUri}`);
    const call = createRequest(paymentsUri, request.post(paymentsUri), headers);
    const response = await call.send(paymentData);
    debug(`${response.status} response for ${paymentsUri}`);

    const result = await obtainResult(call, response, Object.assign({}, headers, { scope: 'payments' }));
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
