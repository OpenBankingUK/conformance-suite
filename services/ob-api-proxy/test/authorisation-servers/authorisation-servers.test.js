const assert = require('assert');

const { drop } = require('../../app/storage.js');
const {
  ASPSP_AUTH_SERVERS_COLLECTION,
  NO_SOFTWARE_STATEMENT_ID,
} = require('../../app/authorisation-servers');
const {
  allAuthorisationServers,
  storeAuthorisationServers,
  updateOpenIdConfigs,
  requestObjectSigningAlgs,
  updateRegisteredConfig,
  getRegisteredConfig,
} = require('../../app/authorisation-servers');

const nock = require('nock');

const authServerId = 'aaaj4NmBD8lQxmLh2O9FLY';
const baseApiDNSUri = 'http://aaa.example.com/some/path/open-banking/v1.1';
const orgId = 'aaa-example-org';
const flattenedObDirectoryAuthServerList = [
  {
    Id: authServerId,
    BaseApiDNSUri: baseApiDNSUri,
    CustomerFriendlyName: 'AAA Example Bank',
    OpenIDConfigEndPointUri: 'http://example.com/openidconfig',
    OBOrganisationId: orgId,
  },
];

const expectedAuthEndpoint = 'http://auth.example.com/authorize';
const expectedRequestAlgorithms = ['HS256', 'RS256'];
const expectedIdTokenAlgorithms = ['HS256', 'PS256'];
const openIdConfig = {
  authorization_endpoint: expectedAuthEndpoint,
  id_token_signing_alg_values_supported: expectedIdTokenAlgorithms,
  request_object_signing_alg_values_supported: expectedRequestAlgorithms,
};

const newClientCredentials = {
  clientId: 'a-client-id',
  clientSecret: 'a-client-secret',
};

const clientCredentials = [
  Object.assign(
    { softwareStatementId: NO_SOFTWARE_STATEMENT_ID },
    newClientCredentials,
  ),
];

const registeredConfig = {
  request_object_signing_alg: ['PS256'],
};

const registeredConfigs = [
  Object.assign(
    { softwareStatementId: NO_SOFTWARE_STATEMENT_ID },
    registeredConfig,
  ),
];

const withOpenIdConfig = {
  id: authServerId,
  obDirectoryConfig: {
    BaseApiDNSUri: baseApiDNSUri,
    CustomerFriendlyName: 'AAA Example Bank',
    OpenIDConfigEndPointUri: 'http://example.com/openidconfig',
    Id: authServerId,
    OBOrganisationId: 'aaa-example-org',
  },
  openIdConfig,
};

const withRegisteredConfig = {
  id: authServerId,
  obDirectoryConfig: {
    BaseApiDNSUri: baseApiDNSUri,
    CustomerFriendlyName: 'AAA Example Bank',
    OpenIDConfigEndPointUri: 'http://example.com/openidconfig',
    Id: authServerId,
    OBOrganisationId: 'aaa-example-org',
  },
  registeredConfigs,
};

const callAndGetLatestConfig = async (fn, authorisationServerId, data) => {
  if (fn) await fn(authorisationServerId, data);
  const list = await allAuthorisationServers();
  return list[0];
};

describe('authorisation servers', () => {
  beforeEach(async () => {
    await drop(ASPSP_AUTH_SERVERS_COLLECTION);
    await storeAuthorisationServers(flattenedObDirectoryAuthServerList);
  });

  afterEach(async () => {
    await drop(ASPSP_AUTH_SERVERS_COLLECTION);
  });

  describe('updateRegisteredConfig', () => {
    it('before called registered config not present', async () => {
      const authServerConfig = await callAndGetLatestConfig();
      assert.ok(!authServerConfig.registeredConfigs, 'registeredConfig not present');
    });

    it('stores registeredConfig in db when not OB provisioned', async () => {
      const authServerConfig = await callAndGetLatestConfig(
        updateRegisteredConfig,
        authServerId,
        registeredConfig,
      );
      assert.ok(authServerConfig.registeredConfigs, 'registeredConfig present');
      assert.deepEqual(authServerConfig, withRegisteredConfig);
    });
  });

  describe('getRegisteredConfig', () => {
    beforeEach(async () => {
      await updateRegisteredConfig(authServerId, registeredConfig);
    });

    describe('called with invalid authServerId', () => {
      it('throws error', async () => {
        try {
          await getRegisteredConfig('invalid-id');
          assert.ok(false);
        } catch (err) {
          assert.equal(err.status, 500);
        }
      });
    });

    it('retrieves registered config for an authorisationServerId', async () => {
      const found = await getRegisteredConfig(authServerId);
      assert.deepEqual(found, registeredConfigs[0]);
    });
  });

  describe('updateOpenIdConfigs', () => {
    nock(/example\.com/)
      .get('/openidconfig')
      .reply(200, openIdConfig);

    it('before called openIdConfig not present', async () => {
      const authServerConfig = await callAndGetLatestConfig();
      assert.ok(!authServerConfig.openIdConfig, 'openIdConfig not present');
    });

    it('retrieves openIdConfig and stores in db', async () => {
      const authServerConfig = await callAndGetLatestConfig(updateOpenIdConfigs);
      assert.ok(authServerConfig.openIdConfig, 'openIdConfig present');
      assert.deepEqual(authServerConfig, withOpenIdConfig);

      const requestAlgorithms = await requestObjectSigningAlgs(authServerId);
      assert.deepEqual(requestAlgorithms, expectedRequestAlgorithms);
    });
  });
});

exports.flattenedObDirectoryAuthServerList = flattenedObDirectoryAuthServerList;
exports.clientCredentials = clientCredentials;
exports.openIdConfig = openIdConfig;
