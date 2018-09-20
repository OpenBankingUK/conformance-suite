const request = require('superagent');
const { createRequest, obtainResult } = require('../ob-util');
const log = require('debug')('log');
const errorLog = require('debug')('error');
const debug = require('debug')('debug');
const assert = require('assert');
const util = require('util');

const buildAccountRequestData = Permissions => ({
  Data: { Permissions },
  Risk: {},
});

const verifyHeaders = (headers) => {
  assert.ok(headers.sessionId, 'sessionId missing from headers');
  if (headers.config) {
    assert.ok(headers.config.api_version, 'api_version missing from headers.config');
  }
};

/*
 * For now only support Client Credentials Grant Type (OAuth 2.0).
 * @resourceServerPath e.g. http://example.com/open-banking/v1.1
 */
const postAccountRequests = async (resourceServerPath, headers) => {
  try {
    verifyHeaders(headers);
    const body = buildAccountRequestData(headers.permissions);
    const apiVersion = headers.config.api_version;
    const accountRequestsUri = `${resourceServerPath}/open-banking/v${apiVersion}/account-requests`;
    log(`POST to ${accountRequestsUri}`);
    const call = createRequest(accountRequestsUri, request.post(accountRequestsUri), headers);
    const response = await call.send(body);
    debug(`${response.status} response for ${accountRequestsUri}`);

    const result = await obtainResult(call, response, Object.assign({}, headers, { scope: 'accounts' }));
    return result;
  } catch (err) {
    errorLog(util.inspect(err));
    const error = new Error(err.message);
    error.status = err.response ? err.response.status : 500;
    throw error;
  }
};

/*
 * For now only support Client Credentials Grant Type (OAuth 2.0).
 * @resourceServerPath e.g. http://example.com/open-banking/v1.1
 */
const getAccountRequest = async (accountRequestId, resourceServerPath, headers) => {
  try {
    verifyHeaders(headers);
    const accountRequestsUri = `${resourceServerPath}/open-banking/v1.1/account-requests/${accountRequestId}`;
    log(`GET to ${accountRequestsUri}`);
    const response = await createRequest(
      accountRequestsUri,
      request.get(accountRequestsUri), headers,
    ).send();
    debug(`${response.status} response for ${accountRequestsUri}`);
    return response.body;
  } catch (err) {
    errorLog(util.inspect(err));
    const error = new Error(err.message);
    error.status = err.response ? err.response.status : 500;
    throw error;
  }
};

const deleteAccountRequest = async (accountRequestId, resourceServerPath, headers) => {
  try {
    verifyHeaders(headers);
    const apiVersion = headers.config.api_version;
    const accountRequestDeleteUri = `${resourceServerPath}/open-banking/v${apiVersion}/account-requests/${accountRequestId}`;
    log(`DELETE to ${accountRequestDeleteUri}`);
    const response = await createRequest(
      accountRequestDeleteUri,
      request.del(accountRequestDeleteUri), headers,
    ).send();
    debug(`${response.status} response for ${accountRequestDeleteUri}`);
    if (response.status === 204) {
      return true;
    }
    errorLog(`deleteAccountRequest, expected 204, got: ${util.inspect(response)}`);
    throw new Error(`Expected 204 response to delete account request, got: ${response.status} body: ${response.body}`);
  } catch (err) {
    errorLog(util.inspect(err));
    const error = new Error(err.message);
    error.status = err.response ? err.response.status : 400;
    throw error;
  }
};

module.exports = {
  buildAccountRequestData,
  postAccountRequests,
  getAccountRequest,
  deleteAccountRequest,
};
