/* eslint camelcase: 0 */
const { createClaims, signWithNone, createJsonWebSignature } = require('./request-jws');
const { base64EncodeJSON, isMock } = require('../ob-util');
const qs = require('qs');

const statePayload = (authorisationServerId, sessionId, scope, interactionId, accountRequestId) => {
  const state = {
    authorisationServerId,
    interactionId,
    sessionId,
    scope,
    accountRequestId,
  };
  return base64EncodeJSON(state);
};

const generateRedirectUri = async (authorisationServerId, requestId, scope,
  sessionId, interactionId, config) => {
  const {
    authorization_endpoint, client_id, issuer, redirect_uri,
  } = config;
  const state = statePayload(authorisationServerId, sessionId, scope, interactionId, requestId);
  const signingAlgs = ['RS256'];
  const payload = createClaims(
    scope, requestId, client_id, issuer,
    redirect_uri, state, createClaims,
  );
  const signature = isMock(config.resource_endpoint) ? signWithNone(payload)
    : await createJsonWebSignature(payload, signingAlgs, config);
  const uri =
    `${authorization_endpoint}?${qs.stringify({
      redirect_uri,
      state,
      client_id,
      response_type: 'code',
      request: signature,
      scope,
    })}`;
  return uri;
};

module.exports = {
  generateRedirectUri,
  statePayload,
};
