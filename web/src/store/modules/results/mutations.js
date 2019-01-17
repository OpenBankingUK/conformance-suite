import * as types from './mutation-types';

export default {
  [types.SET_TEST_CASE_RESULTS](state, testCaseResults) {
    state.testCaseResults = testCaseResults;
  },
};
