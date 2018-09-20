const { setConsent } = require('./consents');
const { obtainAuthorisationCodeAccessToken } = require('./obtain-access-token');
const { session } = require('../session');
const { base64DecodeJSON } = require('../ob-util');
const debug = require('debug')('debug');

const validatePayload = (payload) => {
  const {
    authorisationServerId,
    authorisationCode,
    scope,
    accountRequestId,
  } = payload;
  const missing = [];
  if (!authorisationServerId) {
    missing.push('authorisationServerId');
  }
  if (!authorisationCode) {
    missing.push('authorisationCode');
  }
  if (!scope) {
    missing.push('scope');
  }
  if (!accountRequestId) {
    missing.push('accountRequestId');
  }
  if (missing.length > 0) {
    const err = new Error(`Bad request, ${missing.join(', ')} missing from request payload`);
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
