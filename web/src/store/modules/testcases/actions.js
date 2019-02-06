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

export default {
  /**
   * Step 4: Calls /api/test-cases to get all the test cases, then sets the
   * retrieved test cases in the store.
   * Route: `/wizard/overview-run`.
   */
  async computeTestCases({ commit, dispatch, state }) {
    try {
      const testCases = await api.computeTestCases();
      if (_.isEqual(testCases, state.testCases)) {
        return;
      }

      commit(types.SET_TEST_CASES, testCases);
      commit(types.SET_TEST_CASES_STATUS, '');
      dispatch('status/clearErrors', null, { root: true });
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
    try {
      commit(types.SET_HAS_RUN_STARTED, true);

      await api.executeTestCases();
      dispatch('status/clearErrors', null, { root: true });
      commit(types.SET_TEST_CASES_STATUS, 'PENDING');

      const wsConnection = await createWebSocketConnection();
      commit(types.SET_WEBSOCKET_CONNECTION, wsConnection);

      wsConnection.onerror = (ev) => {
        dispatch('status/setErrors', [ev], { root: true });
      };
      wsConnection.onmessage = (ev) => {
        const { data } = ev;
        const update = JSON.parse(data);

        commit(types.SET_WEBSOCKET_MESSAGE, update);

        const isErrorMsg = _.has(update, 'error'); // update = {"error":"createRequest: setHeaders Replaced Context value Bearer $access_token :replacement not found in context: Bearer $access_token"}
        if (!isErrorMsg) {
          commit(types.UPDATE_TEST_CASE, update); // update = {"test":{"id":"#co0001","pass":true}}
        }
      };
    } catch (err) {
      commit(types.SET_HAS_RUN_STARTED, false);
      dispatch('status/setErrors', [err], { root: true });
    }
    dispatch('config/setWizardStep', constants.WIZARD.STEP_FIVE, { root: true });
  },
};
