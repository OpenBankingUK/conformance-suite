import * as _ from 'lodash';
import * as types from './mutation-types';
import * as configTypes from '../config/mutation-types';
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
      dispatch(configTypes.SET_TEST_CASES_ERROR, []);
      dispatch(configTypes.SET_WIZARD_STEP, constants.WIZARD.STEP_FOUR);
    } catch (err) {
      commit(types.SET_TEST_CASES, []);
      dispatch(configTypes.SET_TEST_CASES_ERROR, [err]);
      dispatch(configTypes.SET_WIZARD_STEP, constants.WIZARD.STEP_FOUR);
    }
  },
  /**
   * Step 5: Calls `/api/run/start`.
   * Route: `/wizard/run-overview`.
   */
  async executeTestCases({ commit, dispatch }) {
    try {
      const execution = await api.executeTestCases();
      commit(types.SET_EXECUTION_RESULTS, execution);
      dispatch(configTypes.SET_EXECUTION_ERROR, []);
      dispatch(configTypes.SET_WIZARD_STEP, constants.WIZARD.STEP_FIVE);
    } catch (err) {
      dispatch(configTypes.SET_EXECUTION_ERROR, [err]);
      dispatch(configTypes.SET_WIZARD_STEP, constants.WIZARD.STEP_FIVE);
    }
  },
};
