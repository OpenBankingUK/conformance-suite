import * as _ from 'lodash';
import * as moment from 'moment';
import * as types from './mutation-types';

export default {
  [types.SET_TEST_CASES](state, testCases) {
    state.testCases = testCases;
  },
  [types.SET_TEST_CASES_COMPLETED](state, value) {
    state.test_cases_completed = value;
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
    const predicate = { '@id': update.test.id };
    // Assume that each testCase has a globally unique id, then find the matching testCase.
    const testCases = _.flatMap(state.testCases, spec => _.get(spec, 'testCases'));
    const testCase = _.find(testCases, predicate);

    if (_.isNil(testCase)) {
      // eslint-disable-next-line no-console
      console.error('Failed to find testCase, testCases=%o, predicate=%o, update=%o', testCases, predicate, update);
      return;
    }

    const {
      id, pass, metrics, fail, detail, refURI,
    } = update.test;

    testCase.id = id;
    testCase.meta.status = pass ? 'PASSED' : 'FAILED';
    const responseSeconds = moment.duration(metrics.response_time).asMilliseconds().toFixed(3);
    testCase.meta.metrics.responseTime = `${responseSeconds.toLocaleString()}ms`;
    testCase.meta.metrics.responseSize = `${metrics.response_size.toLocaleString()}`;
    testCase.error = fail;
    testCase.detail = detail;
    testCase.refURI = refURI;

    if (fail) {
      // Set the row variant, for alternate styling.
      // https://bootstrap-vue.js.org/docs/components/table/#items-record-data-
      _.merge(testCase, {
        _rowVariant: 'danger',
      });
    }
  },
  [types.SET_TEST_CASES_STATUS](state, status) {
    const DEFAULTS = {
      _rowVariant: null,
      _showDetails: false,
      meta: {
        status,
        metrics: {
          responseTime: '',
          responseSize: '',
        },
      },
    };

    state.testCases = _.map(state.testCases, (spec) => {
      const testCases = _.map(spec.testCases, testCase => _.merge({}, testCase, DEFAULTS));
      return _.assign(spec, { testCases });
    });
  },
  [types.SET_CONSENT_URLS](state, urls) {
    state.consentUrls = urls;
  },
  [types.TOGGLE_ROW_DETAILS](state, item) {
    _.merge(item, {
      _showDetails: !_.get(item, '_showDetails'),
    });
  },

  [types.ADD_TOKEN_ACQUIRED](state, value) {
    state.tokens.acquired = [
      ...state.tokens.acquired,
      value,
    ];
  },
  [types.SET_ALL_TOKENS_ACQUIRED](state) {
    state.tokens.all_acquired = true;
  },
};
