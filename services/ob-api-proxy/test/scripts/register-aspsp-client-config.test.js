const assert = require('assert');
const proxyquire = require('proxyquire');
const sinon = require('sinon');

describe('registerAgreedConfig', () => {
  let fakeUpdateRegisteredConfig;
  let registerAgreedConfig;

  beforeEach(async () => {
    fakeUpdateRegisteredConfig = sinon.stub();
    ({ registerAgreedConfig } = proxyquire(
      '../../scripts/register-aspsp-client-config',
      {
        '../app/authorisation-servers': {
          updateRegisteredConfig: fakeUpdateRegisteredConfig,
        },
      },
    ));
  });

  it('registers config agreed with ASPSP for an authServerId and scoped by software statement', async () => {
    await registerAgreedConfig(['authServerId=48fr7qwRKzA0eWKR2Se8YR', 'field=request_object_signing_alg', 'value=["PS256"]']);
    assert(fakeUpdateRegisteredConfig.calledWithExactly(
      '48fr7qwRKzA0eWKR2Se8YR',
      { request_object_signing_alg: ['PS256'] },
    ));
  });
});
