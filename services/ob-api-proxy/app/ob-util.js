/* eslint camelcase: 0 */
const nJwt = require('njwt');
const _ = require('lodash');
const { validate, validateResponseOn } = require('./validator');

process.env.NODE_TLS_REJECT_UNAUTHORIZED = 0; // To enable use of self signed certs

const APPLICATION_JSON = 'application/json; charset=utf-8';

const BASE64 = 'base64';
const base64Encode = string =>
  Buffer.from(string).toString(BASE64);

const base64Decode = encoded =>
  Buffer.from(encoded, BASE64).toString();

const base64EncodeJSON = object =>
  base64Encode(JSON.stringify(object));

const base64DecodeJSON = encoded =>
  JSON.parse(base64Decode(encoded));

const ca = base64Decode(process.env.OB_ISSUING_CA || '');
const isMock = uri => uri.includes('localhost') || uri.includes('reference-mock-server');
const setupMutualTLS = (uri, agent, cert, key) => {
  if (isMock(uri)) return agent;
  return agent.key(key).cert(cert).ca(cert);
};

const setOptionalHeader = (header, value, requestObj) => {
  if (value) requestObj.set(header, value);
};

const setHeaders = (requestObj, headers) => {
  requestObj
    .set('authorization', `Bearer ${headers.accessToken}`)
    .set('content-type', APPLICATION_JSON)
    .set('accept', APPLICATION_JSON)
    .set('x-fapi-interaction-id', headers.interactionId)
    .set('x-fapi-financial-id', headers.fapiFinancialId);

  setOptionalHeader('x-idempotency-key', headers.idempotencyKey, requestObj);
  setOptionalHeader('x-fapi-customer-last-logged-time', headers.customerLastLogged, requestObj);
  setOptionalHeader('x-fapi-customer-ip-address', headers.customerIp, requestObj);
  setOptionalHeader('x-jws-signature', headers.jwsSignature, requestObj);
  setOptionalHeader('x-validation-run-id', headers.validationRunId, requestObj);
  return requestObj;
};

const verifyHeaders = (headers) => {
  const requiredKeys = [
    'accessToken',
    'fapiFinancialId',
    'interactionId',
    'config',
    'config.transport_cert',
    'config.transport_key',
  ];
  const missingKeys = _.filter(requiredKeys, requiredKey => !_.has(headers, requiredKey));
  if (missingKeys.length > 0) {
    const msg = `verifyHeaders: Missing: ${missingKeys.join(', ')} missing from headers`;
    throw new Error(msg);
  }
};

const createRequest = (uri, requestObj, headers) => {
  verifyHeaders(headers);
  const { transport_cert, transport_key } = headers.config;
  const req = setHeaders(setupMutualTLS(uri, requestObj, transport_cert, transport_key), headers);
  return req;
};

const validateRequestResponse = async (req, res, responseBody, details) => {
  const { failedValidation } = await validate(req, res, details);
  return Object.assign(responseBody, { failedValidation });
};

const obtainResult = async (call, response, headers) => {
  let result;
  if (validateResponseOn()) {
    result =
      await validateRequestResponse(call, response.res, response.body, {
        interactionId: headers.interactionId,
        sessionId: headers.sessionId,
        permissions: headers.permissions,
        authorisationServerId: headers.authorisationServerId,
        validationRunId: headers.validationRunId,
        scope: headers.scope,
        swaggerUris: headers.swaggerUris,
      });
  } else {
    result = response.body;
  }
  return result;
};

const createJwt = (claims, signingKey) => {
  const createdJwt = nJwt.create(claims, signingKey, 'RS256');
  // createdJwt.setHeader('kid', <signingKid>); // todo: check whether this is needed
  return createdJwt.compact();
};

// Basic Authentication Scheme: https://tools.ietf.org/html/rfc2617#section-2
const createBasicAuth = (userid, password) => {
  const basicCredentials = base64Encode(`${userid}:${password}`);
  return `Basic ${basicCredentials}`;
};

module.exports = {
  base64EncodeJSON,
  base64Decode,
  base64DecodeJSON,
  setupMutualTLS,
  createBasicAuth,
  createRequest,
  obtainResult,
  caCert: ca,
  validateRequestResponse,
  createJwt,
  isMock,
};
