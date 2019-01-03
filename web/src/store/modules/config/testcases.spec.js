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

import api from '../../../api';
// https://jestjs.io/docs/en/mock-functions#mocking-modules
jest.mock('../../../api');

describe('config/testCases', () => {
  /**
   * Creates a real store so we don't have to mock things out.
   */
  const createRealStore = () => {
    const localVue = createLocalVue();
    localVue.use(Vuex);

    return new Vuex.Store({
      state: cloneDeep(state),
      actions,
      mutations,
      getters,
    });
  };

  it('config/errors.testCases is initially empty', async () => {
    const store = createRealStore();

    expect(store.getters.errors.testCases).toEqual([]);
  });

  it('config/testCases is initially empty', async () => {
    const store = createRealStore();

    expect(store.getters.testCases).toEqual([]);
  });

  describe('config/computeTestCases', () => {
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

    it('config/computeTestCases sets config/testCases, if successful', async () => {
      const store = createRealStore();

      expect(store.getters.errors.testCases).toEqual([]);
      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await store.dispatch('computeTestCases');
      expect(store.getters.errors.testCases).toEqual([]);
      expect(store.getters.testCases).toEqual(OK_RESPONSE);
    });

    it('config/computeTestCases sets config/errors.testCases, if unsuccessful', async () => {
      const store = createRealStore();

      expect(store.getters.errors.testCases).toEqual([]);
      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockRejectedValueOnce(ERROR_RESPONSE);
      await store.dispatch('computeTestCases');
      expect(store.getters.errors.testCases).toEqual([ERROR_RESPONSE]);
      expect(store.getters.testCases).toEqual([]);
    });

    it('config/computeTestCases sets config/testCases and clears config/errors.testCases, if successful', async () => {
      const store = createRealStore();

      expect(store.getters.errors.testCases).toEqual([]);
      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockRejectedValueOnce(ERROR_RESPONSE);
      await store.dispatch('computeTestCases');
      expect(store.getters.errors.testCases).toEqual([ERROR_RESPONSE]);
      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await store.dispatch('computeTestCases');
      expect(store.getters.errors.testCases).toEqual([]);
      expect(store.getters.testCases).toEqual(OK_RESPONSE);
    });

    it('config/computeTestCases clears config/testCases and sets config/errors.testCases, if unsuccessful', async () => {
      const store = createRealStore();

      expect(store.getters.errors.testCases).toEqual([]);
      expect(store.getters.testCases).toEqual([]);

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await store.dispatch('computeTestCases');
      expect(store.getters.errors.testCases).toEqual([]);
      expect(store.getters.testCases).toEqual(OK_RESPONSE);

      api.computeTestCases.mockRejectedValueOnce(ERROR_RESPONSE);
      await store.dispatch('computeTestCases');
      expect(store.getters.errors.testCases).toEqual([ERROR_RESPONSE]);
      expect(store.getters.testCases).toEqual([]);
    });
  });
});
