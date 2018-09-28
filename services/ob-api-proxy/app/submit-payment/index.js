const { submitPayment } = require('./submit-payment');
const { consentAccessToken } = require('../authorise');
const { extractHeaders } = require('../session');
const uuidv4 = require('uuid/v4');
const error = require('debug')('error');
const debug = require('debug')('debug');
const _ = require('lodash');

exports.paymentSubmission = async (req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  try {
    const headers = await extractHeaders(req.headers);
    const { authorisationServerId, username, validationRunId } = headers;
    const idempotencyKey = uuidv4();
    const keys = {
      username, authorisationServerId, scope: 'payments', validationRunId,
    };
    const accessToken = await consentAccessToken(keys);
    const headersWithToken = Object.assign({ idempotencyKey, accessToken }, headers);
    const response = await submitPayment(authorisationServerId, headersWithToken);

    const paymentSubmissionId = _.get(response, 'Data.PaymentSubmissionId');
    const validation_result = _.get(response, 'validation_result'); // eslint-disable-line
    debug('services/ob-api-proxy/app/submit-payment/index.js:paymentSubmission -> response=%O', response);
    debug('services/ob-api-proxy/app/submit-payment/index.js:paymentSubmission -> Payment Submission succesfully completed. paymentSubmissionId=%O', paymentSubmissionId);

    return res
      .status(201) // We can't intercept a 302 !
      .send({ validation_result });
  } catch (err) {
    error(err);
    const status = err.status ? err.status : 500;
    return res.status(status).send({ message: err.message });
  }
};
