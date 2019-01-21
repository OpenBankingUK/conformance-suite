/**
 * This creates a real store so avoid having to mock things.
 * This makes testing much easier.
 *
 * See the recommendation:
 * https://vue-test-utils.vuejs.org/guides/using-with-vuex.html#testing-a-running-store
 */
import { createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import { cloneDeep } from 'lodash';

import actions from './actions';
import mutations from './mutations';
import getters from './getters';
import state from './state';

import constants from '../config/constants';
import api from '../../../api';
// https://jestjs.io/docs/en/mock-functions#mocking-modules
jest.mock('../../../api');

describe('store/modules/results', () => {
  let dispatch;
  /**
     * Creates a real store so we don't have to mock things out.
     */
  const createRealStore = () => {
    const localVue = createLocalVue();
    localVue.use(Vuex);

    const store = new Vuex.Store({
      state: cloneDeep(state),
      actions,
      mutations,
      getters,
    });
    dispatch = jest.fn();
    store.dispatch = dispatch;
    return store;
  };

  it('config/testCaseResults is initially empty', async () => {
    const store = createRealStore();

    expect(store.state.testCaseResults).toEqual({});
  });

  describe('results/computeTestCaseResults', () => {
    const ERROR_RESPONSE = {
      error: 'error generation test cases, discovery model not set',
    };

    const OK_RESPONSE = { response: 'api response' };

    afterEach(() => {
      jest.resetAllMocks();
    });

    it('sets results/testCaseResults, if successful', async () => {
      const store = createRealStore();

      expect(store.state.testCaseResults).toEqual({});

      api.computeTestCaseResults.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCaseResults(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseResultsErrors', [], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_SIX, { root: true });

      expect(store.state.testCaseResults).toEqual(OK_RESPONSE);
    });

    it('sets config/errors.computeTestCaseResults, if unsuccessful', async () => {
      const store = createRealStore();

      expect(store.state.testCaseResults).toEqual({});

      api.computeTestCaseResults.mockRejectedValueOnce(ERROR_RESPONSE);
      await actions.computeTestCaseResults(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseResultsErrors', [ERROR_RESPONSE], { root: true });

      expect(store.state.testCaseResults).toEqual({});
    });

    it('sets results/testCaseResults and clears config/errors.testCaseResults, if successful', async () => {
      const store = createRealStore();

      expect(store.state.testCaseResults).toEqual({});

      api.computeTestCaseResults.mockRejectedValueOnce(ERROR_RESPONSE);
      await actions.computeTestCaseResults(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseResultsErrors', [ERROR_RESPONSE], { root: true });
      expect(store.state.testCaseResults).toEqual({});

      api.computeTestCaseResults.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCaseResults(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseResultsErrors', [], { root: true });
      expect(store.state.testCaseResults).toEqual(OK_RESPONSE);
    });

    it('clears results/testCaseResults and sets config/errors.testCaseResults, if unsuccessful', async () => {
      const store = createRealStore();

      expect(store.state.testCaseResults).toEqual({});

      api.computeTestCaseResults.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCaseResults(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseResultsErrors', [], { root: true });
      expect(store.state.testCaseResults).toEqual(OK_RESPONSE);

      api.computeTestCaseResults.mockRejectedValueOnce(ERROR_RESPONSE);
      await actions.computeTestCaseResults(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseResultsErrors', [ERROR_RESPONSE], { root: true });
      expect(store.state.testCaseResults).toEqual({});
    });
  });
});
