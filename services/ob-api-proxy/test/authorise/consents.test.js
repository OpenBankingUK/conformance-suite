const assert = require('assert');
const proxyquire = require('proxyquire');
const sinon = require('sinon');
const { mungeToken } = require('../../app/authorise/obtain-access-token');

const {
  setConsent,
  consent,
  consentAccessToken,
  consentAccessTokenAndPermissions,
  consentAccountRequestId,
  deleteConsent,
  getConsent,
} = require('../../app/authorise');
const { AUTH_SERVER_USER_CONSENTS_COLLECTION } = require('../../app/authorise/consents');

const { drop } = require('../../app/storage.js');

const username = 'testUsername';
const sessionId = 'testSessionId';
const validationRunId = 'testValidationRunId';
const authorisationServerId = 'a123';
const scope = 'accounts';
const keys = {
  username, authorisationServerId, scope, validationRunId,
};

const accountRequestId = 'xxxxxx-xxxx-43c6-9c75-eaf01821375e';
const authorisationCode = 'spoofAuthCode';
const token = 'testAccessToken';
const tokenPayload = {
  access_token: token,
  expires_in: 3600,
  token_type: 'bearer',
};
const permissions = ['ReadAccountsDetail'];

const accountRequestPayload = {
  username,
  authorisationServerId,
  scope,
  accountRequestId,
  permissions,
};

const consentPayload = {
  username,
  authorisationServerId,
  scope,
  accountRequestId,
  expirationDateTime: null,
  authorisationCode,
  token: mungeToken(tokenPayload),
};

const consentStatus = 'Authorised';

const fapiFinancialId = 'testFapiFinancialId';
const resourcePath = 'http://example.com';
const config = {
  api_version: '1.1',
  client_id: 'testClientId',
  client_secret: undefined,
  fapi_financial_id: fapiFinancialId,
  resource_endpoint: resourcePath,
};

describe('setConsents', () => {
  beforeEach(async () => {
    await drop(AUTH_SERVER_USER_CONSENTS_COLLECTION);
  });

  afterEach(async () => {
    await drop(AUTH_SERVER_USER_CONSENTS_COLLECTION);
  });

  it('stores account request payload and allows to be retrieved', async () => {
    await setConsent(keys, accountRequestPayload);
    const stored = await consent(keys);
    assert.equal(stored.id, `${username}:::${authorisationServerId}:::${scope}:::${validationRunId}`);
  });

  it('stores consent payload, keeping permissions from stored account request with same accountRequestId', async () => {
    await setConsent(keys, accountRequestPayload);
    await setConsent(keys, consentPayload);
    const stored = await consent(keys);
    assert.deepEqual(stored.permissions, accountRequestPayload.permissions);
  });

  it('stores consent payload, without permissions from stored account request with different accountRequestId', async () => {
    const accountRequestWithDifferentId = Object.assign({}, accountRequestPayload, { accountRequestId: 'differentId' });
    await setConsent(keys, accountRequestWithDifferentId);
    await setConsent(keys, consentPayload);
    const stored = await consent(keys);
    assert.equal(stored.permissions, null);
  });

  it('stores payload and allows consent access_token to be retrieved', async () => {
    await setConsent(keys, consentPayload);
    const storedAccessToken = await consentAccessToken(keys);
    assert.equal(storedAccessToken, token);
  });

  it('stores payload and allows consent access_token and permissions to be retrieved', async () => {
    await setConsent(keys, accountRequestPayload);
    await setConsent(keys, consentPayload);
    const data = await consentAccessTokenAndPermissions(keys);
    assert.equal(data.accessToken, token);
    assert.deepEqual(data.permissions, permissions);
  });

  it('stores payload and allows consent accountRequestId to be retrieved', async () => {
    await setConsent(keys, consentPayload);
    const storedAccountRequestId = await consentAccountRequestId(keys);
    assert.equal(storedAccountRequestId, accountRequestId);
  });
});

describe('deleteConsent', () => {
  beforeEach(async () => {
    await drop(AUTH_SERVER_USER_CONSENTS_COLLECTION);
  });

  afterEach(async () => {
    await drop(AUTH_SERVER_USER_CONSENTS_COLLECTION);
  });

  it('stores payload and allows consent to be retrieved by keys id', async () => {
    await setConsent(keys, consentPayload);
    await deleteConsent(keys);
    const result = await getConsent(keys);
    assert.equal(result, null);
  });
});

describe('getConsentStatus', () => {
  const interactionId = 'testInteractionId';
  const accessToken = 'grant-credential-access-token';
  const getAccountRequestStub = sinon.stub();
  let getConsentStatus;

  describe('successful', () => {
    beforeEach(() => {
      getAccountRequestStub.returns({ Data: { Status: consentStatus } });
      ({ getConsentStatus } = proxyquire(
        '../../app/authorise/consents.js',
        {
          './obtain-access-token': {
            obtainClientCredentialsAccessToken: () => accessToken,
          },
          '../setup-account-request/account-requests': {
            getAccountRequest: getAccountRequestStub,
          },
          'uuid/v4': () => interactionId,
        },
      ));
    });

    it('makes remote call to get account request', async () => {
      await getConsentStatus(accountRequestId, authorisationServerId, sessionId, config);
      const headers = {
        accessToken, fapiFinancialId, interactionId, sessionId, authorisationServerId,
      };
      assert(getAccountRequestStub.calledWithExactly(accountRequestId, resourcePath, headers));
    });

    it('gets the status for an existing consent', async () => {
      const actual =
        await getConsentStatus(accountRequestId, authorisationServerId, sessionId, config);
      assert.equal(actual, consentStatus);
    });
  });

  describe('errors', () => {
    it('throws error for missing payload', async () => {
      try {
        await getConsentStatus(accountRequestId, authorisationServerId, sessionId, config);
      } catch (err) {
        assert(err);
      }
    });

    it('throws error for missing Data payload', async () => {
      getAccountRequestStub.returns({});
      try {
        await getConsentStatus(accountRequestId, authorisationServerId, sessionId);
      } catch (err) {
        assert(err);
      }
    });
  });
});
