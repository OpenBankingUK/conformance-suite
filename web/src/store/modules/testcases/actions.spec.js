import actions from './actions';
import * as types from './mutation-types';

import constants from '../config/constants';
import api from '../../../api';

jest.mock('../../../api');

describe('executeTestCases', () => {
  const state = { execution: {} };
  let commit;
  let dispatch;

  describe('when execution sucessful', () => {
    const result = 'mock';
    beforeEach(() => {
      commit = jest.fn();
      dispatch = jest.fn();
      api.executeTestCases.mockResolvedValue(result);
    });

    it('commits testcases', async () => {
      await actions.executeTestCases({ commit, dispatch, state });
      expect(commit).toHaveBeenCalledWith(types.SET_EXECUTION_RESULTS, result);
      expect(dispatch).toHaveBeenCalledWith('config/setExecutionErrors', [], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FIVE, { root: true });
    });
  });

  describe('when execution throws Error', () => {
    let error;

    beforeEach(() => {
      commit = jest.fn();
      error = new Error('some error');
      api.executeTestCases.mockRejectedValue(error);
    });

    it('sets Error', async () => {
      await actions.executeTestCases({ commit, dispatch, state });
      expect(commit).not.toHaveBeenCalled();
      expect(dispatch).toHaveBeenCalledWith('config/setExecutionErrors', [error], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FIVE, { root: true });
    });
  });
});
