const assert = require('assert'); // eslint-disable-line
const proxyquire = require('proxyquire'); // eslint-disable-line
const sinon = require('sinon'); //eslint-disable-line

describe('cacheLatestConfigs', () => {
  let fakeFetchOBAccountPaymentServiceProviders;
  let fakeUpdateOpenIdConfigs;

  beforeEach(async () => {
    fakeFetchOBAccountPaymentServiceProviders = sinon.stub().returns([]);
    fakeUpdateOpenIdConfigs = sinon.stub().returns(null);
    const cacheLatestConfigs = proxyquire('../../scripts/update-auth-server-and-open-id-configs', // eslint-disable-line
      {
        '../app/authorisation-servers': {
          updateOpenIdConfigs: fakeUpdateOpenIdConfigs,
        },
        '../app/ob-directory': {
          fetchOBAccountPaymentServiceProviders: fakeFetchOBAccountPaymentServiceProviders,
        },
      },
    ).cacheLatestConfigs;

    cacheLatestConfigs();
  });

  it('fetches OB account service providers', () => {
    assert(fakeFetchOBAccountPaymentServiceProviders.called);
  });

  it('updates id configs', () => {
    assert(fakeUpdateOpenIdConfigs.called);
  });
});
