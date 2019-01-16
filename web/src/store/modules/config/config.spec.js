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

describe('web/src/store/modules/config', () => {
  describe('config/configuration', () => {
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

    it('configurationErrors initially empty', async () => {
      const store = createRealStore();

      expect(store.getters.configurationErrors).toEqual([]);
    });

    it('configuration.{signing_private,signing_public,transport_private,transport_public} initially empty', async () => {
      const store = createRealStore();

      expect(store.getters.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
      });
    });

    it('setConfigurationSigningPrivate', async () => {
      const store = createRealStore();

      const signingPrivate = 'signingPrivate';
      await store.dispatch('setConfigurationSigningPrivate', signingPrivate);

      expect(store.getters.configuration).toEqual({
        signing_private: signingPrivate,
        signing_public: '',
        transport_private: '',
        transport_public: '',
      });
    });

    it('setConfigurationSigningPublic', async () => {
      const store = createRealStore();

      const signingPublic = 'signingPublic';
      await store.dispatch('setConfigurationSigningPublic', signingPublic);

      expect(store.getters.configuration).toEqual({
        signing_private: '',
        signing_public: signingPublic,
        transport_private: '',
        transport_public: '',
      });
    });

    it('setConfigurationTransportPrivate', async () => {
      const store = createRealStore();

      const transportPrivate = 'transportPrivate';
      await store.dispatch('setConfigurationTransportPrivate', transportPrivate);

      expect(store.getters.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: transportPrivate,
        transport_public: '',
      });
    });

    it('setConfigurationTransportPublic', async () => {
      const store = createRealStore();

      const transportPublic = 'transportPublic';
      await store.dispatch('setConfigurationTransportPublic', transportPublic);

      expect(store.getters.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: transportPublic,
      });
    });

    describe('validateConfiguration', () => {
      afterEach(() => {
        jest.resetAllMocks();
      });

      it('setConfigurationSigningPrivate not called before validateConfiguration', async () => {
        const store = createRealStore();

        expect(store.getters.configurationErrors).toEqual([]);
        await store.dispatch('setConfigurationSigningPublic', 'setConfigurationSigningPublic');
        await store.dispatch('setConfigurationTransportPrivate', 'setConfigurationTransportPrivate');
        await store.dispatch('setConfigurationTransportPublic', 'setConfigurationTransportPublic');
        expect(store.getters.configurationErrors).toEqual([]);

        const valid = await store.dispatch('validateConfiguration');
        expect(valid).toEqual(false);

        expect(store.getters.configurationErrors).toEqual([
          'Signing Private Certificate (.key) empty',
        ]);
      });

      it('setConfigurationSigningPublic not called before validateConfiguration', async () => {
        const store = createRealStore();

        expect(store.getters.configurationErrors).toEqual([]);
        await store.dispatch('setConfigurationSigningPrivate', 'setConfigurationSigningPrivate');
        await store.dispatch('setConfigurationTransportPrivate', 'setConfigurationTransportPrivate');
        await store.dispatch('setConfigurationTransportPublic', 'setConfigurationTransportPublic');
        expect(store.getters.configurationErrors).toEqual([]);

        const valid = await store.dispatch('validateConfiguration');
        expect(valid).toEqual(false);

        expect(store.getters.configurationErrors).toEqual([
          'Signing Public Certificate (.pem) empty',
        ]);
      });

      it('setConfigurationTransportPrivate not called before validateConfiguration', async () => {
        const store = createRealStore();

        expect(store.getters.configurationErrors).toEqual([]);
        await store.dispatch('setConfigurationSigningPublic', 'setConfigurationSigningPublic');
        await store.dispatch('setConfigurationSigningPrivate', 'setConfigurationSigningPrivate');
        await store.dispatch('setConfigurationTransportPublic', 'setConfigurationTransportPublic');
        expect(store.getters.configurationErrors).toEqual([]);

        const valid = await store.dispatch('validateConfiguration');
        expect(valid).toEqual(false);

        expect(store.getters.configurationErrors).toEqual([
          'Transport Private Certificate (.key) empty',
        ]);
      });

      it('setConfigurationTransportPublic not called before validateConfiguration', async () => {
        const store = createRealStore();

        expect(store.getters.configurationErrors).toEqual([]);
        await store.dispatch('setConfigurationSigningPublic', 'setConfigurationSigningPublic');
        await store.dispatch('setConfigurationSigningPrivate', 'setConfigurationSigningPrivate');
        await store.dispatch('setConfigurationTransportPrivate', 'setConfigurationTransportPrivate');
        expect(store.getters.configurationErrors).toEqual([]);

        const valid = await store.dispatch('validateConfiguration');
        expect(valid).toEqual(false);

        expect(store.getters.configurationErrors).toEqual([
          'Transport Public Certificate (.pem) empty',
        ]);
      });

      it('setConfigurationSigningPrivate, setConfigurationSigningPublic, setConfigurationTransportPrivate and setConfigurationTransportPublic not called before validateConfiguration', async () => {
        const store = createRealStore();

        expect(store.getters.configurationErrors).toEqual([]);

        const valid = await store.dispatch('validateConfiguration');
        expect(valid).toEqual(false);

        expect(store.getters.configurationErrors).toEqual([
          'Signing Private Certificate (.key) empty',
          'Signing Public Certificate (.pem) empty',
          'Transport Private Certificate (.key) empty',
          'Transport Public Certificate (.pem) empty',
        ]);
      });

      it('setConfigurationSigningPrivate, setConfigurationSigningPublic, setConfigurationTransportPrivate and setConfigurationTransportPublic called before validateConfiguration', async () => {
        api.validateConfiguration.mockResolvedValue({
          signing_private: 'does_not_matter_what_the_value_is',
          signing_public: 'does_not_matter_what_the_value_is',
          transport_private: 'does_not_matter_what_the_value_is',
          transport_public: 'does_not_matter_what_the_value_is',
        });

        const store = createRealStore();

        expect(store.getters.configurationErrors).toEqual([]);
        await store.dispatch('setConfigurationSigningPublic', 'setConfigurationSigningPublic');
        await store.dispatch('setConfigurationSigningPrivate', 'setConfigurationSigningPrivate');
        await store.dispatch('setConfigurationTransportPrivate', 'setConfigurationTransportPrivate');
        await store.dispatch('setConfigurationTransportPublic', 'setConfigurationTransportPublic');
        expect(store.getters.configurationErrors).toEqual([]);

        const valid = await store.dispatch('validateConfiguration');
        expect(valid).toEqual(true);

        expect(store.getters.configurationErrors).toEqual([]);
      });

      it('setConfigurationSigningPrivate, setConfigurationSigningPublic, setConfigurationTransportPrivate and setConfigurationTransportPublic called with invalid values before validateConfiguration', async () => {
        const errorResponse = {
          error: "error with signing certificate: error with public key: asn1: structure error: tags don't match (16 vs {class:0 tag:2 length:1 isCompound:false}) {optional:false explicit:false application:false private:false defaultValue:\u003cnil\u003e tag:\u003cnil\u003e stringType:0 timeType:0 set:false omitEmpty:false} tbsCertificate @2",
        };
        api.validateConfiguration.mockRejectedValue(errorResponse);

        const store = createRealStore();

        expect(store.getters.configurationErrors).toEqual([]);
        await store.dispatch('setConfigurationSigningPublic', 'not_a_certificate');
        await store.dispatch('setConfigurationSigningPrivate', 'not_a_certificate');
        await store.dispatch('setConfigurationTransportPrivate', 'not_a_certificate');
        await store.dispatch('setConfigurationTransportPublic', 'not_a_certificate');
        expect(store.getters.configurationErrors).toEqual([]);

        const valid = await store.dispatch('validateConfiguration');
        expect(valid).toEqual(false);

        expect(store.getters.configurationErrors).toEqual([
          errorResponse,
        ]);
      });

      it('validateConfiguration clears previous errors', async () => {
        const store = createRealStore();

        // This will generate an error because we have not called any of the methods
        // that sets the values for the configuration.
        expect(store.getters.configurationErrors).toEqual([]);
        expect(await store.dispatch('validateConfiguration')).toEqual(false);
        expect(store.getters.configurationErrors).toEqual([
          'Signing Private Certificate (.key) empty',
          'Signing Public Certificate (.pem) empty',
          'Transport Private Certificate (.key) empty',
          'Transport Public Certificate (.pem) empty',
        ]);

        api.validateConfiguration.mockResolvedValue({
          signing_private: 'does_not_matter_what_the_value_is',
          signing_public: 'does_not_matter_what_the_value_is',
          transport_private: 'does_not_matter_what_the_value_is',
          transport_public: 'does_not_matter_what_the_value_is',
        });

        await store.dispatch('setConfigurationSigningPublic', 'setConfigurationSigningPublic');
        await store.dispatch('setConfigurationSigningPrivate', 'setConfigurationSigningPrivate');
        await store.dispatch('setConfigurationTransportPrivate', 'setConfigurationTransportPrivate');
        await store.dispatch('setConfigurationTransportPublic', 'setConfigurationTransportPublic');
        // This will clear out the previous errors, and will result in configurationErrors
        // being empty since they are not any errors.
        expect(await store.dispatch('validateConfiguration')).toEqual(true);

        expect(store.getters.configurationErrors).toEqual([]);
      });

      it('setConfigurationErrors sets errors', async () => {
        const error = new Error('e');
        const store = createRealStore();

        await store.dispatch('setConfigurationErrors', [error]);
        expect(store.getters.configurationErrors).toEqual([error]);
      });
    });
  });

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

    it('config/errors.testCaseResults is initially empty', async () => {
      const store = createRealStore();

      expect(store.getters.errors.testCaseResults).toEqual([]);
    });

    it('config/testCaseResults is initially empty', async () => {
      const store = createRealStore();

      expect(store.state.testCaseResults).toEqual({});
    });

    describe('config/computeTestCaseResults', () => {
      const ERROR_RESPONSE = {
        error: 'error generation test cases, discovery model not set',
      };

      const OK_RESPONSE = { response: 'api response' };

      afterEach(() => {
        jest.resetAllMocks();
      });

      it('config/computeTestCaseResults sets config/testCaseResults, if successful', async () => {
        const store = createRealStore();

        expect(store.getters.errors.testCaseResults).toEqual([]);
        expect(store.state.testCaseResults).toEqual({});

        api.computeTestCaseResults.mockResolvedValueOnce(OK_RESPONSE);
        await store.dispatch('computeTestCaseResults');
        expect(store.getters.errors.testCaseResults).toEqual([]);
        expect(store.state.testCaseResults).toEqual(OK_RESPONSE);
      });

      it('config/computeTestCaseResults sets config/errors.computeTestCaseResults, if unsuccessful', async () => {
        const store = createRealStore();

        expect(store.getters.errors.testCaseResults).toEqual([]);
        expect(store.state.testCaseResults).toEqual({});

        api.computeTestCaseResults.mockRejectedValueOnce(ERROR_RESPONSE);
        await store.dispatch('computeTestCaseResults');
        expect(store.getters.errors.testCaseResults).toEqual([ERROR_RESPONSE]);
        expect(store.state.testCaseResults).toEqual({});
      });

      it('config/computeTestCaseResults sets config/testCaseResults and clears config/errors.testCaseResults, if successful', async () => {
        const store = createRealStore();

        expect(store.getters.errors.testCaseResults).toEqual([]);
        expect(store.state.testCaseResults).toEqual({});

        api.computeTestCaseResults.mockRejectedValueOnce(ERROR_RESPONSE);
        await store.dispatch('computeTestCaseResults');
        expect(store.getters.errors.testCaseResults).toEqual([ERROR_RESPONSE]);
        expect(store.state.testCaseResults).toEqual({});

        api.computeTestCaseResults.mockResolvedValueOnce(OK_RESPONSE);
        await store.dispatch('computeTestCaseResults');
        expect(store.getters.errors.testCaseResults).toEqual([]);
        expect(store.state.testCaseResults).toEqual(OK_RESPONSE);
      });

      it('config/computeTestCaseResults clears config/testCaseResults and sets config/errors.testCaseResults, if unsuccessful', async () => {
        const store = createRealStore();

        expect(store.getters.errors.testCaseResults).toEqual([]);
        expect(store.state.testCaseResults).toEqual({});

        api.computeTestCaseResults.mockResolvedValueOnce(OK_RESPONSE);
        await store.dispatch('computeTestCaseResults');
        expect(store.getters.errors.testCaseResults).toEqual([]);
        expect(store.state.testCaseResults).toEqual(OK_RESPONSE);

        api.computeTestCaseResults.mockRejectedValueOnce(ERROR_RESPONSE);
        await store.dispatch('computeTestCaseResults');
        expect(store.getters.errors.testCaseResults).toEqual([ERROR_RESPONSE]);
        expect(store.state.testCaseResults).toEqual({});
      });
    });
  });
});
