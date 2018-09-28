const request = require('supertest');

const authorization = 'abc';
const fapiFinancialId = 'xyz';
const authServerId = 'testAuthServerId';
const validationRunId = 'validationRunId';
const { mungeToken } = require('../app/authorise/obtain-access-token');
const { base64EncodeJSON } = require('../app/ob-util');

process.env.DEBUG = 'error';
process.env.AUTHORIZATION = authorization;
process.env.OB_DIRECTORY_HOST = 'http://example.com';

const { app } = require('../app/index.js');
const { session } = require('../app/session');
const assert = require('assert');

const nock = require('nock');

const username = 'alice';
const password = 'wonderland';

const ACCOUNT_SWAGGER = 'https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/ee715e094a59b37aeec46aef278f528f5d89eb03/dist/v1.1/account-info-swagger.json';

const config = {
  client_id: 'clientId',
  client_secret: 'clientSecret',
  fapi_financial_id: fapiFinancialId,
  resource_endpoint: 'http://example.com',
  transport_cert: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
  transport_key: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
};

const requestHeaders = {
  reqheaders: {
    'authorization': authorization,
    'x-authorization-server-id': authServerId,
    'x-swagger-uris': ACCOUNT_SWAGGER,
    'x-config': base64EncodeJSON(config),
  },
};

nock(/example\.com/, requestHeaders)
  .get('/open-banking/v1.1/accounts')
  .times(5)
  .reply(200, {
    Data: {
      Account: [
        {
          AccountId: '22290',
          Currency: 'GBP',
          Account: {
            SchemeName: 'SortCodeAccountNumber',
            Identification: '30854645679085',
            Name: 'Ms Smith',
            SecondaryIdentification: '341267',
          },
        },
      ],
    },
    Links: {
      Self: '/accounts',
    },
    Meta: {
      TotalPages: 1,
    },
  });

// bad payload to trigger validation error
nock(/example\.com/, requestHeaders)
  .get('/open-banking/v1.1/accounts/22290/balances')
  .reply(200, { Data: { Balance: {} }, Links: { Self: '' }, Meta: {} });

nock(/example\.com/)
  .get('/open-banking/v1.1/accounts/non-existing')
  .reply(404);

const login = application => request(application)
  .post('/login')
  .set('Accept', 'x-www-form-urlencoded')
  .send({ u: username, p: password });

describe('Session Creation (Login)', () => {
  it('returns "Access-Control-Allow-Origin: *" header', (done) => {
    login(app)
      .end((err, res) => {
        const header = res.headers['access-control-allow-origin'];
        assert.equal(header, '*');
        done();
      });
  });

  it('returns a guid in the body as a json payload for /login', (done) => {
    login(app)
      .end((err, res) => {
        const mySid = res.body.sid;
        const isGuid = (mySid.length === 36);
        assert.equal(true, isGuid);
        done();
      });
  });

  it('returns a 200 for second user kate login /login', (done) => {
    request(app)
      .post('/login')
      .set('Accept', 'x-www-form-urlencoded')
      .send({ u: 'kate', p: 'lookingglass' })
      .end((err, res) => {
        assert.equal(res.status, 200);
        done();
      });
  });


  it('returns an unauthorised status for an invalid set of credentials at /login', (done) => {
    request(app)
      .post('/login')
      .set('Accept', 'x-www-form-urlencoded')
      .send({ u: 'foo', p: 'baarx' })
      .end((err, res) => {
        assert.equal(res.status, 401);
        done();
      });
  });

  it('returns 500 error status for username "trigger-error"', (done) => {
    request(app)
      .post('/login')
      .set('Accept', 'x-www-form-urlencoded')
      .send({ u: 'trigger-error', p: 'baarx' })
      .end((err, res) => {
        assert.equal(res.status, 500);
        done();
      });
  });
});

describe('Cross Origin Requests Handled Correctly', () => {
  it('returns "Access-Control-Allow-Origin: *" header', (done) => {
    login(app).end(() => {
      request(app)
        .post('/logout')
        .end((e, r) => {
          const header = r.headers['access-control-allow-origin'];
          assert.equal(header, '*');
          done();
        });
    });
  });
});

describe('Session Deletion (Logout)', () => {
  it('destroys a valid session at /logout', (done) => {
    login(app).end((err, res) => {
      const sessionId = res.body.sid;

      request(app)
        .post('/logout')
        .set('Accept', 'application/json')
        .set('authorization', sessionId)
        .end((e, r) => {
          assert.equal(r.status, 200);
          assert.equal(r.body.sid, sessionId);
          done();
        });
    });
  });

  it('does not destroy an invalid session at /logout', (done) => {
    request(app)
      .get('/logout')
      .set('Accept', 'application/json')
      .set('authorization', 'jkaghrtegdkhsugf')
      .end((err, res) => {
        assert.equal(res.status, 204);
        done();
      });
  });

  after(async () => {
    await session.deleteAll();
  });
});
const scope = 'accounts';
const accountRequestId = 'xxxxxx-xxxx-43c6-9c75-eaf01821375e';
const authorisationCode = 'spoofAuthCode';
const token = 'testAccessToken';
const tokenPayload = {
  access_token: token,
  expires_in: 3600,
};

const consentPayload = {
  username,
  authorisationServerId: authServerId,
  scope,
  accountRequestId,
  expirationDateTime: null,
  authorisationCode,
  token: mungeToken(tokenPayload),
};

const { AUTH_SERVER_USER_CONSENTS_COLLECTION } = require('../app/authorise/consents');
const { ASPSP_AUTH_SERVERS_COLLECTION } = require('../app/authorisation-servers');
const { setConsent } = require('../app/authorise');
const { setAuthServerConfig } = require('../app/authorisation-servers');
const { drop } = require('../app/storage.js');

const resourceApiHost = 'http://example.com';

const loginAsync = async (application) => {
  const res = await login(application);
  const sessionId = res.body.sid;
  return { sessionId, res };
};

const requestResource = async (sessionId, url, application) => {
  const req = request(application)
    .get(url)
    .set('Accept', 'application/json')
    .set('x-authorization-server-id', authServerId)
    .set('x-validation-run-id', validationRunId)
    .set('x-config', base64EncodeJSON(config))
    .set('x-swagger-uris', ACCOUNT_SWAGGER);
  if (sessionId) {
    req.set('authorization', sessionId);
  }
  return req;
};

describe('Proxy', () => {
  beforeEach(async () => {
    await setAuthServerConfig(authServerId, {
      obDirectoryConfig: {
        OBOrganisationId: fapiFinancialId,
        BaseApiDNSUri: resourceApiHost,
      },
    });
    await setConsent({
      username,
      authorisationServerId: authServerId,
      scope,
      validationRunId,
    }, consentPayload);
  });

  afterEach(async () => {
    await session.deleteAll();
    await drop(AUTH_SERVER_USER_CONSENTS_COLLECTION);
    await drop(ASPSP_AUTH_SERVERS_COLLECTION);
    delete process.env.DEBUG;
    delete process.env.OB_DIRECTORY_HOST;
    delete process.env.AUTHORIZATION;
    delete process.env.VALIDATE_RESPONSE;
  });

  it('returns proxy 200 response for /open-banking/v1.1/accounts with valid session', async () => {
    const { sessionId } = await loginAsync(app);
    const r = await requestResource(sessionId, '/open-banking/v1.1/accounts', app);
    assert.equal(r.status, 200, `expected 200, got: ${r.status} ${JSON.stringify(r.body)}`);
    assert.equal(r.body.Data.Account[0].AccountId, '22290');
  });

  it('sets failedValidation false on response when VALIDATE_RESPONSE is true and validation passes', async () => {
    process.env.VALIDATE_RESPONSE = 'true';
    const { sessionId } = await loginAsync(app);
    const r = await requestResource(sessionId, '/open-banking/v1.1/accounts', app);
    assert.equal(r.body.validation_result.report.failedValidation, false, `Expect failedValidation: false in response JSON, got: ${r.body.failedValidation}`);
  }).timeout(1000 * 5);

  it('sets failedValidation true on response when VALIDATE_RESPONSE is true and validation fails', async () => {
    process.env.VALIDATE_RESPONSE = 'true';
    const { sessionId } = await loginAsync(app);
    const r = await requestResource(sessionId, '/open-banking/v1.1/accounts/22290/balances', app);
    assert.equal(r.body.validation_result.report.failedValidation, true, `Expect failedValidation: true in response JSON, got: ${r.body.failedValidation}`);
  });

  it('does not set failedValidation on response when VALIDATE_RESPONSE is false', async () => {
    process.env.VALIDATE_RESPONSE = 'false';
    const { sessionId } = await loginAsync(app);
    const r = await requestResource(sessionId, '/open-banking/v1.1/accounts', app);
    assert.equal(r.body.failedValidation, undefined);
  });

  it('returns 400 response for missing x-authorization-server-id', (done) => {
    loginAsync(app).then(({ sessionId }) => {
      request(app)
        .get('/open-banking/v1.1/accounts')
        .set('Accept', 'application/json')
        .set('authorization', sessionId)
        .end((e, r) => {
          assert.equal(r.status, 400);
          done();
        });
    });
  });

  it('returns proxy 404 reponse for /open-banking/non-existing', async () => {
    const { sessionId } = await loginAsync(app);
    const r = await requestResource(sessionId, '/open-banking/v1.1/accounts/non-existing', app);
    assert.equal(r.status, 404);
  });

  it('returns 404 for path != /open-banking', async () => {
    const { sessionId } = await loginAsync(app);
    const r = await requestResource(sessionId, '/open-banking-invalid', app);
    assert.equal(r.status, 404);
  });

  it('returns proxy 401 unauthorised response for /open-banking/* with missing authorization header', async () => {
    await loginAsync(app);
    const r = await requestResource(null, '/open-banking/v1.1/accounts', app);
    assert.equal(r.status, 401);
    assert.equal(r.headers['access-control-allow-origin'], '*');
  });

  it('returns proxy 401 unauthorised response for /open-banking/* with invalid authorization header', async () => {
    await loginAsync(app);
    const r = await requestResource('invalid-token', '/open-banking/v1.1/accounts', app);
    assert.equal(r.status, 401);
    assert.equal(r.headers['access-control-allow-origin'], '*');
  });
});
