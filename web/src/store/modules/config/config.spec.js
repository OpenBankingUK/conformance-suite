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
      await actions.setConfigurationSigningPrivate(store, signingPrivate);

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
      await actions.setConfigurationSigningPublic(store, signingPublic);

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
      await actions.setConfigurationTransportPrivate(store, transportPrivate);

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
      await actions.setConfigurationTransportPublic(store, transportPublic);

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

        await actions.setConfigurationSigningPublic(store, 'setConfigurationSigningPublic');
        await actions.setConfigurationTransportPrivate(store, 'setConfigurationTransportPrivate');
        await actions.setConfigurationTransportPublic(store, 'setConfigurationTransportPublic');

        const valid = await actions.validateConfiguration(store);
        expect(valid).toEqual(false);

        const errors = [
          'Signing Private Certificate (.key) empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors, { root: true });
      });

      it('setConfigurationSigningPublic not called before validateConfiguration', async () => {
        const store = createRealStore();

        await actions.setConfigurationSigningPrivate(store, 'setConfigurationSigningPrivate');
        await actions.setConfigurationTransportPrivate(store, 'setConfigurationTransportPrivate');
        await actions.setConfigurationTransportPublic(store, 'setConfigurationTransportPublic');

        const valid = await actions.validateConfiguration(store);
        expect(valid).toEqual(false);

        const errors = [
          'Signing Public Certificate (.pem) empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors, { root: true });
      });

      it('setConfigurationTransportPrivate not called before validateConfiguration', async () => {
        const store = createRealStore();

        await actions.setConfigurationSigningPublic(store, 'setConfigurationSigningPublic');
        await actions.setConfigurationSigningPrivate(store, 'setConfigurationSigningPrivate');
        await actions.setConfigurationTransportPublic(store, 'setConfigurationTransportPublic');

        const valid = await actions.validateConfiguration(store);
        expect(valid).toEqual(false);

        const errors = [
          'Transport Private Certificate (.key) empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors, { root: true });
      });

      it('setConfigurationTransportPublic not called before validateConfiguration', async () => {
        const store = createRealStore();

        await actions.setConfigurationSigningPublic(store, 'setConfigurationSigningPublic');
        await actions.setConfigurationSigningPrivate(store, 'setConfigurationSigningPrivate');
        await actions.setConfigurationTransportPrivate(store, 'setConfigurationTransportPrivate');

        const valid = await actions.validateConfiguration(store);
        expect(valid).toEqual(false);

        const errors = [
          'Transport Public Certificate (.pem) empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors, { root: true });
      });

      it('setConfigurationSigningPrivate, setConfigurationSigningPublic, setConfigurationTransportPrivate and setConfigurationTransportPublic not called before validateConfiguration', async () => {
        const store = createRealStore();

        const valid = await actions.validateConfiguration(store);
        expect(valid).toEqual(false);

        const errors = [
          'Signing Private Certificate (.key) empty',
          'Signing Public Certificate (.pem) empty',
          'Transport Private Certificate (.key) empty',
          'Transport Public Certificate (.pem) empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors, { root: true });
      });

      it('setConfigurationSigningPrivate, setConfigurationSigningPublic, setConfigurationTransportPrivate and setConfigurationTransportPublic called before validateConfiguration', async () => {
        api.validateConfiguration.mockResolvedValue({
          signing_private: 'does_not_matter_what_the_value_is',
          signing_public: 'does_not_matter_what_the_value_is',
          transport_private: 'does_not_matter_what_the_value_is',
          transport_public: 'does_not_matter_what_the_value_is',
        });

        const store = createRealStore();

        await actions.setConfigurationSigningPublic(store, 'setConfigurationSigningPublic');
        await actions.setConfigurationSigningPrivate(store, 'setConfigurationSigningPrivate');
        await actions.setConfigurationTransportPrivate(store, 'setConfigurationTransportPrivate');
        await actions.setConfigurationTransportPublic(store, 'setConfigurationTransportPublic');

        const valid = await actions.validateConfiguration(store);
        expect(valid).toEqual(true);
      });

      it('setConfigurationSigningPrivate, setConfigurationSigningPublic, setConfigurationTransportPrivate and setConfigurationTransportPublic called with invalid values before validateConfiguration', async () => {
        const errorResponse = {
          error: "error with signing certificate: error with public key: asn1: structure error: tags don't match (16 vs {class:0 tag:2 length:1 isCompound:false}) {optional:false explicit:false application:false private:false defaultValue:\u003cnil\u003e tag:\u003cnil\u003e stringType:0 timeType:0 set:false omitEmpty:false} tbsCertificate @2",
        };
        api.validateConfiguration.mockRejectedValue(errorResponse);

        const store = createRealStore();

        await actions.setConfigurationSigningPublic(store, 'not_a_certificate');
        await actions.setConfigurationSigningPrivate(store, 'not_a_certificate');
        await actions.setConfigurationTransportPrivate(store, 'not_a_certificate');
        await actions.setConfigurationTransportPublic(store, 'not_a_certificate');

        const valid = await actions.validateConfiguration(store);
        expect(valid).toEqual(false);

        expect(dispatch).toHaveBeenCalledWith('status/setErrors', [errorResponse], { root: true });
      });

      it('validateConfiguration clears previous errors', async () => {
        const store = createRealStore();

        // This will generate an error because we have not called any of the methods
        // that sets the values for the configuration.

        expect(await actions.validateConfiguration(store)).toEqual(false);
        const errors = [
          'Signing Private Certificate (.key) empty',
          'Signing Public Certificate (.pem) empty',
          'Transport Private Certificate (.key) empty',
          'Transport Public Certificate (.pem) empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors, { root: true });

        api.validateConfiguration.mockResolvedValue({
          signing_private: 'does_not_matter_what_the_value_is',
          signing_public: 'does_not_matter_what_the_value_is',
          transport_private: 'does_not_matter_what_the_value_is',
          transport_public: 'does_not_matter_what_the_value_is',
        });

        await actions.setConfigurationSigningPublic(store, 'setConfigurationSigningPublic');
        await actions.setConfigurationSigningPrivate(store, 'setConfigurationSigningPrivate');
        await actions.setConfigurationTransportPrivate(store, 'setConfigurationTransportPrivate');
        await actions.setConfigurationTransportPublic(store, 'setConfigurationTransportPublic');
        // This will clear out the previous errors, and will result in configurationErrors
        // being empty since they are not any errors.
        expect(await actions.validateConfiguration(store)).toEqual(true);
        expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
      });
    });
  });
});
