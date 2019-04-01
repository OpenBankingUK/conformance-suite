import * as _ from 'lodash';
import * as types from './mutation-types';
import constants from '../config/constants';

import api from '../../../api';

/**
 * Setup WebSocket connection to the backend to retrieve results.
 */
const createWebSocketConnection = () => {
  // There are other ways of doing this, see:
  // https://vuex.vuejs.org/guide/plugins.html#committing-mutations-inside-plugins
  const getUrl = () => {
    const { location } = window;
    const { host, protocol } = location;
    const isSecure = protocol === 'https:';
    const scheme = isSecure ? 'wss' : 'ws';

    const url = `${scheme}://${host}/api/run/ws`;
    return url;
  };

  return new Promise((resolve, reject) => {
    const url = getUrl();

    const wsConnection = new WebSocket(url);
    wsConnection.onopen = () => {
      resolve(wsConnection);
    };
    wsConnection.onclose = (ev) => {
      reject(ev);
    };
  });
};

/*
Takes: [
  {
    specIdentifier: SPEC_NAME,
    namedPermissions: [
      {
        name: 'to1002',
        consentUrl: CONSENT_URL,
      },
      {
        name: 'to1002',
        consentUrl: CONSENT_URL2,
      },
    ],
  },
]
Returns: { [SPEC_NAME]: [CONSENT_URL, CONSENT_URL2] }
*/
const consentUrls = (specTokens) => {
  const map = {};
  specTokens.forEach((item) => {
    const list = item.namedPermissions.map(p => p.consentUrl);
    const urls = _.filter(list, u => u); // remove null/undefined urls
    map[item.specIdentifier] = urls;
  });
  return map;
};

export default {
  /**
   * Step 4: Calls /api/test-cases to get all the test cases, then sets the
   * retrieved test cases in the store.
   * Route: `/wizard/overview-run`.
   */
  async computeTestCases({ commit, dispatch, state }) {
    try {
      const setShowLoading = flag => dispatch('status/setShowLoading', flag, { root: true });
      await api.stopTestRun(setShowLoading); // ensure any previous run is stopped
      const testCases = await api.computeTestCases(setShowLoading);
      if (_.isEqual(testCases.specCases, state.testCases)) {
        return;
      }

      commit(types.SET_TEST_CASES, testCases.specCases);
      commit(types.SET_TEST_CASES_STATUS, '');

      if (testCases.specTokens) {
        commit(types.SET_CONSENT_URLS, consentUrls(testCases.specTokens));
      }
      dispatch('status/clearErrors', null, { root: true });

      // Setup WebSocket now as we need to know when tokens have been acquired.
      setShowLoading(true);
      const wsConnection = await createWebSocketConnection();
      commit(types.SET_WEBSOCKET_CONNECTION, wsConnection);

      wsConnection.onerror = (ev) => {
        setShowLoading(false);
        dispatch('status/setErrors', [ev], { root: true });
      };
      wsConnection.onmessage = (ev) => {
        // TODO: call `setShowLoading(false)` once using _.once or similar.
        setShowLoading(false);
        const { data } = ev;
        const update = JSON.parse(data);

        commit(types.SET_WEBSOCKET_MESSAGE, update);

        if (_.has(update, 'type') && update.type === 'ResultType_TestCasesCompleted') {
          commit(types.SET_TEST_CASES_COMPLETED, update.value);
        } else if (_.has(update, 'type') && update.type === 'ResultType_AcquiredAccessToken') {
          commit(types.ADD_TOKEN_ACQUIRED, update);
        } else if (_.has(update, 'type') && update.type === 'ResultType_AcquiredAllAccessTokens') {
          commit(types.SET_ALL_TOKENS_ACQUIRED);
        } else if (_.has(update, 'stopped') && update.stopped) {
          // do nothing
        } else {
          // update = {"error":"createRequest: setHeaders Replaced Context value Bearer $access_token :replacement not found in context: Bearer $access_token"}
          const isErrorMsg = _.has(update, 'error');
          if (!isErrorMsg) {
            // update = {"test":{"id":"#co0001","pass":true}}
            commit(types.UPDATE_TEST_CASE, update);
          }
        }
      };
    } catch (err) {
      commit(types.SET_TEST_CASES, []);
      dispatch('status/setErrors', [err], { root: true });
    }
    dispatch('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
  },
  /**
   * Step 5: Calls POST `/api/run` then setups WebSocket connection to `/api/run/ws`.
   * Route: `/wizard/overview-run`.
   */
  async executeTestCases({ commit, dispatch }) {
    const setShowLoading = flag => dispatch('status/setShowLoading', flag, { root: true });
    try {
      commit(types.SET_HAS_RUN_STARTED, true);

      await api.executeTestCases(setShowLoading);
      dispatch('status/clearErrors', null, { root: true });
      commit(types.SET_TEST_CASES_STATUS, 'PENDING');
    } catch (err) {
      setShowLoading(false);
      commit(types.SET_HAS_RUN_STARTED, false);
      dispatch('status/setErrors', [err], { root: true });
    }
    dispatch('config/setWizardStep', constants.WIZARD.STEP_FIVE, { root: true });
  },
};
