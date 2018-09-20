const { createClaims, createJsonWebSignature } = require('./request-jws');
const { generateRedirectUri } = require('./authorise-uri');
const {
  setConsent,
  getConsent,
  consent,
  consentAccessToken,
  consentAccessTokenAndPermissions,
  consentAccountRequestId,
  deleteConsent,
} = require('./consents');
const { authorisationCodeGrantedHandler } = require('./authorisation-code-granted');
const { obtainClientCredentialsAccessToken } = require('./obtain-access-token');

module.exports = {
  authorisationCodeGrantedHandler,
  createClaims,
  createJsonWebSignature,
  generateRedirectUri,
  obtainClientCredentialsAccessToken,
  setConsent,
  getConsent,
  consent,
  consentAccountRequestId,
  consentAccessToken,
  consentAccessTokenAndPermissions,
  deleteConsent,
};
