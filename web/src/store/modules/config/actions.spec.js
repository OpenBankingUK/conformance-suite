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
