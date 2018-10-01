const { setupPayment } = require('./setup-payment');
const { generateRedirectUri } = require('../authorise');
const { extractHeaders } = require('../session');
const uuidv4 = require('uuid/v4');
const error = require('debug')('error');
const debug = require('debug')('debug');
const _ = require('lodash');

exports.paymentAuthoriseConsent = async (req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  try {
    const headers = await extractHeaders(req.headers);
    const { authorisationServerId } = headers;
    const { CreditorAccount } = req.body;
    const { InstructedAmount } = req.body;
    const idempotencyKey = uuidv4();
    const response = await setupPayment(
      authorisationServerId, Object.assign({ idempotencyKey }, headers),
      CreditorAccount, InstructedAmount,
    );
    debug('services/ob-api-proxy/app/setup-payment/index.js:paymentAuthoriseConsent -> response=%j', response);

    const paymentId = _.get(response, 'Data.PaymentId');
    const validation_result = _.get(response, 'validation_result'); // eslint-disable-line

    const uri = await generateRedirectUri(
      authorisationServerId, paymentId,
      'openid payments', headers.sessionId, headers.interactionId, headers.config,
    );

    debug('services/ob-api-proxy/app/setup-payment/index.js:paymentAuthoriseConsent -> authorize uri=%j', uri);
    debug('services/ob-api-proxy/app/setup-payment/index.js:paymentAuthoriseConsent -> validation_result=%j', validation_result);
    return res
      .status(200) // We can't intercept a 302 !
      .send({
        uri,
        validation_result,
      });
  } catch (err) {
    error(err);
    const status = err.status ? err.status : 500;
    return res.status(status).send({ message: err.message });
  }
};
