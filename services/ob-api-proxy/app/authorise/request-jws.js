/* eslint camelcase: 0 */
const jws = require('jws');
const jose = require('node-jose');

/**
 * Issuer of the token.
 * OpenID Connect protocol mandates this MUST include the client ID of the TPP.
 * Should contain the ClientID of the TPP’s OAuth Client.
 * Required as per FAPI RW / OpenID Standard.
 * For now return raw clientId
 */
const issuer = clientId => clientId;

/**
 * Used to help mitigate against replay attacks.
 * Required by FAPI Read Write (Hybrid explicitly required –
 * required by OIDC Core for Hybrid Flow).
 * Hybrid Flow support is optional in the OB Security Profile.
 * For now return dummy value.
 */
const nonce = () => 'dummy-nonce';

const claims = requestId => ({
  userinfo: {
    openbanking_intent_id: {
      value: requestId, // not sure this
      essential: true,
    },
  },
  id_token: {
    openbanking_intent_id: {
      value: requestId,
      essential: true,
    },
    acr: {
      essential: true,
    },
  },
});

const createClaims = (
  scope, requestId, clientId, authServerIssuer,
  registeredRedirectUrl, state, useOpenidConnect, // eslint-disable-line
) => ({
  aud: authServerIssuer,
  iss: issuer(clientId),
  response_type: 'code',
  client_id: clientId,
  redirect_uri: registeredRedirectUrl,
  scope,
  state,
  nonce: nonce(),
  max_age: 86400,
  claims: claims(requestId),
});

const signWithNone = payload => jws.sign({
  header: { alg: 'none' },
  payload,
});

const signWithFapiAlg = async (alg, payload, key, kid) => {
  const privateSigningKey = await jose.JWK.asKey(key, 'pem');
  const signatureConfig = {
    fields: { alg, kid },
    format: 'compact',
  };

  const result = await jose.JWS.createSign(signatureConfig, privateSigningKey)
    .update(JSON.stringify(payload), 'utf-8')
    .final();

  return result;
};

const signWithOtherAlg = (alg, payload, key) => jws.sign({
  header: { alg },
  payload,
  privateKey: key,
});

/**
 * FAPI says JWS signatures shall use the PS256 or ES256 algorithms for signing.
 * See: https://openid.net/specs/openid-financial-api-part-2.html#rfc.section.8.6
 *
 * As the Open Banking certificates are RSA certificates, only PS256 will be supported.
 * Open Banking is also permitting RS256 in MIT for now to ease implementation.
 * Some reference banks are also permitting HS256 for now to ease implementation.
 */
const createJsonWebSignature = (
  payload, signingAlgs,
  { client_secret, signing_key, signing_kid },
) => {
  switch (true) {
    case signing_key && signingAlgs.includes('RS256'):
      return signWithOtherAlg('RS256', payload, signing_key);
    case signing_key && signingAlgs.includes('PS256'):
      return signWithFapiAlg('PS256', payload, signing_key, signing_kid);
    case signingAlgs.includes('HS256'):
      return signWithOtherAlg('HS256', payload, client_secret);
    case signingAlgs.includes('none'):
      return signWithNone(payload);
    case !signing_key && (signingAlgs.includes('RS256') || signingAlgs.includes('RS256')):
      throw new Error(`cannot create JSON web signature for ${signingAlgs} as signing key was undefined or blank`);
    default:
      throw new Error(`cannot create JSON web signature as ${signingAlgs} not supported`);
  }
};

module.exports = {
  signWithNone,
  createClaims,
  createJsonWebSignature,
};
