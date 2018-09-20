const { setupPayment } = require('./setup-payment');
const { generateRedirectUri } = require('../authorise');
const { extractHeaders } = require('../session');
const uuidv4 = require('uuid/v4');
const error = require('debug')('error');
const debug = require('debug')('debug');

exports.paymentAuthoriseConsent = async (req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  try {
    const headers = await extractHeaders(req.headers);
    const { authorisationServerId } = headers;
    const { CreditorAccount } = req.body;
    const { InstructedAmount } = req.body;
    const idempotencyKey = uuidv4();
    const paymentId = await setupPayment(
      authorisationServerId, Object.assign({ idempotencyKey }, headers),
      CreditorAccount, InstructedAmount,
    );

    const uri = await generateRedirectUri(
      authorisationServerId, paymentId,
      'openid payments', headers.sessionId, headers.interactionId, headers.config,
    );

    debug(`authorize URL is: ${uri}`);
    return res.status(200).send({ uri }); // We can't intercept a 302 !
  } catch (err) {
    error(err);
    const status = err.status ? err.status : 500;
    return res.status(status).send({ message: err.message });
  }
};
