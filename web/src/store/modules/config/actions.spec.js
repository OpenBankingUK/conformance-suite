import actions from './actions';
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
      const state = {
        discoveryModel: {},
      };
      actions.setDiscoveryModel({ commit, state }, '{');
      expect(commit).toHaveBeenCalledWith('DISCOVERY_MODEL_PROBLEMS', [{
        error: 'Unexpected end of JSON input',
        key: null,
      }]);
    });

    it('does not commit discovery model', () => {
      const state = {
        discoveryModel: {},
      };
      actions.setDiscoveryModel({ commit, state }, '{');
      expect(commit).not.toHaveBeenCalledWith('SET_DISCOVERY_MODEL', '{');
    });
  });

  describe('with valid JSON string', () => {
    it('commits parsed JSON', () => {
      const state = {
        discoveryModel: {},
      };
      actions.setDiscoveryModel({ commit, state }, '{"a": 1}');
      expect(commit).toHaveBeenCalledWith('SET_DISCOVERY_MODEL', { a: 1 });
    });

    it('commits null problems', () => {
      const state = {
        discoveryModel: {},
      };
      actions.setDiscoveryModel({ commit, state }, '{"a": 1}');
      expect(commit).toHaveBeenCalledWith('DISCOVERY_MODEL_PROBLEMS', null);
    });
  });
});
