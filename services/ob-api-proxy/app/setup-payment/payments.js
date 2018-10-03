const superagent = require('superagent');
const { createRequest, obtainResult } = require('../ob-util');
const debug = require('debug')('debug');
const error = require('debug')('error');
const _ = require('lodash');

const verifyHeaders = (headers) => {
  const requiredKeys = [
    'idempotencyKey',
    'sessionId',
  ];

  const missingKeys = _.filter(requiredKeys, requiredKey => !_.has(headers, requiredKey));
  if (missingKeys.length > 0) {
    const msg = `verifyHeaders: missingKeys=${missingKeys.join(', ')} missing from headers=${JSON.stringify(headers)}`;
    throw new Error(msg);
  }
};


/**
 * @description Dual purpose: payments and payment-submissions
 */
const postPayments = async (resourceServerPath, paymentPathEndpoint, headers, paymentData) => {
  try {
    verifyHeaders(headers);

    const paymentsUri = `${resourceServerPath}${paymentPathEndpoint}`;
    debug('services/ob-api-proxy/app/setup-payment/payments.js:postPayments -> POST to paymentsUri=%j', paymentsUri);
    const request = createRequest(paymentsUri, superagent.post(paymentsUri), headers);
    const response = await request.send(paymentData);
    debug('services/ob-api-proxy/app/setup-payment/payments.js:postPayments -> response.status=%j, paymentsUri=%j', response.status, paymentsUri);

    const result = await obtainResult(request, response, Object.assign({}, headers, { scope: 'payments' }));
    debug('services/ob-api-proxy/app/setup-payment/payments.js:postPayments -> result=%j', result);

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
