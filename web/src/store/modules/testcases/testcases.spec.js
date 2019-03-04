/**
 * This creates a real store so avoid having to mock things.
 * This makes testing much easier.
 *
 * See the recommendation:
 * https://vue-test-utils.vuejs.org/guides/using-with-vuex.html#testing-a-running-store
 */
import { createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import * as _ from 'lodash';

import { WebSocket, Server } from 'mock-socket'; // https://github.com/thoov/mock-socket

import actions from './actions';
import mutations from './mutations';
import getters from './getters';
import state from './state';

import constants from '../config/constants';
import api from '../../../api';

// https://jestjs.io/docs/en/mock-functions#mocking-modules
jest.mock('../../../api');

/*
 * By default the global WebSocket object is stubbed out. However,
 * if you need to stub something else out you can like so:
 */
window.WebSocket = WebSocket; // Here we stub out the window object

describe('store/modules/testcases', () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  const CONSENT_URL = 'http://example.com';
  const CONSENT_URL2 = 'http://example.com/2';
  const SPEC_NAME = 'Account and Transaction API Specification';
  const OK_RESPONSE = {
    specTokens: [
      {
        specIdentifier: SPEC_NAME,
        namedPermissions: [
          {
            name: 'to1002',
            codeSet: {
              codes: [
                'ReadAccountsBasic',
              ],
              testIds: [
                '#co0001',
                '#t1001',
              ],
            },
            consentUrl: CONSENT_URL,
          },
          {
            name: 'to1003',
            codeSet: {
              codes: [
                'ReadAccountsDetail',
              ],
              testIds: [
                '#co0002',
                '#t1002',
              ],
            },
            consentUrl: CONSENT_URL2,
          },
        ],
      },
    ],
    specCases: [
      {
        apiSpecification: {
          name: SPEC_NAME,
          url: 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.1',
          version: 'v3.1',
          schemaVersion: 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json',
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
    ],
  };

  let dispatch;
  /**
   * Creates a real store so we don't have to mock things out.
   */
  const createRealStore = () => {
    const localVue = createLocalVue();
    localVue.use(Vuex);
    const store = new Vuex.Store({
      state: _.cloneDeep(state),
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

    expect(store.state.testCases).toEqual([]);
  });

  describe('testcases/computeTestCases', () => {
    const ERROR_RESPONSE = {
      error: 'error generation test cases, discovery model not set',
    };

    // Add additional `meta.status` field to each individual testCase in `OK_RESPONSE`.
    const EXPECTED_TESTCASES_STATE = _.map(OK_RESPONSE.specCases, (spec) => {
      const testCases = _.map(spec.testCases, testCase => _.merge({}, testCase, {
        _rowVariant: null,
        _showDetails: false,
        meta: {
          status: '',
          metrics: {
            responseSize: '',
            responseTime: '',
          },

        },
      }));
      return _.merge({}, spec, { testCases });
    });

    const EXPECTED_CONSENT_URLS_STATE = {
      [SPEC_NAME]: [CONSENT_URL, CONSENT_URL2],
    };

    it('testcases/computeTestCases sets testcases/testCases, if successful', async () => {
      expect.assertions(10);
      const fakeURL = 'ws://localhost/api/run/ws';
      const mockServer = new Server(fakeURL);
      const store = createRealStore();

      expect(store.state.testCases).toEqual([]);
      expect(store.state.ws.connection).toBeNull();

      mockServer.on('connection', (socket) => {
        expect(store.state.hasRunStarted).toEqual(false);
        expect(store.state.ws.connection).toBeDefined();
        expect(store.state.testCases).toEqual(EXPECTED_TESTCASES_STATE);
        expect(store.state.consentUrls).toEqual(EXPECTED_CONSENT_URLS_STATE);

        expect(dispatch).toHaveBeenCalledWith('status/setShowLoading', true, { root: true });
        expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });

        socket.close();
        mockServer.stop(() => {
          expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
        });
      });

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCases(store);

      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
    });

    it('testcases/computeTestCases sets config/errors.testCases, if unsuccessful', async () => {
      expect.assertions(6);

      const store = createRealStore();

      expect(store.state.testCases).toEqual([]);
      expect(store.state.ws.connection).toBeNull();

      api.computeTestCases.mockRejectedValueOnce(ERROR_RESPONSE);
      await actions.computeTestCases(store);

      expect(store.state.ws.connection).toBeDefined();
      expect(store.state.testCases).toEqual([]);

      expect(dispatch).toHaveBeenCalledWith('status/setErrors', [ERROR_RESPONSE], { root: true });
      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FOUR, { root: true });
    });
  });

  describe('testcases/executeTestCases', () => {
    // Add additional `meta.status` field to each individual testCase in `OK_RESPONSE`.
    const EXPECTED_TESTCASES_STATE_PENDING = _.map(OK_RESPONSE.specCases, (spec) => {
      const testCases = _.map(spec.testCases, testCase => _.merge({}, testCase, {
        _rowVariant: null,
        _showDetails: false,
        meta: {
          status: 'PENDING',
          metrics: {
            responseSize: '',
            responseTime: '',
          },
        },
      }));
      return _.merge({}, spec, { testCases });
    });
    const EXPECTED_TESTCASES_STATE_NOT_STARTED = _.map(OK_RESPONSE.specCases, (spec) => {
      const testCases = _.map(spec.testCases, testCase => _.merge({}, testCase, {
        _rowVariant: null,
        _showDetails: false,
        meta: {
          status: '',
          metrics: {
            responseSize: '',
            responseTime: '',
          },
        },
      }));
      return _.merge({}, spec, { testCases });
    });

    it('testcases/state.testsCases have \'PENDING\' state when testcases/executeTestCases is called', async () => {
      expect.assertions(9);

      const fakeURL = 'ws://localhost/api/run/ws';
      const mockServer = new Server(fakeURL);
      const store = createRealStore();

      expect(store.state.testCases).toEqual([]);
      expect(store.state.hasRunStarted).toEqual(false);
      expect(store.state.ws.connection).toBeNull();

      mockServer.on('connection', async (socket) => {
        expect(store.state.testCases).toEqual(EXPECTED_TESTCASES_STATE_NOT_STARTED);
        expect(store.state.ws.connection).toBeDefined();

        api.executeTestCases.mockResolvedValueOnce({});
        await actions.executeTestCases(store);

        expect(store.state.hasRunStarted).toEqual(true);
        expect(store.state.testCases).toEqual(EXPECTED_TESTCASES_STATE_PENDING);

        socket.close();
        mockServer.stop(() => {
          expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
        });
      });

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCases(store);

      expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
    });

    it('testcases/state.testsCases are left in \'NOT_STARTED\' state when testcases/executeTestCases is called', async () => {
      const fakeURL = 'ws://localhost/api/run/ws';
      const mockServer = new Server(fakeURL);
      const store = createRealStore();

      const ERROR_RESPONSE = {
        error: 'testcases/state.testsCases are left in \'NOT_STARTED\' state when testcases/executeTestCases is called',
      };

      expect(store.state.testCases).toEqual([]);
      expect(store.state.hasRunStarted).toEqual(false);
      expect(store.state.ws.connection).toBeNull();

      mockServer.on('connection', async (socket) => {
        expect(store.state.testCases).toEqual(EXPECTED_TESTCASES_STATE_NOT_STARTED);

        api.executeTestCases.mockRejectedValueOnce(ERROR_RESPONSE);
        await actions.executeTestCases(store);

        expect(store.state.testCases).toEqual(EXPECTED_TESTCASES_STATE_NOT_STARTED);
        expect(store.state.hasRunStarted).toEqual(false);

        // expect(dispatch).toHaveBeenCalledWith('status/setErrors', [ERROR_RESPONSE], { root: true });
        // expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FIVE, { root: true });

        socket.close();
        mockServer.stop(() => {
          expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
        });
      });

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCases(store);

      expect(dispatch).toHaveBeenCalledWith('config/setWizardStep', constants.WIZARD.STEP_FIVE, { root: true });
    });

    it('testcases/state.testsCases are updated when update arrives on the WebSocket', async () => {
      const fakeURL = 'ws://localhost/api/run/ws';
      const mockServer = new Server(fakeURL);
      const store = createRealStore();

      expect(store.state.testCases).toEqual([]);
      expect(store.state.hasRunStarted).toEqual(false);
      expect(store.state.ws.connection).toBeNull();

      mockServer.on('connection', async (socket) => {
        expect(store.state.testCases).toEqual(EXPECTED_TESTCASES_STATE_NOT_STARTED);

        api.executeTestCases.mockResolvedValueOnce({});
        await actions.executeTestCases(store);

        expect(store.state.hasRunStarted).toEqual(true);
        expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
        expect(store.state.testCases).toEqual(EXPECTED_TESTCASES_STATE_PENDING);

        const message = JSON.stringify({ test: { id: '#t1000', pass: true } });
        socket.send(message);

        await new Promise(resolve => setTimeout(resolve, 1000));
        expect(store.state.testCases[0].testCases[0].meta.status).toEqual('PASSED');

        socket.close();
        mockServer.stop(() => {
          expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
        });
      });

      api.computeTestCases.mockResolvedValueOnce(OK_RESPONSE);
      await actions.computeTestCases(store);
    });
  });
});
