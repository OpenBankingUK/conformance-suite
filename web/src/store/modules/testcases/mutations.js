import * as _ from 'lodash';
import * as moment from 'moment';
import * as types from './mutation-types';

export default {
  [types.SET_TEST_CASES](state, testCases) {
    state.testCases = testCases;
  },
  [types.SET_HAS_RUN_STARTED](state, hasRunStarted) {
    state.hasRunStarted = hasRunStarted;
  },
  [types.SET_WEBSOCKET_CONNECTION](state, connection) {
    state.ws.connection = connection;
  },
  [types.SET_WEBSOCKET_MESSAGE](state, message) {
    state.ws.messages = [
      ...state.ws.messages,
      message,
    ];
  },
  [types.UPDATE_TEST_CASE](state, update) {
    const { id, pass, metrics } = update.test;
    const predicate = { '@id': id };

    // Assume that each testCase has a globally unique id, then find the matching testCase.
    const testCases = _.flatMap(state.testCases, spec => _.get(spec, 'testCases'));
    const testCase = _.find(testCases, predicate);

    if (testCase) {
      testCase.meta.status = pass ? 'PASSED' : 'FAILED';
      const responseSeconds = moment.duration(metrics.response_time / 1000000).asSeconds().toFixed(6);
      testCase.meta.metrics.responseTime = `${responseSeconds}s`;
      testCase.meta.metrics.responseSize = `${metrics.response_size}B`;
    } else {
      // eslint-disable-next-line no-console
      console.error('Failed to find testCase, testCases=%o, predicate=%o, update=%o', testCases, predicate, update);
    }
  },
  [types.SET_TEST_CASES_STATUS](state, status) {
    const DEFAULTS = { meta: { status, metrics: { responseTime: '', responseSize: '' } } };

    state.testCases = _.map(state.testCases, (spec) => {
      const testCases = _.map(spec.testCases, testCase => _.merge({}, testCase, DEFAULTS));
      return _.merge({}, spec, { testCases });
    });
  },
};
