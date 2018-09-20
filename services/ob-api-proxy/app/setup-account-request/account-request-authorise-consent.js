const { setupAccountRequest } = require('./setup-account-request');
const { deleteRequest } = require('./delete-account-request');
const { generateRedirectUri, setConsent } = require('../authorise');
const { extractHeaders } = require('../session');

const uuidv4 = require('uuid/v4');
const error = require('debug')('error');
const debug = require('debug')('debug');

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
    const { accountRequestId, permissions } = await setupAccountRequest(headersWithPermissions);
    const interactionId2 = uuidv4();
    const uri = await generateRedirectUri(
      authorisationServerId, accountRequestId,
      'openid accounts', sessionId, interactionId2, config,
    );

    debug(`authorize URL is: ${uri}`);
    await storePermissions(
      username, authorisationServerId, accountRequestId,
      validationRunId, permissions,
    );
    return res.status(200).send({ uri }); // We can't intercept a 302 !
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
