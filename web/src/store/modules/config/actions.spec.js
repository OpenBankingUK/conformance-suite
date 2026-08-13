import api from '../../../api';
import actions from './actions';
import * as types from './mutation-types.js';

jest.mock('../../../api');

describe('setDiscoveryTemplates', () => {
  let commit;

  beforeEach(() => {
    commit = jest.fn();
  });

  const name = 'ob-v3.0-random';

  it('commits discovery template with matched image', async () => {
    const matchingImage = `./${name}.png`;
    const discoveryImages = {};
    discoveryImages[matchingImage] = 'mockImage';
    const discoveryTemplates = [{ discoveryModel: { name } }];

    const data = { discoveryTemplates, discoveryImages };
    await actions.setDiscoveryTemplates({ commit }, data);
    expect(commit).toHaveBeenCalledWith(types.SET_DISCOVERY_TEMPLATES, [
      {
        model: discoveryTemplates[0],
        image: 'mockImage',
      },
    ]);
  });

  it('commits discovery template with no-image default when no matching image', async () => {
    const nonMatchingImage = './an-image.png';
    const defaultImage = './no-image-discovery-icon.png';
    const discoveryImages = {};
    discoveryImages[nonMatchingImage] = 'mockImage';
    discoveryImages[defaultImage] = 'defaultImage';
    const discoveryTemplates = [{ discoveryModel: { name } }];

    const data = { discoveryTemplates, discoveryImages };
    await actions.setDiscoveryTemplates({ commit }, data);
    expect(commit).toHaveBeenCalledWith(types.SET_DISCOVERY_TEMPLATES, [
      {
        model: discoveryTemplates[0],
        image: 'defaultImage',
      },
    ]);
  });

  it('commits discovery template with v4.0 fallback image for v4.0.1 template', async () => {
    const fallbackImage = './ob-v4.0-generic.png';
    const discoveryImages = {};
    discoveryImages[fallbackImage] = 'v4FallbackImage';
    discoveryImages['./no-image-discovery-icon.png'] = 'defaultImage';
    const discoveryTemplates = [{ discoveryModel: { name: 'ob-v4.0.1-generic' } }];

    const data = { discoveryTemplates, discoveryImages };
    await actions.setDiscoveryTemplates({ commit }, data);
    expect(commit).toHaveBeenCalledWith(types.SET_DISCOVERY_TEMPLATES, [
      {
        model: discoveryTemplates[0],
        image: 'v4FallbackImage',
      },
    ]);
  });

  it('commits discovery template with ozone fallback image for v4.0.1 mobile template', async () => {
    const fallbackImage = './ob-v4.0-ozone.png';
    const discoveryImages = {};
    discoveryImages[fallbackImage] = 'v4OzoneFallbackImage';
    discoveryImages['./no-image-discovery-icon.png'] = 'defaultImage';
    const discoveryTemplates = [{ discoveryModel: { name: 'ob-v4.0.1-ozone-mobile' } }];

    const data = { discoveryTemplates, discoveryImages };
    await actions.setDiscoveryTemplates({ commit }, data);
    expect(commit).toHaveBeenCalledWith(types.SET_DISCOVERY_TEMPLATES, [
      {
        model: discoveryTemplates[0],
        image: 'v4OzoneFallbackImage',
      },
    ]);
  });
});

describe('validateDiscoveryConfig', () => {
  const state = { discoveryModel: {} };
  let commit;
  let dispatch;

  describe('when validation passes', () => {
    beforeEach(() => {
      commit = jest.fn();
      dispatch = jest.fn();
      api.validateDiscoveryConfig.mockReturnValueOnce({
        success: true,
        problems: [],
        response: {
          token_endpoints: {
            'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json': 'https://modelobank2018.o3bank.co.uk:4201/token',
            'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/payment-initiation-swagger.json': 'https://modelobank2018.o3bank.co.uk:4201/token',
          },
          authorization_endpoints: {
            'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json': 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
            'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/payment-initiation-swagger.json': 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
          },
        },
      });
    });

    afterEach(() => {
      jest.resetAllMocks();
    });

    it('commits null validation problems', async () => {
      await actions.validateDiscoveryConfig({ commit, dispatch, state });
      expect(commit).toHaveBeenCalledWith(types.DISCOVERY_MODEL_PROBLEMS, null);
    });

    it('dispatches clearErrors', async () => {
      await actions.validateDiscoveryConfig({ commit, dispatch, state });
      expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
    });
  });

  describe('when validation fails with problem messages', () => {
    const problems = [
      {
        key: 'DiscoveryModel.Version',
        error: 'Field validation for \'Version\' failed on the \'required\' tag',
      },
      {
        key: 'DiscoveryModel.DiscoveryItems',
        error: 'Field validation for \'DiscoveryItems\' failed on the \'required\' tag',
      },
    ];

    beforeEach(() => {
      commit = jest.fn();
      dispatch = jest.fn();
      api.validateDiscoveryConfig.mockReturnValueOnce({
        success: false,
        problems,
      });
    });

    afterEach(() => {
      jest.resetAllMocks();
    });

    it('commits array of validation problem strings', async () => {
      await actions.validateDiscoveryConfig({ commit, dispatch, state });
      expect(commit).toHaveBeenCalledWith(types.DISCOVERY_MODEL_PROBLEMS, problems);
    });

    it('dispatches setErrors', async () => {
      await actions.validateDiscoveryConfig({ commit, dispatch, state });
      const expected = [problems[0].error, problems[1].error];
      expect(dispatch).toHaveBeenCalledWith('status/setErrors', expected, { root: true });
    });
  });

  describe('when validation throws Error', () => {
    beforeEach(() => {
      commit = jest.fn();
      dispatch = jest.fn();
      api.validateDiscoveryConfig.mockRejectedValueOnce(new Error('some error'));
    });

    afterEach(() => {
      jest.resetAllMocks();
    });

    it('commits Error message in problems array', async () => {
      await actions.validateDiscoveryConfig({ commit, dispatch, state });
      expect(commit).toHaveBeenCalledWith(types.DISCOVERY_MODEL_PROBLEMS, [
        { key: null, error: 'some error' },
      ]);
    });

    it('dispatches setErrors', async () => {
      await actions.validateDiscoveryConfig({ commit, dispatch, state });
      expect(dispatch).toHaveBeenCalledWith('status/setErrors', ['some error'], { root: true });
    });
  });
});

describe('v4 standing order frequency configuration', () => {
  let commit;
  let dispatch;

  beforeEach(() => {
    commit = jest.fn();
    dispatch = jest.fn();
    api.validateConfiguration.mockReset();
  });

  const validConfigurationState = overrides => ({
    configuration: {
      signing_private: 'signing-private',
      signing_public: 'signing-public',
      transport_private: 'transport-private',
      transport_public: 'transport-public',
      resource_ids: {
        account_ids: [{ account_id: 'account-id' }],
        statement_ids: [{ statement_id: 'statement-id' }],
      },
      transaction_from_date: '2022-01-01T00:00:00+01:00',
      transaction_to_date: '2022-01-01T00:00:00+01:00',
      client_id: 'client-id',
      client_secret: 'client-secret',
      token_endpoint: 'https://server/token',
      response_type: 'code id_token',
      token_endpoint_auth_method: 'client_secret_basic',
      request_object_signing_alg: 'PS256',
      authorization_endpoint: 'https://server/auth',
      resource_base_url: 'https://server',
      x_fapi_financial_id: 'financial-id',
      issuer: 'https://server/issuer',
      redirect_url: 'https://server/callback',
      payment_frequency: 'EvryDay',
      v4_standing_order_frequency: {
        Type: 'WEEK',
        PointInTime: '03',
      },
      ...overrides,
    },
  });

  it('imports a Type-only v4 standing order frequency without merging the default PointInTime', () => {
    const state = {
      configuration: {
        payment_frequency: 'EvryDay',
        v4_standing_order_frequency: {
          Type: 'WEEK',
          PointInTime: '03',
        },
      },
    };

    actions.setConfigurationJSON({ commit, dispatch, state }, '{"v4_standing_order_frequency":{"Type":"MNTH"}}');

    expect(commit).toHaveBeenCalledWith(types.SET_CONFIGURATION, {
      payment_frequency: 'EvryDay',
      v4_standing_order_frequency: {
        Type: 'MNTH',
      },
    });
    expect(commit).toHaveBeenCalledWith(types.SET_V4_STANDING_ORDER_FREQUENCY, {
      Type: 'MNTH',
    });
  });

  it('rejects v4 standing order frequency with PointInTime and CountPerPeriod together', async () => {
    const state = validConfigurationState({
      v4_standing_order_frequency: {
        Type: 'WEEK',
        PointInTime: '03',
        CountPerPeriod: 1,
      },
    });

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(false);
    expect(dispatch).toHaveBeenCalledWith('status/setErrors', [
      'v4_standing_order_frequency.PointInTime and CountPerPeriod are mutually exclusive',
    ], { root: true });
  });

  it.each(['1', '03', '99', '-1'])('accepts spec-aligned PointInTime value %s', async (pointInTime) => {
    const state = validConfigurationState({
      v4_standing_order_frequency: {
        Type: 'WEEK',
        PointInTime: pointInTime,
      },
    });
    api.validateConfiguration.mockResolvedValueOnce({});

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(true);
  });

  it.each(['100', '-10', 'AA'])('rejects invalid PointInTime value %s', async (pointInTime) => {
    const state = validConfigurationState({
      v4_standing_order_frequency: {
        Type: 'WEEK',
        PointInTime: pointInTime,
      },
    });

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(false);
    expect(dispatch).toHaveBeenCalledWith('status/setErrors', [
      'v4_standing_order_frequency.PointInTime must be numeric text up to two characters, including negative single-digit values',
    ], { root: true });
  });

  it.each([1, 2147483647])('accepts positive Int32 CountPerPeriod value %d', async (countPerPeriod) => {
    const state = validConfigurationState({
      v4_standing_order_frequency: {
        Type: 'WEEK',
        CountPerPeriod: countPerPeriod,
      },
    });
    api.validateConfiguration.mockResolvedValueOnce({});

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(true);
  });

  it.each([0, -1, 1.5, 2147483648])('rejects invalid CountPerPeriod value %d', async (countPerPeriod) => {
    const state = validConfigurationState({
      v4_standing_order_frequency: {
        Type: 'WEEK',
        CountPerPeriod: countPerPeriod,
      },
    });

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(false);
    expect(dispatch).toHaveBeenCalledWith('status/setErrors', [
      'v4_standing_order_frequency.CountPerPeriod must be a positive Int32',
    ], { root: true });
  });

  it('rejects invalid payment_frequency values', async () => {
    const state = validConfigurationState({
      payment_frequency: 'WEEK',
    });

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(false);
    expect(dispatch).toHaveBeenCalledWith('status/setErrors', [
      'Payment frequency must be a valid Frequency_1 value',
    ], { root: true });
  });

  it('rejects legacy v4 Frequency_1 Type combined with structured fields', async () => {
    const state = validConfigurationState({
      v4_standing_order_frequency: {
        Type: 'IntrvlWkDay:01:03',
        PointInTime: '03',
      },
    });

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(false);
    expect(dispatch).toHaveBeenCalledWith('status/setErrors', [
      'v4_standing_order_frequency legacy Frequency_1 Type cannot be combined with PointInTime or CountPerPeriod',
    ], { root: true });
  });

  it('accepts legacy v4 Frequency_1 Type without structured fields', async () => {
    const state = validConfigurationState({
      v4_standing_order_frequency: {
        Type: 'IntrvlWkDay:01:03',
      },
    });
    api.validateConfiguration.mockResolvedValueOnce({});

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(true);
  });

  it('accepts new OBFrequency6Code values', async () => {
    const state = validConfigurationState({
      v4_standing_order_frequency: {
        Type: 'LXMH',
      },
    });
    api.validateConfiguration.mockResolvedValueOnce({});

    const valid = await actions.validateConfiguration({ commit, dispatch, state });

    expect(valid).toEqual(true);
  });
});

[
  {
    action: 'setDiscoveryModel',
    property: 'discoveryModel',
    successMutation: types.SET_DISCOVERY_MODEL,
    errorMutation: types.DISCOVERY_MODEL_PROBLEMS,
    expectedErrorState: [{
      error: 'Unexpected end of JSON input',
      key: null,
    }],
    expectedErrors: ['Unexpected end of JSON input'],
    initialState: {},
    validJSON: '{"a": 1}',
    expectedState: { a: 1 },
  },
  {
    action: 'setConfigurationJSON',
    property: 'configuration',
    successMutation: types.SET_CONFIGURATION,
    errorMutation: null,
    expectedErrors: ['Unexpected end of JSON input'],
    initialState: {
      signing_private: '',
      signing_public: '',
      transport_private: '',
      transport_public: '',
    },
    validJSON: '{"a": 1, "signing_private": "test"}',
    expectedState: {
      signing_private: 'test',
      signing_public: '',
      transport_private: '',
      transport_public: '',
    },
  },
].forEach(({
  action, property, successMutation, errorMutation, expectedErrors, initialState,
  validJSON, expectedState, expectedErrorState,
}) => {
  describe(action, () => {
    const state = {
      [property]: initialState,
    };
    let commit;
    let dispatch;
    beforeEach(() => {
      commit = jest.fn();
      dispatch = jest.fn();
    });

    describe('with JSON string equal to current state', () => {
      it('does not commit value', () => {
        actions[action]({ commit, dispatch, state }, '{}');
        expect(commit).not.toHaveBeenCalledWith(successMutation, '{}');
      });
    });

    describe('with invalid JSON string', () => {
      it('commits problems', () => {
        actions[action]({ commit, dispatch, state }, '{');
        if (errorMutation) {
          expect(commit).toHaveBeenCalledWith(errorMutation, expectedErrorState);
        }
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', expectedErrors, { root: true });
      });

      it('does not commit value', () => {
        actions[action]({ commit, dispatch, state }, '{');
        expect(commit).not.toHaveBeenCalledWith(successMutation, '{');
      });
    });

    describe('with valid JSON string', () => {
      it('commits parsed JSON', () => {
        actions[action]({ commit, dispatch, state }, validJSON);
        expect(commit).toHaveBeenCalledWith(successMutation, expectedState);
      });

      it('commits null problems', () => {
        actions[action]({ commit, dispatch, state }, validJSON);
        if (errorMutation) {
          expect(commit).toHaveBeenCalledWith(errorMutation, null);
        }
        expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
      });
    });
  });
});
