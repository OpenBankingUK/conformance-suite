const request = require('supertest');
const fs = require('fs');
const path = require('path');

const {
  flattenedObDirectoryAuthServerList,
  clientCredentials,
  openIdConfig,
} = require('./authorisation-servers/authorisation-servers.test');
const { extractAuthorisationServers } = require('../app/ob-directory');

const accessToken = 'AN_ACCESS_TOKEN';

const { drop } = require('../app/storage.js');

const { app } = require('../app/index.js');
const { session } = require('../app/session');
const {
  ASPSP_AUTH_SERVERS_COLLECTION,
  storeAuthorisationServers,
  updateClientCredentials,
  updateOpenIdConfig,
} = require('../app/authorisation-servers');

const assert = require('assert');
const nock = require('nock');

nock(/secure-url\.com/)
  .get('/private_key.pem')
  .reply(200, fs.readFileSync(path.join(__dirname, 'test_private_key.pem')));

nock(/auth\.com/)
  .post('/as/token.oauth2')
  .reply(200, {
    access_token: accessToken,
    token_type: 'Bearer',
    expires_in: 1000,
  });

const aspspPayload = {
  Resources: [
    {
      'urn:openbanking:competentauthorityclaims:1.0': {
        AuthorityId: 'FCA',
        MemberState: 'GB',
        RegistrationId: '123',
      },
      'AuthorisationServers': [
        {
          Id: 'aaaj4NmBD8lQxmLh2O9FLY',
          BaseApiDNSUri: 'http://aaa.example.com',
          CustomerFriendlyLogoUri: 'string',
          CustomerFriendlyName: 'AAA Example Bank',
          OpenIDConfigEndPointUri: 'http://aaa.example.com/openid/config',
        },
      ],
      'urn:openbanking:organisation:1.0': {
        OrganisationCommonName: 'AAA Group PLC',
        OBOrganisationId: 'aaax5nTR33811QyQfi',
      },
      'id': 'aaax5nTR33811QyQfi',
    },
    {
      'urn:openbanking:competentauthorityclaims:1.0': {
        AuthorityId: 'FCA',
        MemberState: 'GB',
        RegistrationId: '456',
      },
      'AuthorisationServers': [
        {
          Id: 'bbbX7tUB4fPIYB0k1m',
          BaseApiDNSUri: 'http://bbb.example.com',
          CustomerFriendlyLogoUri: 'string',
          CustomerFriendlyName: 'BBB Example Bank',
          OpenIDConfigEndPointUri: 'http://bbb.example.com/openid/config',
        },
        {
          Id: 'cccbN8iAsMh74sOXhk',
          BaseApiDNSUri: 'http://ccc.example.com',
          CustomerFriendlyLogoUri: 'string',
          CustomerFriendlyName: 'CCC Example Bank',
          OpenIDConfigEndPointUri: 'http://ccc.example.com/openid/config',
        },
      ],
      'urn:openbanking:organisation:1.0': {
        OrganisationCommonName: 'BBBCCC Group PLC',
        OBOrganisationId: 'bbbcccUB4fPIYB0k1m',
      },
      'id': 'bbbcccUB4fPIYB0k1m',
    },
    {
      id: 'fPIYB0k1moGhX7tUB4',
    },
  ],
};

describe('extractAuthorisationServers', () => {
  it('returns flattened ASPSP auth server list', () => {
    const list = extractAuthorisationServers(aspspPayload);
    const expected = [
      {
        BaseApiDNSUri: 'http://aaa.example.com',
        CustomerFriendlyLogoUri: 'string',
        CustomerFriendlyName: 'AAA Example Bank',
        Id: 'aaaj4NmBD8lQxmLh2O9FLY',
        OpenIDConfigEndPointUri: 'http://aaa.example.com/openid/config',
        OBOrganisationId: 'aaax5nTR33811QyQfi',
        OrganisationCommonName: 'AAA Group PLC',
        AuthorityId: 'FCA',
        MemberState: 'GB',
        RegistrationId: '123',
      },
      {
        BaseApiDNSUri: 'http://bbb.example.com',
        CustomerFriendlyLogoUri: 'string',
        CustomerFriendlyName: 'BBB Example Bank',
        Id: 'bbbX7tUB4fPIYB0k1m',
        OpenIDConfigEndPointUri: 'http://bbb.example.com/openid/config',
        OBOrganisationId: 'bbbcccUB4fPIYB0k1m',
        OrganisationCommonName: 'BBBCCC Group PLC',
        AuthorityId: 'FCA',
        MemberState: 'GB',
        RegistrationId: '456',
      },
      {
        BaseApiDNSUri: 'http://ccc.example.com',
        CustomerFriendlyLogoUri: 'string',
        CustomerFriendlyName: 'CCC Example Bank',
        Id: 'cccbN8iAsMh74sOXhk',
        OpenIDConfigEndPointUri: 'http://ccc.example.com/openid/config',
        OBOrganisationId: 'bbbcccUB4fPIYB0k1m',
        OrganisationCommonName: 'BBBCCC Group PLC',
        AuthorityId: 'FCA',
        MemberState: 'GB',
        RegistrationId: '456',
      },
    ];
    assert.deepEqual(list, expected);
  });
});

const login = application => request(application)
  .post('/login')
  .set('Accept', 'x-www-form-urlencoded')
  .send({ u: 'alice', p: 'wonderland' });

describe('Directory', () => {
  beforeEach(async () => {
    await drop(ASPSP_AUTH_SERVERS_COLLECTION);
    const config = flattenedObDirectoryAuthServerList[0];
    await storeAuthorisationServers([config]);
    await updateClientCredentials(config.Id, clientCredentials);
    await updateOpenIdConfig(config.Id, openIdConfig);
  });

  afterEach(async () => {
    await session.deleteAll();
    await drop(ASPSP_AUTH_SERVERS_COLLECTION);
  });

  const expectedResult = [
    {
      id: 'aaaj4NmBD8lQxmLh2O9FLY',
      name: 'AAA Example Bank',
    },
  ];

  it('returns 200 response for /account-payment-service-provider-authorisation-servers', (done) => {
    login(app).end((err, res) => {
      const sessionId = res.body.sid;

      request(app)
        .get('/account-payment-service-provider-authorisation-servers')
        .set('Accept', 'application/json')
        .set('x-validation-run-id', 'validationRunId')
        .set('authorization', sessionId)
        .end((e, r) => {
          assert.equal(r.status, 200);
          const header = r.headers['access-control-allow-origin'];
          assert.equal(header, '*');
          assert.deepEqual(r.body, expectedResult);
          done();
        });
    });
  });
});
