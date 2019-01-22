import * as types from './mutation-types';

export default {
  [types.SET_EXECUTION_RESULTS](state, execution) {
    state.execution = execution;
  },
  [types.SET_TEST_CASES](state, testCases) {
    state.testCases = testCases;
  },
  [types.SET_HAS_RUN_STARTED](state, hasRunStarted) {
    state.hasRunStarted = hasRunStarted;
  },
};
