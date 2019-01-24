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

describe('store/modules/testcases', () => {
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

  it('testcases/testCases is initially empty', async () => {
    const store = createRealStore();

    expect(store.getters.testCases).toEqual([]);
  });

  describe('testcases/computeTestCases', () => {
    const ERROR_RESPONSE = {
      error: 'error generation test cases, discovery model not set',
    };

    const OK_RESPONSE = [
      {
        apiSpecification: {
          name: 'Account and Transaction API Specification',
          url: 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0',
          version: 'v3.0',
          schemaVersion: 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json',
        },
        testCases: [
          {
            '@id': '#t1000',
            name: 'Create Account Access Consents',
            input: {
              method: 'POST',
              endpoint: '/account-access-consents',
              contextGet: {},
            },
            expect: {
              'status-code': 201,
              'schema-validation': true,
              contextPut: {},
            },
          },
        ],
      },
    ];

    afterEach(() => {
      jest.resetAllMocks();
    });

    it('testcases/computeTestCases sets testcases/testCases, if successful', async () => {
      const store = createRealStore();

      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCases(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseErrors', [], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
      expect(store.getters.testCases).toEqual(OK_RESPONSE);
    });

    it('testcases/computeTestCases sets config/errors.testCases, if unsuccessful', async () => {
      const store = createRealStore();

      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockRejectedValueOnce(ERROR_RESPONSE);
      await actions.computeTestCases(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseErrors', [ERROR_RESPONSE], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
      expect(store.getters.testCases).toEqual([]);
    });

    it('testcases/computeTestCases sets testcases/testCases and clears config/errors.testCases, if successful', async () => {
      const store = createRealStore();

      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockRejectedValueOnce(ERROR_RESPONSE);
      await actions.computeTestCases(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseErrors', [ERROR_RESPONSE], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCases(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseErrors', [], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
      expect(store.getters.testCases).toEqual(OK_RESPONSE);
    });

    it('testcases/computeTestCases clears testcases/testCases and sets config/errors.testCases, if unsuccessful', async () => {
      const store = createRealStore();

      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCases(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseErrors', [], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
      expect(store.getters.testCases).toEqual(OK_RESPONSE);

      api.computeTestCases.mockRejectedValueOnce(ERROR_RESPONSE);
      await actions.computeTestCases(store);
      expect(dispatch).toHaveBeenCalledWith('config/setTestCaseErrors', [ERROR_RESPONSE], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
      expect(store.getters.testCases).toEqual([]);
    });
  });
});
