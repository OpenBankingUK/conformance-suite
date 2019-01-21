import * as _ from 'lodash';
import * as types from './mutation-types';
import constants from '../config/constants';

import api from '../../../api';

export default {
  /**
   * Step 4: Calls /api/test-cases to get all the test cases, then sets the
   * retrieved test cases in the store.
   * Route: `/wizard/run-overview`.
   */
  async computeTestCases({ commit, dispatch, state }) {
    try {
      const testCases = await api.computeTestCases();
      if (_.isEqual(testCases, state.testCases)) {
        return;
      }

      commit(types.SET_TEST_CASES, testCases);
      dispatch('config/setTestCaseErrors', [], { root: true });
    } catch (err) {
      commit(types.SET_TEST_CASES, []);
      dispatch('config/setTestCaseErrors', [err], { root: true });
    }
    dispatch('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
  },
  /**
   * Step 5: Calls `/api/run/start`.
   * Route: `/wizard/run-overview`.
   */
  async executeTestCases({ commit, dispatch }) {
    try {
      const execution = await api.executeTestCases();
      commit(types.SET_EXECUTION_RESULTS, execution);
      dispatch('config/setExecutionErrors', [], { root: true });
    } catch (err) {
      dispatch('config/setExecutionErrors', [err], { root: true });
    }
    dispatch('config/setWizardStep', constants.WIZARD.STEP_FIVE, { root: true });
  },
};
