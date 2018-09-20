const proxyquire = require('proxyquire');
const assert = require('assert');
const nock = require('nock');
const sinon = require('sinon');
const qs = require('qs');

const clientId = 's6BhdRkqt3';
const clientSecret = '7Fjfp0ZBr1KtDRbnfVdmIw';
const basicAuth = 'Basic czZCaGRSa3F0Mzo3RmpmcDBaQnIxS3REUmJuZlZkbUl3';
const jwtToken = 'jwt.token';

const createJwtStub = sinon.stub().returns(jwtToken);

const obtainAccessToken = () =>
  proxyquire('../../app/authorise/obtain-access-token', {
    '../ob-util': {
      createJwt: createJwtStub,
    },
  });

const response = {
  access_token: 'accessToken',
  expires_in: 3600,
  token_type: 'bearer',
  scope: 'accounts payments',
};

const assertError = (error, name, status, message) => {
  assert.equal(error.name, name);
  assert.equal(error.message, message);
  assert.equal(error.status, status);
};

describe('POST /token', () => {
  const redirectionUrl = 'http://example.com/redirect';
  const authorisationCode = 'mockAuthorisationCode';
  const scenarios = [{
    authMethod: 'private_key_jwt',
    config: {
      client_id: clientId,
      client_secret: '',
      signing_key: 'x',
      token_endpoint: 'http://example.com/token',
      token_endpoint_auth_method: 'private_key_jwt',
    },
    clientCredentialsRequest: {
      scope: 'accounts payments',
      grant_type: 'client_credentials',
      client_assertion_type: 'urn:ietf:params:oauth:client-assertion-type:jwt-bearer',
      client_assertion: jwtToken,
    },
    authorisationCodeRequest: {
      grant_type: 'authorization_code',
      redirect_uri: redirectionUrl,
      code: authorisationCode,
      client_id: clientId,
      client_assertion_type: 'urn:ietf:params:oauth:client-assertion-type:jwt-bearer',
      client_assertion: jwtToken,
    },
  },
  {
    authMethod: 'client_secret_basic',
    config: {
      client_id: clientId,
      client_secret: clientSecret,
      signing_key: 'x',
      token_endpoint: 'http://example.com/token',
      token_endpoint_auth_method: 'client_secret_basic',
    },
    clientCredentialsRequest: {
      scope: 'accounts payments',
      grant_type: 'client_credentials',
    },
    authorisationCodeRequest: {
      grant_type: 'authorization_code',
      redirect_uri: redirectionUrl,
      code: authorisationCode,
      client_id: clientId,
    },
  }];

  scenarios.forEach(({
    authMethod, config, clientCredentialsRequest, authorisationCodeRequest,
  }) => {
    describe(`obtainClientCredentialsAccessToken ${authMethod}`, () => {
      const { obtainClientCredentialsAccessToken } = obtainAccessToken();

      it('returns token when 200 OK', async () => {
        const body = qs.stringify(clientCredentialsRequest);
        nock(/example\.com/)
          .post('/token', body)
          .reply(200, function check() {
            assert.equal(this.req.headers['content-type'], 'application/x-www-form-urlencoded');

            if (authMethod === 'client_secret_basic') {
              assert.equal(this.req.headers['authorization'], basicAuth);
            }

            return response;
          });

        const token = await obtainClientCredentialsAccessToken(config);

        assert.equal(token, 'accessToken');
      });

      it('throws error with response status when non 200 response', async () => {
        nock(/example\.com/)
          .post('/token')
          .reply(403, { error: 'message' });
        try {
          await obtainClientCredentialsAccessToken(config);
          assert.ok(false);
        } catch (error) {
          assertError(error, 'Error', 403, 'Forbidden {"error":"message"}');
        }
      });

      it('throws error with status set to 500 when error sending request', async () => {
        try {
          const badConfig = Object.assign({}, config, { token_endpoint: 'bad-url' });
          await obtainClientCredentialsAccessToken(badConfig);
          assert.ok(false);
        } catch (error) {
          assertError(error, 'Error', 500, 'getaddrinfo ENOTFOUND bad-url bad-url:80 ');
        }
      });
    });

    describe(`obtainAuthorisationCodeAccessToken ${authMethod}`, () => {
      const { obtainAuthorisationCodeAccessToken } = obtainAccessToken();

      it('returns token from obtainAuthorisationCodeAccessToken when 200 OK', async () => {
        const body = qs.stringify(authorisationCodeRequest);

        nock(/example\.com/)
          .post('/token', body)
          .reply(200, function check() {
            assert.equal(this.req.headers['content-type'], 'application/x-www-form-urlencoded');

            if (authMethod === 'client_secret_basic') {
              assert.equal(this.req.headers['authorization'], basicAuth);
            }

            return response;
          });
        const { accessToken } = await obtainAuthorisationCodeAccessToken(
          redirectionUrl, authorisationCode,
          config,
        );
        assert.equal(accessToken, 'accessToken');
      });

      it('returns token from obtainAuthorisationCodeAccessToken when 200 OK and no scope or expires_in on response', async () => {
        const body = qs.stringify(authorisationCodeRequest);

        nock(/example\.com/)
          .post('/token', body)
          .reply(200, function () { // eslint-disable-line
            assert.equal(this.req.headers['content-type'], 'application/x-www-form-urlencoded');
            if (authMethod === 'client_secret_basic') {
              assert.equal(this.req.headers['authorization'], basicAuth);
            }
            return {
              access_token: 'accessToken',
              token_type: 'bearer',
            };
          });
        const { accessToken, tokenExpiresAt } = await obtainAuthorisationCodeAccessToken(
          redirectionUrl, authorisationCode,
          config,
        );
        assert.equal(accessToken, 'accessToken');
        assert.equal(tokenExpiresAt, null);
      });

      it('throws error with response status when non 200 response', async () => {
        nock(/example\.com/)
          .post('/token')
          .reply(403, { error: 'message' });
        try {
          await obtainAuthorisationCodeAccessToken(
            redirectionUrl, authorisationCode,
            config,
          );
          assert.ok(false);
        } catch (error) {
          assertError(error, 'Error', 403, 'Forbidden {"error":"message"}');
        }
      });

      it('throws error with status set to 500 when error sending request', async () => {
        try {
          await obtainAuthorisationCodeAccessToken(
            redirectionUrl, authorisationCode,
            Object.assign({}, config, { token_endpoint: 'bad-url' }),
          );
          assert.ok(false);
        } catch (error) {
          assertError(error, 'Error', 500, 'getaddrinfo ENOTFOUND bad-url bad-url:80 ');
        }
      });
    });
  });
});
