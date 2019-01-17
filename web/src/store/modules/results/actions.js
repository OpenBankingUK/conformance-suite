import * as types from './mutation-types';
import constants from '../config/constants';

import api from '../../../api';

export default {
  /**
   * Calls /api/report to get all the test cases, then sets the
   * retrieved test cases in the store.
   */
  async computeTestCaseResults({ commit, dispatch }) {
    try {
      const testCaseResults = await api.computeTestCaseResults();
      commit(types.SET_TEST_CASE_RESULTS, testCaseResults);
      dispatch('config/setTestCaseResultsErrors', []);
      dispatch('config/setWizardStep', constants.WIZARD.STEP_SIX);
    } catch (err) {
      commit(types.SET_TEST_CASE_RESULTS, {});
      dispatch('config/setTestCaseResultsErrors', [
        err,
      ]);
    }
  },
};
