/* eslint camelcase: 0 */
const request = require('superagent');
const { setupMutualTLS, createJwt, createBasicAuth } = require('../ob-util');
const debug = require('debug')('debug');

const mungeToken = (body) => {
  const accessToken = body.access_token;
  const tokenType = body.token_type;
  const tokenExpiresAt = body.expires_in ?
    new Date().getTime() + (parseInt(body.expires_in, 10) * 1000) : null;
  return { accessToken, tokenType, tokenExpiresAt };
};

const postToken = async (
  {
    client_id, client_secret,
    signing_key,
    token_endpoint, token_endpoint_auth_method,
    transport_cert, transport_key,
  }, payload) => {
  try {
    const isJwt = token_endpoint_auth_method === 'private_key_jwt';
    const isSecretBasic = token_endpoint_auth_method === 'client_secret_basic';
    const mtlsRequest =
      setupMutualTLS(token_endpoint, request.post(token_endpoint), transport_cert, transport_key).type('form');
    const body = Object.assign({}, payload);

    if (isJwt) {
      // https://openid.net/specs/openid-connect-core-1_0.html#IDToken
      const claims = {
        iss: client_id,
        sub: client_id,
        aud: token_endpoint,
      };
      const createdJwt = createJwt(claims, signing_key);
      body.client_assertion_type = 'urn:ietf:params:oauth:client-assertion-type:jwt-bearer';
      body.client_assertion = createdJwt;
    } else if (isSecretBasic) {
      mtlsRequest.set('authorization', createBasicAuth(client_id, client_secret));
    }

    const response = await mtlsRequest.send(body);
    const data = mungeToken(response.body);

    debug(`accessToken: ${JSON.stringify(data)}`);

    return data;
  } catch (err) {
    const errMsg = `${err.message} ${err.response ? JSON.stringify(err.response.body) : ''}`;
    const e = new Error(errMsg);
    e.status = err.response ? err.response.status : 500;
    throw e;
  }
};

/*
  Client Credentials Grant
  https://tools.ietf.org/html/rfc6749#section-4.4

   +---------+                                  +---------------+
   |         |                                  |               |
   |         |>--(A)- Client Authentication --->| Authorization |
   | Client  |                                  |     Server    |
   |         |<--(B)---- Access Token ---------<|               |
   |         |                                  |               |
   +---------+                                  +---------------+

  (A)  The client authenticates with the authorization server and
      requests an access token from the token endpoint.

  (B)  The authorization server authenticates the client, and if valid,
      issues an access token.

  Access Token Request
  https://tools.ietf.org/html/rfc6749#section-4.4.2

  grant_type - REQUIRED.  Value MUST be set to "client_credentials".

  scope - OPTIONAL.  The scope of the access request.
 */
const obtainClientCredentialsAccessToken = async (config) => {
  const accessTokenRequest = {
    scope: 'accounts payments', // include both scopes for client credentials grant
    grant_type: 'client_credentials',
  };

  const { accessToken } = await postToken(config, accessTokenRequest);

  return accessToken;
};

/*
  Authorization Code Grant
  https://tools.ietf.org/html/rfc6749#section-4.1

  Headless Flow (deviates from rfc6749):


  +----|-----+          Client Identifier      +---------------+
  |         -+----(A)-- & Redirection URI ---->|               |
  |  Elixir  |                                 | Authorization |
  |  app    -+                                 |     Server    |
  |          |                                 |               |
  |         -+----(B)-- Authorization Code ---<|               |
  +-|----|---+                                 +---------------+
    |    |                                         ^      v
   (A)  (B)                                        |      |
    |    |                                         |      |
    ^    v                                         |      |
  +---------+                                      |      |
  |         |>---(C)-- Authorization Code ---------'      |
  |  Node   |          & Redirection URI                  |
  |  app    |                                             |
  |         |<---(D)----- Access Token -------------------'
  +---------+       (w/ Optional Refresh Token)

  Access Token Request
  https://tools.ietf.org/html/rfc6749#section-4.1.3

  grant_type - REQUIRED.  Value MUST be set to "authorization_code".

  code - REQUIRED.  The authorization code received from the
        authorization server.

  redirect_uri - REQUIRED, if the "redirect_uri" parameter was included in the
        authorization request, and their values MUST be identical.

  client_id - REQUIRED, if the client is not authenticating with the
        authorization server.
*/
const obtainAuthorisationCodeAccessToken = async (
  redirectionUrl, authorisationCode, config) => {
  const accessTokenRequest = {
    grant_type: 'authorization_code',
    redirect_uri: redirectionUrl,
    code: authorisationCode,
    client_id: config.client_id,
  };

  return postToken(config, accessTokenRequest);
};

module.exports = {
  mungeToken,
  obtainClientCredentialsAccessToken,
  obtainAuthorisationCodeAccessToken,
};
