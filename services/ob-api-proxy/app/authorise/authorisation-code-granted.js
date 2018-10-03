const { setConsent } = require('./consents');
const { obtainAuthorisationCodeAccessToken } = require('./obtain-access-token');
const { session } = require('../session');
const { base64DecodeJSON } = require('../ob-util');
const debug = require('debug')('debug');
const _ = require('lodash');

const validatePayload = (payload) => {
  const requiredKeys = [
    'authorisationServerId',
    'authorisationCode',
    'scope',
    'accountRequestId',
  ];

  const missingKeys = _.filter(requiredKeys, requiredKey => !_.has(payload, requiredKey));
  if (missingKeys.length > 0) {
    const msg = `validatePayload: missingKeys=${missingKeys.join(', ')} missing from request payload=${JSON.stringify(payload)}`;

    const err = new Error(msg);
    err.status = 400;
    throw err;
  }

  return payload;
};

exports.authorisationCodeGrantedHandler = async (req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  try {
    debug(`#authorisationCodeGrantedHandler request payload: ${JSON.stringify(req.body)}`);
    const config = req.headers['x-config'] && base64DecodeJSON(req.headers['x-config']);
    const {
      authorisationServerId,
      authorisationCode,
      scope,
      accountRequestId,
    } = validatePayload(req.body);

    const tokenPayload = await obtainAuthorisationCodeAccessToken(
      config.redirect_uri,
      authorisationCode,
      config,
    );
    debug(`tokenPayload: ${JSON.stringify(tokenPayload)}`);

    const sessionId = req.headers.authorization;
    debug(`sessionId: ${sessionId}`);

    const validationRunId = req.headers['x-validation-run-id'];

    const username = await session.getUsername(sessionId);
    debug(`username: ${username}`);

    const consentPayload = {
      username,
      authorisationServerId,
      scope,
      accountRequestId,
      expirationDateTime: null,
      authorisationCode,
      token: tokenPayload,
    };
    debug(`consentPayload: ${JSON.stringify(consentPayload)}`);

    await setConsent({
      username, authorisationServerId, scope, validationRunId,
    }, consentPayload);
    return res.status(204).send();
  } catch (err) {
    debug(err);
    const status = err.status ? err.status : 500;
    return res.status(status).send({ message: err.message });
  }
};
