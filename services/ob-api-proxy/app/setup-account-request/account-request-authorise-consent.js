const { setupAccountRequest } = require('./setup-account-request');
const { deleteRequest } = require('./delete-account-request');
const { generateRedirectUri, setConsent } = require('../authorise');
const { extractHeaders } = require('../session');

const uuidv4 = require('uuid/v4');
const error = require('debug')('error');
const debug = require('debug')('debug');
const _ = require('lodash');

const DefaultPermissions = [
  'ReadAccountsDetail',
  'ReadBalances',
  'ReadBeneficiariesDetail',
  'ReadDirectDebits',
  'ReadProducts',
  'ReadStandingOrdersDetail',
  'ReadTransactionsCredits',
  'ReadTransactionsDebits',
  'ReadTransactionsDetail',
];
// ExpirationDateTime: // not populated - the permissions will be open ended
// TransactionFromDateTime: // not populated - request from the earliest available transaction
// TransactionToDateTime: // not populated - request to the latest available transactions

const storePermissions = async (username, authorisationServerId, accountRequestId,
  validationRunId, permissions) => {
  const keys = {
    username, authorisationServerId, scope: 'accounts', validationRunId,
  };
  const accountRequest = { accountRequestId, permissions };
  await setConsent(keys, accountRequest);
};

const accountRequestAuthoriseConsent = async (req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  try {
    const headers = await extractHeaders(req.headers);
    const {
      authorisationServerId, username, sessionId, validationRunId, config,
    } = headers;
    const permissionsList = headers.permissions || DefaultPermissions;
    const headersWithPermissions = Object.assign({ permissions: permissionsList }, headers);
    // const { accountRequestId, permissions } = await setupAccountRequest(headersWithPermissions);
    const response = await setupAccountRequest(headersWithPermissions);
    const accountRequestId = _.get(response, 'Data.AccountRequestId');
    const permissions = _.get(response, 'Data.Permissions');
    const validation_result = _.get(response, 'validation_result'); // eslint-disable-line

    const interactionId2 = uuidv4();
    const uri = await generateRedirectUri(
      authorisationServerId, accountRequestId,
      'openid accounts', sessionId, interactionId2, config,
    );

    debug('services/ob-api-proxy/app/setup-account-request/account-request-authorise-consent.js:accountRequestAuthoriseConsent -> authorize uri=%O', uri);
    debug('services/ob-api-proxy/app/setup-account-request/account-request-authorise-consent.js:accountRequestAuthoriseConsent -> response=%O', response);

    await storePermissions(
      username, authorisationServerId, accountRequestId,
      validationRunId, permissions,
    );

    return res
      .status(200) // We can't intercept a 302 !
      .send({
        uri,
        validation_result,
      });
  } catch (err) {
    error(err);
    const status = err.status ? err.status : 500;
    return res.status(status).send({ message: err.message });
  }
};

const accountRequestRevokeConsent = async (req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  try {
    const headers = await extractHeaders(req.headers);
    const status = await deleteRequest(headers);
    return res.sendStatus(status);
  } catch (err) {
    error(`accountRequestRevokeConsent: ${err}`);
    const status = err.status ? err.status : 500;
    return res.status(status).send({ message: err.message });
  }
};

module.exports = {
  accountRequestAuthoriseConsent,
  accountRequestRevokeConsent,
  DefaultPermissions,
};
