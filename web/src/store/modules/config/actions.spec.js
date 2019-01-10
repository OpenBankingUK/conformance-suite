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

[
  {
    action: 'setDiscoveryModel',
    property: 'discoveryModel',
    successMutation: 'SET_DISCOVERY_MODEL',
    errorMutation: 'DISCOVERY_MODEL_PROBLEMS',
  },
].forEach((scenario) => {
  describe(scenario.action, () => {
    const state = {
      [scenario.property]: {},
    };
    let commit;
    beforeEach(() => {
      commit = jest.fn();
    });

    describe('with invalid JSON string', () => {
      it('commits problems', () => {
        actions[scenario.action]({ commit, state }, '{');
        expect(commit).toHaveBeenCalledWith(scenario.errorMutation, [{
          error: 'Unexpected end of JSON input',
          key: null,
        }]);
      });

      it('does not commit value', () => {
        actions[scenario.action]({ commit, state }, '{');
        expect(commit).not.toHaveBeenCalledWith(scenario.successMutation, '{');
      });
    });

    describe('with valid JSON string', () => {
      it('commits parsed JSON', () => {
        actions[scenario.action]({ commit, state }, '{"a": 1}');
        expect(commit).toHaveBeenCalledWith(scenario.successMutation, { a: 1 });
      });

      it('commits null problems', () => {
        actions[scenario.action]({ commit, state }, '{"a": 1}');
        expect(commit).toHaveBeenCalledWith(scenario.errorMutation, null);
      });
    });
  });
});
