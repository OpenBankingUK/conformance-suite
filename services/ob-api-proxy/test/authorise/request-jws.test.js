const assert = require('assert');
const sinon = require('sinon');
const proxyquire = require('proxyquire');
const { base64Decode } = require('../../app/ob-util');

const jwsStub = {};
const { createClaims, createJsonWebSignature } = proxyquire('../../app/authorise/request-jws', {
  'node-jose': { JWS: jwsStub },
});
const { statePayload } = require('../../app/authorise/authorise-uri');

const payload = { example: 'claims' };

/**
  To generate signing cert/key on command line:

  org="my-org-id-from-open-banking-directory"
  client="id-for-my-client-from-open-banking-directory
  openssl req -new -newkey rsa:2048 -nodes -sha256 \
    -out signing.csr -keyout signing.key \
    -subj "/C=GB/O=Open Banking Limited/OU=${org}/CN=${client}"
*/

// npm run base64-cert-or-key signing.key
const signingKey = base64Decode('LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdkYxbTQwc2FBeS8vYkJtVFA0WG5LdE1aSXFScy9Xc241aXN3UjIrVkNFSGdGNFlIClNic1hhbTB1eG00V25veDJQbDdlaDBiZUtEMmVXcDRwaCt3VHYyeWJCdkN5UXMxV1MwbUFnZFYxTlFJRVQ4UHUKbmxyOW5VMERGeWhuN0RHSHp1MTRrMFFvQk8vVFZkSmdXOVJLRCtRUkw3SW1BMzZWNGFrWkgwdUlRa1RTTlM4dgpGcWlhbFEyRHArMlFVN0FyVnZYdHd1eHVQMEdBcFA4NTJyRG5xeUlacVltcXJudThWM2k3Y04xbkM3WHh5RW5oCmN4UWpSb1lwWHcxMS8wLzRUOVFGNHhYVzMvaVp3VzJ1VktTdWtMNElsc1M1UThSUEhXeW4xdGkwNnRwaVpVcisKbGxwVTg5dXMzM0VwaXJQYnNDSWtvR2ZkWW5WRnhMOGRsTGp2MXdJREFRQUJBb0lCQUYyOFNTZ1F4bmdSbVpUTQp3VmJhSnFoTDluVWp4OHp3VnlHV0dtZGlJcExDWFdhM1hzY1ZJRmpvemw4V2g1RU1xd2JzcE9aQ29PajdpT0xsClZCdDhvbk1lODZLbmdyMzFldHpxVGRYT1NJNUJXNjNwL2NPMTJnRStRcXh5Z2d5cXRUK0hNdnB0NzFCTm5DaFkKRVhXQkZmNEVhMzBGdFI4R0RrWUdwU2JLcXByMjBvbzZGLzFUZ3paM2RKazlMbTg3K2pmV1RmTlRYaXRTemJJRApMUUMwNkFBVkJ1V241RlJjNmVTTENKUTRzbDI0Vy84THYyVjFsaXhkRzVNdmlpWUFJc0thVTF6N3g1MCt3ZnBiCkN5YWNmc2V0YkJNVmltczNWTmlhTjA4YWVtY3Uvak9jSklNcTF4aFZpNFJmVFoxaVgrMk5MVDZ3bGJhZStzd2YKTms0TENnRUNnWUVBOWFkZEVyTDlVa1FCTFI4RGNzSll1elo1QVdxSVFZYkhTVlpVdHBzZ3BaYzZuRW9hbWZzdwpCRUhkN1I0eHcrVnBNbGk1Tm9QYUlFaTN3MDR2QThYakx4Ny9kNXplbWZnS3VkQzRYdkxFc251UWMxbkVMSExjCjdDZ21iTm9rZW8zRkkrYzlYeldibVZBdnVSMWR2Y1AzNHYwYWpHb2RVRHFqKzNvSGNZYjNDdGNDZ1lFQXhFeFoKai9yTlNtdWpNeGFhSVYwYlhIa1BJZEhpcU9WMnZ4RnVuckVIMTNqVnA5RjNrd2ljaTNLam5MNGZNK3JGTmVBdApqQVZnM3NnT3dCMVRUUTlYUVprb3VFSndZT2Q0azNQY0ZWSkJ4N2x1RTlWSnVteWVjb0k4emIvSTRsMHlsMEpXCm40aFA4UFBoelVPZ3FOVkxvQWZPQ3UxNHRqOHFMUmNuTmtTZG93RUNnWUVBb2JUVFVzemFicjN2WEVsdkZxc1MKaCtKNjAxRFNjdmdLMVo3cjB1elpGOGd1UDlXVUgwcTN1QVczMWpBcktENHErb1puSFppOERNWnhtVEl0UnJtTQpMR2VtV1pHOUF2UEI4OEdPckluNHExa2xwSmt4eHVTeHd3OUhCQjZ4SnErT1YyMFAvRTJvcU1xZEw2bENIUG9VCmdxcUVRR3hWOFlzNGlRRXlSeXhHRVM4Q2dZQkJjNk81V2tyeE1ZcXRFakE2UjYxRDNDbXJnU3d1WExTSGFPeVYKaFRtMEl0bzZwcUZVS1Y3cE1FUlZreDhjVkgrRlEwWnNsYTZER2ZteEhSWVZiN1FNYjJFZ2J5YkJhT3pQWGFaWQpoYURoVTNiY3JoVnpUNXhWV2cra0d2cUVYOGJxb0hmNW9aM21IYXVBb2JnRUUzcXYxV3BpUW1RcGdFNHowckNFCmE4U1VBUUtCZ0VEQjlMY1NnQ0NRNnFtTithZ1NDRCtsVW84REJ4aE1LaFRNWWtoRFNPNkw5TFEvUXMxTFRibm4KVFZ2YXlEQS95V1hWOU1CREltVm9peDRYTElLQVlpbTVCVWw0ekcvRkxDdFMzU0w2dzU4dlBERlB6eGtObmNmQgo2dHpySm1PbHRBY3Vuc3J5MEV6UjRDSE9QU2F1Tk85Z2pwMjYydXVFcVZIMEFhbExLNXptCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==');
const signingKid = 'TAuIpZ8mDCMZXXXXXXXXXXXXXXX';

const notSigned = 'eyJhbGciOiJub25lIn0.eyJleGFtcGxlIjoiY2xhaW1zIn0.';
const hs256signed = 'eyJhbGciOiJIUzI1NiJ9.eyJleGFtcGxlIjoiY2xhaW1zIn0.713fdImuJuKedy8Lm5E2alESqRpT4cma1mXZVdX07k4';
const rs256signed = 'eyJhbGciOiJSUzI1NiJ9.eyJleGFtcGxlIjoiY2xhaW1zIn0.dWC5Q4R0c7JZgQ6VuiYpYsPNEeNKPW9Zj55VEDDrj0RqWZc72v4sK-en2qZGW_7ETGFhfUBvViP4Y746E-6eDsZFtzkFGclkXCDvPYr8ey0XsS0aASAaAjTMmY98FdTL-6FVEXwxZ8t6rvle77TGVBdq0wk0DTMishnpD_2roIQALWzDfPXtV31GY0L9L9uOcHopGYgSTAw81CLgWvr98an-li1Q8D8okp4yt_bmPw8mlKN8HI-bUiCStKIiYtze_h029VYphlD8ASqrmTjCQpRiJldTi3OroEoKPQk7MaYZH2kJuJMXnIKXdxfEcKcMh1YsUFhedd2ktlot73O8mw';
const ps256signed = 'eyJhbGciOiJQUzI1NiIsImtpZCI6IiJ9.eyJleGFtcGxlIjoiY2xhaW1zIn0.ky0VU7QFI3YtN-PrnZ6DJevEohLcD0A5THtLui2v2I8isdG_J8JdD6gvBGk8RR1w9RARuJbccUg5s7hR61GLmB49yDzyDqJ6Wc9nPvlABp_TMiq0va7aSuPKHfFNAXsHdM3fYSwkMJKEIFRQvIeazhlKi9i81qHHTjHHzvQMvcxVjThghYJvXOKYzPQug1hyHQsS0RBpdcI8NYOgK-HCGFGpJn7SeAEqcCFOuxOTImzJ57KEdUtlTT2WusUQof2dOhnak9lWUmvIjSODgpiVanVccUPmmPBCSoxR7sn9yBocaikMNTnUdGS7bhgxBy4Q_oDPCx2r1EGfJRW1WrUK7w';
const signatureConfigPS256 = {
  fields: { alg: 'PS256', kid: signingKid },
  format: 'compact',
};
const config = { client_secret: 'testClientSecret', signing_key: signingKey, signing_kid: signingKid };

describe('createJsonWebSignature when signing key blank and algos require signing key', () => {
  it('throws error', () => {
    try {
      createJsonWebSignature(payload, ['RS256'], Object.assign({}, config, { signing_key: '' }));
      assert.ok(false);
    } catch (error) {
      assert.equal(error.name, 'Error');
      assert.equal(error.message, 'cannot create JSON web signature for RS256 as signing key was undefined or blank');
    }
  });
});

describe('createJsonWebSignature when signing algo not recognized', () => {
  it('throws error', () => {
    try {
      createJsonWebSignature(payload, ['Random'], config);
      assert.ok(false);
    } catch (error) {
      assert.equal(error.name, 'Error');
      assert.equal(error.message, 'cannot create JSON web signature as Random not supported');
    }
  });
});

describe('createJsonWebSignature when signing key present', () => {
  before(() => {
    jwsStub.createSign = sinon.stub().returns(jwsStub);
    jwsStub.update = sinon.stub().returns(jwsStub);
    jwsStub.final = sinon.stub().returns(ps256signed);
  });

  describe('with signing algs ["none"]', () => {
    it('creates JWS', async () => {
      const jws = await createJsonWebSignature(payload, ['none'], config);
      assert.equal(jws, notSigned);
    });
  });

  describe('with signing algs ["HS256"]', () => {
    it('creates JWS', async () => {
      const jws = await createJsonWebSignature(payload, ['HS256'], config);
      assert.equal(jws, hs256signed);
    });
  });

  describe('with signing algs ["RS256"]', () => {
    it('creates JWS', async () => {
      const jws = await createJsonWebSignature(payload, ['RS256'], config);
      assert.equal(jws, rs256signed);
    });
  });

  describe('with signing algs ["PS256"]', () => {
    it('creates JWS', async () => {
      const jws = await createJsonWebSignature(payload, ['PS256'], config);
      assert.equal(jws, ps256signed);
      assert(jwsStub.createSign.calledWith(signatureConfigPS256));
      assert(jwsStub.update.calledWithExactly(JSON.stringify(payload), 'utf-8'));
    });
  });

  describe('with non-supported signing algs', () => {
    it('throws error', async () => {
      try {
        await createJsonWebSignature(payload, ['RS512'], config);
        assert.ok(false);
      } catch (error) {
        assert.equal(error.name, 'Error');
      }
    });
  });

  describe('with signing algs ["PS256", "HS256"]', () => {
    it('creates JWS with PS256', async () => {
      const jws = await createJsonWebSignature(payload, ['PS256', 'HS256'], config);
      assert.equal(jws, ps256signed);
    });
  });

  describe('with signing algs ["PS256", "none"]', () => {
    it('creates JWS with PS256', async () => {
      const jws = await createJsonWebSignature(payload, ['PS256', 'none'], config);
      assert.equal(jws, ps256signed);
    });
  });
});

describe('createJsonWebSignature when no signing key present', () => {
  const noSigningKeyConfig = Object.assign({}, config, { signing_key: '' });
  describe('with signing algs ["ES256", "none"]', () => {
    it('creates JWS defaulting to none', async () => {
      const jws = await createJsonWebSignature(payload, ['ES256', 'none'], noSigningKeyConfig);
      assert.equal(jws, notSigned);
    });
  });

  describe('with signing algs ["ES256", "HS256"]', () => {
    it('creates JWS defaulting to HS256', async () => {
      const jws = await createJsonWebSignature(payload, ['ES256', 'HS256'], noSigningKeyConfig);
      assert.equal(jws, hs256signed);
    });
  });
});

describe('createClaims', () => {
  const accountRequestId = 'testAccountRequestId';
  const clientId = 'testClientId';
  const scope = 'openid accounts';
  const authServerIssuer = 'http://aspsp.example.com';
  const registeredRedirectUrl = 'http://tpp.example.com/handle-authorise';
  const authorisationServerId = 'testAuthorisationServerId';
  const sessionId = 'testSessionId';
  const state = statePayload(authorisationServerId, sessionId);

  const expectedClaims = audience => ({
    aud: audience,
    iss: clientId,
    response_type: 'code',
    client_id: clientId,
    redirect_uri: registeredRedirectUrl,
    scope,
    state,
    nonce: 'dummy-nonce',
    max_age: 86400,
    claims:
    {
      userinfo:
      {
        openbanking_intent_id: { value: accountRequestId, essential: true },
      },
      id_token:
      {
        openbanking_intent_id: { value: accountRequestId, essential: true },
        acr: {
          essential: true,
        },
      },
    },
  });

  it('creates claims JSON successfully when useOpenidConnect is false', () => {
    const useOpenidConnect = false;
    const claims = createClaims(
      scope, accountRequestId, clientId, authServerIssuer,
      registeredRedirectUrl, state, useOpenidConnect,
    );
    assert.deepEqual(claims, expectedClaims(authServerIssuer));
  });

  it('creates claims JSON successfully when useOpenidConnect is true', () => {
    const useOpenidConnect = true;
    const claims = createClaims(
      scope, accountRequestId, clientId, authServerIssuer,
      registeredRedirectUrl, state, useOpenidConnect,
    );
    assert.deepEqual(claims, expectedClaims(authServerIssuer));
  });
});
