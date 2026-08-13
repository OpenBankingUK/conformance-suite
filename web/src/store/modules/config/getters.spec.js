import { getters } from './index';

describe('discoveryTemplates', () => {
  let state;
  const name = 'ob-v3.0-random';

  beforeEach(() => {
    state = {
      discoveryTemplates: [
        {
          model: {
            discoveryModel: { name },
          },
          image: 'imageData',
        },
      ],
    };
  });

  it('returns template with image set to PNG matching template name', async () => {
    const list = await getters.discoveryTemplates(state);
    expect(list[0].model).toEqual(state.discoveryTemplates[0].model);
    expect(list[0].image).toEqual('imageData');
  });
});

describe('discoveryProblems', () => {
  let state;

  beforeEach(() => {
    state = {
      problems: [
        {
          key: 'DiscoveryModel.DiscoveryItems[0].APISpecification.Name',
          error: 'Field validation for \'Name\' failed on the \'required\' tag',
        },
        {
          key: null,
          error: 'Unexpected token { in JSON at position 108',
        },
      ],
    };
  });

  it('returns object with JSON `path` and `parent` property for Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[0].path).toEqual('discoveryModel.discoveryItems[0].apiSpecification.name');
    expect(list[0].parent).toEqual('discoveryModel.discoveryItems[0].apiSpecification');
  });

  it('returns object with `error` property for Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[0].error).toEqual(`Field validation for 'Name' failed on the 'required' tag`); // eslint-disable-line
  });

  it('returns object with `error` property and null `path` for non Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[1].path).toEqual(null);
    expect(list[1].error).toEqual(`Unexpected token { in JSON at position 108`); // eslint-disable-line
  });

  it('returns null when null problems', () => {
    const list = getters.discoveryProblems({ problems: null });
    expect(list).toEqual(null);
  });
});

describe('paymentVersionState', () => {
  const discoveryModel = version => ({
    discoveryModel: {
      discoveryItems: [
        {
          apiSpecification: {
            name: 'Payment Initiation API',
            version,
            schemaVersion: `https://example.com/${version}/payment-initiation-openapi.json`,
          },
        },
      ],
    },
  });

  it('classifies v3-only payment discovery', () => {
    const state = { discoveryModel: discoveryModel('v3.1.6') };

    expect(getters.paymentVersionState(state)).toEqual('v3');
    expect(getters.showPaymentFrequency(state, getters)).toEqual(true);
    expect(getters.showV4StandingOrderFrequency(state, getters)).toEqual(false);
  });

  it('classifies v4-only payment discovery', () => {
    const state = { discoveryModel: discoveryModel('v4.0.1') };

    expect(getters.paymentVersionState(state)).toEqual('v4');
    expect(getters.showPaymentFrequency(state, getters)).toEqual(false);
    expect(getters.showV4StandingOrderFrequency(state, getters)).toEqual(true);
  });

  it('classifies mixed payment discovery', () => {
    const state = {
      discoveryModel: {
        discoveryModel: {
          discoveryItems: [
            discoveryModel('v3.1.6').discoveryModel.discoveryItems[0],
            discoveryModel('v4.0.1').discoveryModel.discoveryItems[0],
          ],
        },
      },
    };

    expect(getters.paymentVersionState(state)).toEqual('mixed');
    expect(getters.showPaymentFrequency(state, getters)).toEqual(true);
    expect(getters.showV4StandingOrderFrequency(state, getters)).toEqual(true);
    expect(getters.showFrequencyCompatibilityLabels(state, getters)).toEqual(true);
  });

  it('shows both controls when discovery is unavailable', () => {
    const state = { discoveryModel: null };

    expect(getters.paymentVersionState(state)).toEqual('unknown');
    expect(getters.showPaymentFrequency(state, getters)).toEqual(true);
    expect(getters.showV4StandingOrderFrequency(state, getters)).toEqual(true);
  });
});

describe('Config', () => {
  let state;

  beforeEach(() => {
    state = {
      discoveryModel: {
        discoveryModel: {
          tokenAcquisition: 'headless',
          callbackProxyUrl: 'https://callback-proxy.io',
        },
      },
    };
  });

  describe('getters', () => {
    it('discoveryModel', () => {
      expect(getters.discoveryModel(state)).toEqual(state.discoveryModel);
    });
    it('discoveryModelString', () => {
      expect(getters.discoveryModelString(state)).toEqual(JSON.stringify(state.discoveryModel, null, 2));
    });
    it('tokenAcquisition', () => {
      expect(getters.tokenAcquisition(state)).toEqual('headless');
    });
    it('callbackProxyUrl', () => {
      expect(getters.callbackProxyUrl(state)).toEqual('https://callback-proxy.io');
    });
  });
});
