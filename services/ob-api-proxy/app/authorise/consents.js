const { get, set, remove } = require('../storage');
const { obtainClientCredentialsAccessToken } = require('./obtain-access-token');
const { getAccountRequest } = require('../setup-account-request/account-requests');
const uuidv4 = require('uuid/v4');
const debug = require('debug')('debug');
const _ = require('lodash');

const AUTH_SERVER_USER_CONSENTS_COLLECTION = 'authorisationServerUserConsents';

const validateCompositeKey = (obj) => {
  const requiredKeys = [
    'username',
    'authorisationServerId',
    'scope',
    'validationRunId',
  ];

  const missingKeys = _.filter(requiredKeys, requiredKey => !_.has(obj, requiredKey));
  if (missingKeys.length > 0) {
    const msg = `validateCompositeKey: missingKeys=${missingKeys.join(', ')} from consent lookup keys obj=${JSON.stringify(obj)}`;
    throw new Error(msg);
  }
};

const generateCompositeKey = (obj) => {
  validateCompositeKey(obj);
  return `${obj.username}:::${obj.authorisationServerId}:::${obj.scope}:::${obj.validationRunId}`;
};

const deleteConsent = async (keys) => {
  debug(`#deleteConsent keys: [${JSON.stringify(keys)}]`);
  const compositeKey = generateCompositeKey(keys);
  await remove(AUTH_SERVER_USER_CONSENTS_COLLECTION, compositeKey);
};

const consentPayload = async compositeKey =>
  get(AUTH_SERVER_USER_CONSENTS_COLLECTION, compositeKey);

const getConsent = async (keys) => {
  const compositeKey = generateCompositeKey(keys);
  debug(`consent#id (compositeKey): ${compositeKey}`);
  return consentPayload(compositeKey);
};

const consent = async (keys) => {
  const payload = await getConsent(keys);
  debug(`consent#payload: ${JSON.stringify(payload)}`);
  if (!payload) {
    const err = new Error(`User [${keys.username}] has not yet given consent to access their ${keys.scope}`);
    err.status = 500;
    throw err;
  }
  return payload;
};

const setConsent = async (keys, payload) => {
  debug(`#setConsent keys: [${JSON.stringify(keys)}]`);
  debug(`#setConsent payload: [${JSON.stringify(payload)}]`);
  const stored = await getConsent(keys);
  const toStore = Object.assign(
    {},
    payload,
    stored && stored.accountRequestId === payload.accountRequestId
      ? { permissions: stored.permissions }
      : {},
  );
  const compositeKey = generateCompositeKey(keys);
  debug(`#setConsent compositeKey: [${compositeKey}]`);
  return set(AUTH_SERVER_USER_CONSENTS_COLLECTION, toStore, compositeKey);
};

const consentAccessToken = async (keys) => {
  const existing = await consent(keys);
  return existing.token.accessToken;
};

const consentAccessTokenAndPermissions = async (keys) => {
  const existing = await consent(keys);
  const { accessToken } = existing.token;
  const { permissions } = existing;

  return { accessToken, permissions };
};

const getConsentStatus = async (accountRequestId, authorisationServerId, sessionId, config) => {
  const accessToken = await obtainClientCredentialsAccessToken(config);
  debug(`getConsentStatus#accessToken: ${accessToken}`);

  const fapiFinancialId = config.fapi_financial_id;
  debug(`getConsentStatus#fapiFinancialId: ${fapiFinancialId}`);
  const interactionId = uuidv4();
  const headers = {
    accessToken, fapiFinancialId, interactionId, sessionId, authorisationServerId,
  };
  const response = await getAccountRequest(accountRequestId, config.resource_endpoint, headers);
  debug(`getConsentStatus#getAccountRequest: ${JSON.stringify(response)}`);

  if (!response || !response.Data) {
    const err = new Error(`Bad account request response: "${JSON.stringify(response)}"`);
    err.status = 500;
    throw err;
  }
  const result = response.Data.Status;
  debug(`getConsentStatus#Status: ${result}`);
  return result;
};

const consentAccountRequestId = async (keys) => {
  const existing = await consent(keys);
  return existing.accountRequestId;
};

module.exports = {
  generateCompositeKey,
  setConsent,
  getConsent,
  consent,
  consentAccessToken,
  consentAccessTokenAndPermissions,
  getConsentStatus,
  consentAccountRequestId,
  deleteConsent,
  AUTH_SERVER_USER_CONSENTS_COLLECTION,
};
