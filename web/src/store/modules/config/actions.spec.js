import actions from './actions';
import getters from './getters';
import router from '../../../router';
import * as types from './mutation-types';

import discovery from '../../../api/discovery';

jest.mock('../../../api/discovery');

describe('validateDiscoveryConfig', () => {
  const state = { discoveryModel: {} };
  let commit;

  describe('when validation passes', () => {
    beforeEach(() => {
      commit = jest.fn();
      discovery.validateDiscoveryConfig.mockResolvedValue({
        success: true,
        problems: [],
      });
    });
    it('commits null validation problems', async () => {
      await actions.validateDiscoveryConfig({ commit, state });
      expect(commit).toHaveBeenCalledWith(types.DISCOVERY_MODEL_PROBLEMS, null);
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
      discovery.validateDiscoveryConfig.mockResolvedValue({
        success: false,
        problems,
      });
    });
    it('commits array of validation problem strings', async () => {
      await actions.validateDiscoveryConfig({ commit, state });
      expect(commit).toHaveBeenCalledWith(types.DISCOVERY_MODEL_PROBLEMS, problems);
    });
  });
  describe('when validation throws Error', () => {
    beforeEach(() => {
      commit = jest.fn();
      discovery.validateDiscoveryConfig.mockRejectedValue(new Error('some error'));
    });
    it('commits Error message in problems array', async () => {
      await actions.validateDiscoveryConfig({ commit, state });
      expect(commit).toHaveBeenCalledWith(types.DISCOVERY_MODEL_PROBLEMS, [
        { key: null, error: 'some error' },
      ]);
    });
  });
});

describe('setDiscoveryModel', () => {
  let commit;
  beforeEach(() => {
    commit = jest.fn();
  });

  describe('with invalid JSON string', () => {
    it('commits problems', () => {
      actions.setDiscoveryModel({ commit }, '{');
      expect(commit).toHaveBeenCalledWith('DISCOVERY_MODEL_PROBLEMS', ['Unexpected end of JSON input']);
    });
    it('does not commit discovery model', () => {
      actions.setDiscoveryModel({ commit }, '{');
      expect(commit).not.toHaveBeenCalledWith('SET_DISCOVERY_MODEL', '{');
    });
  });
  describe('with valid JSON string', () => {
    it('commits parsed JSON', () => {
      actions.setDiscoveryModel({ commit }, '{"a": 1}');
      expect(commit).toHaveBeenCalledWith('SET_DISCOVERY_MODEL', { a: 1 });
    });
    it('commits null problems', () => {
      actions.setDiscoveryModel({ commit }, '{"a": 1}');
      expect(commit).toHaveBeenCalledWith('DISCOVERY_MODEL_PROBLEMS', null);
    });
  });
});

describe('Config', () => {
  describe('actions', () => {
    let dispatch;
    let commit;
    let routerSpy;

    beforeEach(() => {
      dispatch = jest.fn();
      commit = jest.fn();
      routerSpy = jest.spyOn(router, 'push');
    });

    it('setConfig', () => {
      const config = '{"a": 1}';
      actions.setConfig({ commit }, config);
      expect(commit).toHaveBeenCalledWith('SET_CONFIG', config);
    });

    it('resetValidationsRun', () => {
      actions.resetValidationsRun({ commit });
      expect(commit).toHaveBeenCalledWith('reporter/SET_WEBSOCKET_LAST_UPDATE', null, { root: true });
      expect(commit).toHaveBeenCalledWith('validations/SET_VALIDATION_DISCOVERY_MODEL', null, { root: true });
    });

    it('startValidation', () => {
      actions.startValidation({ dispatch, getters });
      expect(dispatch).toHaveBeenNthCalledWith(1, 'resetValidationsRun');
      expect(dispatch).toHaveBeenNthCalledWith(
        2,
        'validations/validate', {
          discoveryModel: getters.getDiscoveryModel,
          config: getters.getConfig,
        },
        { root: true },
      );
      expect(routerSpy).toHaveBeenCalledWith('/reports');
    });

    it('updateDiscoveryModel', () => {
      const discoveryModel = '{"a": 1}';
      actions.updateDiscoveryModel({ commit }, discoveryModel);
      expect(commit).toHaveBeenCalledWith('UPDATE_DISCOVERY_MODEL', discoveryModel);
    });

    it('deleteDiscoveryModel', () => {
      const discoveryModel = '{"a": 1}';
      actions.deleteDiscoveryModel({ commit }, discoveryModel);
      expect(commit).toHaveBeenCalledWith('DELETE_DISCOVERY_MODEL', discoveryModel);
    });
  });
});
