/**
 * This creates a real store so avoid having to mock things.
 * This makes testing much easier.
 *
 * See the recommendation:
 * https://vue-test-utils.vuejs.org/guides/using-with-vuex.html#testing-a-running-store
 */
import { createLocalVue } from '@vue/test-utils';
import { cloneDeep } from 'lodash';
import Vuex from 'vuex';
import api from '../../../api';
import actions from './actions';
import {
  getters, mutations,
  // import state - please don't remove comment.
  state,
} from './index';
import * as types from './mutation-types.js';

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

    it('configuration.{signing_private,signing_public,transport_private,transport_public, client_id, client_secret, token_endpoint, authorization_endpoint, resource_base_url, x_fapi_financial_id} initially empty', async () => {
      const store = createRealStore();

      expect(store.getters.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '',
        client_secret: '',
        token_endpoint: '',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',


        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });
    });

    it('setConfigurationSigningPrivate', async () => {
      const store = createRealStore();

      const signingPrivate = 'signingPrivate';
      await actions.setConfigurationSigningPrivate(store, signingPrivate);

      expect(store.getters.configuration.signing_private).toEqual(signingPrivate);
    });

    it('setConfigurationSigningPublic', async () => {
      const store = createRealStore();

      const signingPublic = 'signingPublic';
      await actions.setConfigurationSigningPublic(store, signingPublic);

      expect(store.getters.configuration.signing_public).toEqual(signingPublic);
    });

    it('setConfigurationTransportPrivate', async () => {
      const store = createRealStore();

      const transportPrivate = 'transportPrivate';
      await actions.setConfigurationTransportPrivate(store, transportPrivate);

      expect(store.getters.configuration.transport_private).toEqual(transportPrivate);
    });

    it('setConfigurationTransportPublic', async () => {
      const store = createRealStore();

      const transportPublic = 'transportPublic';
      await actions.setConfigurationTransportPublic(store, transportPublic);

      expect(store.getters.configuration.transport_public).toEqual(transportPublic);
    });

    it('commits client_id, client_secret, token_endpoint, request_object_signing_alg, authorization_endpoint, resource_base_url, x_fapi_financial_id, issuer, redirect_url and resource_ids', async () => {
      const store = createRealStore();

      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '',
        client_secret: '',
        token_endpoint: '',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_TOKEN_ENDPOINT_AUTH_METHODS, ['tls_client_auth', 'client_secret_basic']);
      store.commit(types.SET_CLIENT_ID, '8672384e-9a33-439f-8924-67bb14340d71');
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '',
        token_endpoint: '',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_CLIENT_SECRET, '2cfb31a3-5443-4e65-b2bc-ef8e00266a77');
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '2cfb31a3-5443-4e65-b2bc-ef8e00266a77',
        token_endpoint: '',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_TOKEN_ENDPOINT, 'https://modelobank2018.o3bank.co.uk:4201/token');
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '2cfb31a3-5443-4e65-b2bc-ef8e00266a77',
        token_endpoint: 'https://modelobank2018.o3bank.co.uk:4201/token',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_TOKEN_ENDPOINT_AUTH_METHOD, 'client_secret_basic');
      expect(store.state.configuration.token_endpoint_auth_method).toEqual('client_secret_basic');

      store.commit(types.SET_AUTHORIZATION_ENDPOINT, 'https://modelobankauth2018.o3bank.co.uk:4101/auth');
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '2cfb31a3-5443-4e65-b2bc-ef8e00266a77',
        token_endpoint: 'https://modelobank2018.o3bank.co.uk:4201/token',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_RESOURCE_BASE_URL, 'https://ob19-rs1.o3bank.co.uk:4501');
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '2cfb31a3-5443-4e65-b2bc-ef8e00266a77',
        token_endpoint: 'https://modelobank2018.o3bank.co.uk:4201/token',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
        resource_base_url: 'https://ob19-rs1.o3bank.co.uk:4501',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_X_FAPI_FINANCIAL_ID, '0015800001041RHAAY');
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '2cfb31a3-5443-4e65-b2bc-ef8e00266a77',
        token_endpoint: 'https://modelobank2018.o3bank.co.uk:4201/token',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
        resource_base_url: 'https://ob19-rs1.o3bank.co.uk:4501',
        x_fapi_financial_id: '0015800001041RHAAY',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_ISSUER, 'https://modelobankauth2018.o3bank.co.uk:4101');
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '2cfb31a3-5443-4e65-b2bc-ef8e00266a77',
        token_endpoint: 'https://modelobank2018.o3bank.co.uk:4201/token',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
        resource_base_url: 'https://ob19-rs1.o3bank.co.uk:4501',
        x_fapi_financial_id: '0015800001041RHAAY',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: 'https://modelobankauth2018.o3bank.co.uk:4101',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.ADD_RESOURCE_ACCOUNT_ID, { account_id: 'account-id' });
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '2cfb31a3-5443-4e65-b2bc-ef8e00266a77',
        token_endpoint: 'https://modelobank2018.o3bank.co.uk:4201/token',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
        resource_base_url: 'https://ob19-rs1.o3bank.co.uk:4501',
        x_fapi_financial_id: '0015800001041RHAAY',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: 'https://modelobankauth2018.o3bank.co.uk:4101',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }, { account_id: 'account-id' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.ADD_RESOURCE_STATEMENT_ID, { statement_id: 'statement-id' });
      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '8672384e-9a33-439f-8924-67bb14340d71',
        client_secret: '2cfb31a3-5443-4e65-b2bc-ef8e00266a77',
        token_endpoint: 'https://modelobank2018.o3bank.co.uk:4201/token',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: 'https://modelobankauth2018.o3bank.co.uk:4101/auth',
        resource_base_url: 'https://ob19-rs1.o3bank.co.uk:4501',
        x_fapi_financial_id: '0015800001041RHAAY',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: 'https://modelobankauth2018.o3bank.co.uk:4101',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }, { account_id: 'account-id' }],
          statement_ids: [{ statement_id: '' }, { statement_id: 'statement-id' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_REQUEST_OBJECT_SIGNING_ALG, 'PS256');
      expect(store.state.configuration.request_object_signing_alg).toEqual('PS256');
    });

    it('sets resource_ids', async () => {
      const store = createRealStore();

      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '',
        client_secret: '',
        token_endpoint: '',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      const acctIDs = [{ account_id: '123' }, { account_id: '456' }];
      store.commit(types.SET_RESOURCE_ACCOUNT_IDS, acctIDs);

      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '',
        client_secret: '',
        token_endpoint: '',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '123' }, { account_id: '456' }],
          statement_ids: [{ statement_id: '' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      const stmtIDs = [{ statement_id: '123' }, { statement_id: '456' }];
      store.commit(types.SET_RESOURCE_STATEMENT_IDS, stmtIDs);

      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',
        transaction_from_date: '',
        transaction_to_date: '',
        client_id: '',
        client_secret: '',
        token_endpoint: '',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '123' }, { account_id: '456' }],
          statement_ids: [{ statement_id: '123' }, { statement_id: '456' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });

      store.commit(types.SET_TRANSACTION_FROM_DATE, '2016-01-01T10:40:00+02:00');
      store.commit(types.SET_TRANSACTION_TO_DATE, '2016-01-01T10:40:00+02:00');

      expect(store.state.configuration).toEqual({
        signing_private: '',
        signing_public: '',
        transport_private: '',
        transport_public: '',
        tpp_signature_kid: '',
        tpp_signature_issuer: '',
        tpp_signature_tan: 'openbanking.org.uk',

        transaction_from_date: '2016-01-01T10:40:00+02:00',
        transaction_to_date: '2016-01-01T10:40:00+02:00',
        client_id: '',
        client_secret: '',
        token_endpoint: '',
        response_type: '',
        token_endpoint_auth_method: 'client_secret_basic',
        request_object_signing_alg: '',
        authorization_endpoint: '',
        resource_base_url: '',
        x_fapi_financial_id: '',
        send_x_fapi_customer_ip_address: false,
        x_fapi_customer_ip_address: '',
        issuer: '',
        redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
        resource_ids: {
          account_ids: [{ account_id: '123' }, { account_id: '456' }],
          statement_ids: [{ statement_id: '123' }, { statement_id: '456' }],
        },
        creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        international_creditor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        cbpii_debtor_account: {
          scheme_name: '',
          identification: '',
          name: '',
        },
        currency_of_transfer: 'USD',
        instructed_amount: {
          currency: 'GBP',
          value: '1.00',
        },
        payment_frequency: 'EvryDay',
        first_payment_date_time: '2022-01-01T00:00:00+01:00',
        requested_execution_date_time: '2022-01-01T00:00:00+01:00',
        acr_values_supported: [],
        conditional_properties: [],
      });
    });

    describe('validateDiscoveryConfig', () => {
      afterEach(() => {
        jest.resetAllMocks();
      });

      it('commits authorization_endpoint, token_endpoint, issuer  after success', async () => {
        const store = createRealStore();

        expect(store.state.configuration).toEqual({
          signing_private: '',
          signing_public: '',
          transport_private: '',
          transport_public: '',
          tpp_signature_kid: '',
          tpp_signature_issuer: '',
          tpp_signature_tan: 'openbanking.org.uk',
          transaction_from_date: '',
          transaction_to_date: '',
          client_id: '',
          client_secret: '',
          token_endpoint: '',
          response_type: '',
          token_endpoint_auth_method: 'client_secret_basic',
          request_object_signing_alg: '',
          authorization_endpoint: '',
          resource_base_url: '',
          x_fapi_financial_id: '',
          send_x_fapi_customer_ip_address: false,
          x_fapi_customer_ip_address: '',
          issuer: '',
          redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
          resource_ids: {
            account_ids: [{ account_id: '' }],
            statement_ids: [{ statement_id: '' }],
          },
          creditor_account: {
            scheme_name: '',
            identification: '',
            name: '',
          },
          international_creditor_account: {
            scheme_name: '',
            identification: '',
            name: '',
          },
          cbpii_debtor_account: {
            scheme_name: '',
            identification: '',
            name: '',
          },
          currency_of_transfer: 'USD',
          instructed_amount: {
            currency: 'GBP',
            value: '1.00',
          },
          payment_frequency: 'EvryDay',
          first_payment_date_time: '2022-01-01T00:00:00+01:00',
          requested_execution_date_time: '2022-01-01T00:00:00+01:00',
          acr_values_supported: [],
          conditional_properties: [],
        });

        api.validateDiscoveryConfig.mockReturnValueOnce({
          success: true,
          problems: [],
          response: {
            token_endpoints: {
              'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json': '',
            },
            default_token_endpoint_auth_method: {
              'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json':
              'client_secret_basic',
            },
            token_endpoint_auth_methods: {
              'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json':
              ['tls_client_auth', 'client_secret_basic'],
            },
            authorization_endpoints: {
              'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json': 'https://modelobankauth2018.o3bank.co.uk:4101/auth_1',
              'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/payment-initiation-swagger.json': 'https://modelobankauth2018.o3bank.co.uk:4101/auth_2',
            },
            issuers: {
              'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json': 'https://modelobankauth2018.o3bank.co.uk:4101_1',
              'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json': 'https://modelobankauth2018.o3bank.co.uk:4101_2',
            },
            request_object_signing_alg_values_supported: {
              'schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json':
                ['PS256', 'RS256'],
            },
            default_transaction_from_date: '2016-01-01T10:40:00+02:00',
            default_transaction_to_date: '2025-12-31T10:40:00+02:00',
            response_types_supported: [
              'code',
              'code id_token',
            ],
            acr_values_supported: [],
            conditional_properties: [],
          },
        });

        await actions.validateDiscoveryConfig(store);

        expect(store.state.token_endpoint_auth_methods).toEqual(['tls_client_auth', 'client_secret_basic']);
        expect(store.state.request_object_signing_alg_values_supported).toEqual(['PS256', 'RS256']);
        expect(store.state.response_types_supported).toEqual(['code', 'code id_token']);
        expect(store.state.configuration).toEqual({
          signing_private: '',
          signing_public: '',
          transport_private: '',
          transport_public: '',
          tpp_signature_kid: '',
          tpp_signature_issuer: '',
          tpp_signature_tan: 'openbanking.org.uk',
          transaction_from_date: '2016-01-01T10:40:00+02:00',
          transaction_to_date: '2025-12-31T10:40:00+02:00',
          client_id: '',
          client_secret: '',
          token_endpoint: '',
          response_type: '',
          token_endpoint_auth_method: 'client_secret_basic',
          request_object_signing_alg: '',
          authorization_endpoint: 'https://modelobankauth2018.o3bank.co.uk:4101/auth_1',
          resource_base_url: '',
          x_fapi_financial_id: '',
          send_x_fapi_customer_ip_address: false,
          x_fapi_customer_ip_address: '',
          issuer: 'https://modelobankauth2018.o3bank.co.uk:4101_1',
          redirect_url: 'https://127.0.0.1:8443/conformancesuite/callback',
          resource_ids: {
            account_ids: [{ account_id: '' }],
            statement_ids: [{ statement_id: '' }],
          },
          creditor_account: {
            scheme_name: '',
            identification: '',
            name: '',
          },
          international_creditor_account: {
            scheme_name: '',
            identification: '',
            name: '',
          },
          cbpii_debtor_account: {
            scheme_name: '',
            identification: '',
            name: '',
          },
          currency_of_transfer: 'USD',
          instructed_amount: {
            currency: 'GBP',
            value: '1.00',
          },
          payment_frequency: 'EvryDay',
          first_payment_date_time: '2022-01-01T00:00:00+01:00',
          requested_execution_date_time: '2022-01-01T00:00:00+01:00',
          acr_values_supported: [],
          conditional_properties: [],
        });
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
          'Account IDs empty',
          'Statement IDs empty',
          'Transaction From Date empty',
          'Transaction To Date empty',
          'Client ID empty',
          'Client Secret empty',
          'Token Endpoint empty',
          'response_type empty',
          'Request object signing algorithm empty',
          'Authorization Endpoint empty',
          'Resource Base URL empty',
          'x-fapi-financial-id empty',
          'issuer empty',
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
          'Account IDs empty',
          'Statement IDs empty',
          'Transaction From Date empty',
          'Transaction To Date empty',
          'Client ID empty',
          'Client Secret empty',
          'Token Endpoint empty',
          'response_type empty',
          'Request object signing algorithm empty',
          'Authorization Endpoint empty',
          'Resource Base URL empty',
          'x-fapi-financial-id empty',
          'issuer empty',
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
          'Account IDs empty',
          'Statement IDs empty',
          'Transaction From Date empty',
          'Transaction To Date empty',
          'Client ID empty',
          'Client Secret empty',
          'Token Endpoint empty',
          'response_type empty',
          'Request object signing algorithm empty',
          'Authorization Endpoint empty',
          'Resource Base URL empty',
          'x-fapi-financial-id empty',
          'issuer empty',
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
          'Account IDs empty',
          'Statement IDs empty',
          'Transaction From Date empty',
          'Transaction To Date empty',
          'Client ID empty',
          'Client Secret empty',
          'Token Endpoint empty',
          'response_type empty',
          'Request object signing algorithm empty',
          'Authorization Endpoint empty',
          'Resource Base URL empty',
          'x-fapi-financial-id empty',
          'issuer empty',
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
          'Account IDs empty',
          'Statement IDs empty',
          'Transaction From Date empty',
          'Transaction To Date empty',
          'Client ID empty',
          'Client Secret empty',
          'Token Endpoint empty',
          'response_type empty',
          'Request object signing algorithm empty',
          'Authorization Endpoint empty',
          'Resource Base URL empty',
          'x-fapi-financial-id empty',
          'issuer empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors, { root: true });
      });

      it('setConfigurationSigningPrivate, setConfigurationSigningPublic, setConfigurationTransportPrivate and setConfigurationTransportPublic called before validateConfiguration', async () => {
        api.validateConfiguration.mockReturnValueOnce({
          signing_private: 'does_not_matter_what_the_value_is',
          signing_public: 'does_not_matter_what_the_value_is',
          transport_private: 'does_not_matter_what_the_value_is',
          transport_public: 'does_not_matter_what_the_value_is',
        });

        const store = createRealStore();
        store.commit(types.SET_TRANSACTION_FROM_DATE, '2016-01-01T10:40:00+02:00');
        store.commit(types.SET_TRANSACTION_TO_DATE, '2016-01-01T10:40:00+02:00');
        store.commit(types.SET_CLIENT_ID, '8672384e-9a33-439f-8924-67bb14340d71');
        store.commit(types.SET_CLIENT_SECRET, '2cfb31a3-5443-4e65-b2bc-ef8e00266a77');
        store.commit(types.SET_TOKEN_ENDPOINT, 'https://modelobank2018.o3bank.co.uk:4201/token');
        store.commit(types.SET_RESPONSE_TYPE, 'code id_token');
        store.commit(types.SET_AUTHORIZATION_ENDPOINT, 'https://modelobankauth2018.o3bank.co.uk:4101/auth');
        store.commit(types.SET_RESOURCE_BASE_URL, 'https://ob19-rs1.o3bank.co.uk:4501');
        store.commit(types.SET_X_FAPI_FINANCIAL_ID, '0015800001041RHAAY');
        store.commit(types.SET_ISSUER, 'https://modelobankauth2018.o3bank.co.uk:4101');
        store.commit(types.SET_REQUEST_OBJECT_SIGNING_ALG, 'PS256');

        await actions.setConfigurationSigningPublic(store, 'setConfigurationSigningPublic');
        await actions.setConfigurationSigningPrivate(store, 'setConfigurationSigningPrivate');
        await actions.setConfigurationTransportPrivate(store, 'setConfigurationTransportPrivate');
        await actions.setConfigurationTransportPublic(store, 'setConfigurationTransportPublic');

        await actions.removeResourceAccountID(store, 0);
        await actions.removeResourceStatementID(store, 0);
        await store.commit(types.ADD_RESOURCE_ACCOUNT_ID, { account_id: 'account-id' });
        await store.commit(types.ADD_RESOURCE_STATEMENT_ID, { statement_id: 'statement-id' });
        await store.commit(types.SET_TRANSACTION_FROM_DATE, '2016-01-01T10:40:00+02:00');
        await store.commit(types.SET_TRANSACTION_TO_DATE, '2025-12-31T10:40:00+02:00');

        const valid = await actions.validateConfiguration(store);
        expect(valid).toEqual(true);
      });

      it('setConfigurationSigningPrivate, setConfigurationSigningPublic, setConfigurationTransportPrivate and setConfigurationTransportPublic called with invalid values before validateConfiguration', async () => {
        const errorResponse = {
          error: "error with signing certificate: error with public key: asn1: structure error: tags don't match (16 vs {class:0 tag:2 length:1 isCompound:false}) {optional:false explicit:false application:false private:false defaultValue:\u003cnil\u003e tag:\u003cnil\u003e stringType:0 timeType:0 set:false omitEmpty:false} tbsCertificate @2",
        };
        api.validateConfiguration.mockRejectedValueOnce(errorResponse);

        const store = createRealStore();
        store.commit(types.SET_TRANSACTION_FROM_DATE, '2016-01-01T10:40:00+02:00');
        store.commit(types.SET_TRANSACTION_TO_DATE, '2016-01-01T10:40:00+02:00');
        store.commit(types.SET_CLIENT_ID, '8672384e-9a33-439f-8924-67bb14340d71');
        store.commit(types.SET_CLIENT_SECRET, '2cfb31a3-5443-4e65-b2bc-ef8e00266a77');
        store.commit(types.SET_TOKEN_ENDPOINT, 'https://modelobank2018.o3bank.co.uk:4201/token');
        store.commit(types.SET_RESPONSE_TYPE, 'code id_token');
        store.commit(types.SET_AUTHORIZATION_ENDPOINT, 'https://modelobankauth2018.o3bank.co.uk:4101/auth');
        store.commit(types.SET_RESOURCE_BASE_URL, 'https://ob19-rs1.o3bank.co.uk:4501');
        store.commit(types.SET_X_FAPI_FINANCIAL_ID, '0015800001041RHAAY');
        store.commit(types.SET_ISSUER, 'https://modelobankauth2018.o3bank.co.uk:4101');
        store.commit(types.SET_REQUEST_OBJECT_SIGNING_ALG, 'PS256');

        await actions.setConfigurationSigningPublic(store, 'not_a_certificate');
        await actions.setConfigurationSigningPrivate(store, 'not_a_certificate');
        await actions.setConfigurationTransportPrivate(store, 'not_a_certificate');
        await actions.setConfigurationTransportPublic(store, 'not_a_certificate');

        await actions.removeResourceAccountID(store, 0);
        await actions.removeResourceStatementID(store, 0);
        await store.commit(types.ADD_RESOURCE_ACCOUNT_ID, { account_id: 'account-id' });
        await store.commit(types.ADD_RESOURCE_STATEMENT_ID, { statement_id: 'statement-id' });
        await store.commit(types.SET_TRANSACTION_FROM_DATE, '2016-01-01T10:40:00+02:00');
        await store.commit(types.SET_TRANSACTION_TO_DATE, '2025-12-31T10:40:00+02:00');

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
          'Account IDs empty',
          'Statement IDs empty',
          'Transaction From Date empty',
          'Transaction To Date empty',
          'Client ID empty',
          'Client Secret empty',
          'Token Endpoint empty',
          'response_type empty',
          'Request object signing algorithm empty',
          'Authorization Endpoint empty',
          'Resource Base URL empty',
          'x-fapi-financial-id empty',
          'issuer empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors, { root: true });

        api.validateConfiguration.mockReturnValueOnce({
          signing_private: 'does_not_matter_what_the_value_is',
          signing_public: 'does_not_matter_what_the_value_is',
          transport_private: 'does_not_matter_what_the_value_is',
          transport_public: 'does_not_matter_what_the_value_is',
        });

        await actions.setConfigurationSigningPublic(store, 'setConfigurationSigningPublic');
        await actions.setConfigurationSigningPrivate(store, 'setConfigurationSigningPrivate');
        await actions.setConfigurationTransportPrivate(store, 'setConfigurationTransportPrivate');
        await actions.setConfigurationTransportPublic(store, 'setConfigurationTransportPublic');

        await actions.removeResourceAccountID(store, 0);
        await actions.removeResourceStatementID(store, 0);
        await store.commit(types.ADD_RESOURCE_ACCOUNT_ID, { account_id: 'account-id' });
        await store.commit(types.ADD_RESOURCE_STATEMENT_ID, { statement_id: 'statement-id' });

        store.commit(types.SET_TRANSACTION_FROM_DATE, '2016-01-01T10:40:00+02:00');
        store.commit(types.SET_TRANSACTION_TO_DATE, '2016-01-01T10:40:00+02:00');
        store.commit(types.SET_CLIENT_ID, '8672384e-9a33-439f-8924-67bb14340d71');
        store.commit(types.SET_CLIENT_SECRET, '2cfb31a3-5443-4e65-b2bc-ef8e00266a77');
        store.commit(types.SET_TOKEN_ENDPOINT, 'https://modelobank2018.o3bank.co.uk:4201/token');
        store.commit(types.SET_RESPONSE_TYPE, 'code id_token');
        store.commit(types.SET_AUTHORIZATION_ENDPOINT, 'https://modelobankauth2018.o3bank.co.uk:4101/auth');
        store.commit(types.SET_RESOURCE_BASE_URL, 'https://ob19-rs1.o3bank.co.uk:4501');
        store.commit(types.SET_X_FAPI_FINANCIAL_ID, '0015800001041RHAAY');
        store.commit(types.SET_ISSUER, 'https://modelobankauth2018.o3bank.co.uk:4101');
        store.commit(types.SET_REQUEST_OBJECT_SIGNING_ALG, 'PS256');
        await store.commit(types.SET_TRANSACTION_FROM_DATE, '2016-01-01T10:40:00+02:00');
        await store.commit(types.SET_TRANSACTION_TO_DATE, '2025-12-31T10:40:00+02:00');

        // This will clear out the previous errors, and will result in configurationErrors
        // being empty since they are not any errors.
        expect(await actions.validateConfiguration(store)).toEqual(true);
        expect(dispatch).toHaveBeenCalledWith('status/clearErrors', null, { root: true });
      });

      it('validateConfiguration returns invalid transaction from/start date errors', async () => {
        const store = createRealStore();

        await actions.setConfigurationSigningPublic(store, 'setConfigurationSigningPublic');
        await actions.setConfigurationSigningPrivate(store, 'setConfigurationSigningPrivate');
        await actions.setConfigurationTransportPrivate(store, 'setConfigurationTransportPrivate');
        await actions.setConfigurationTransportPublic(store, 'setConfigurationTransportPublic');

        await actions.removeResourceAccountID(store, 0);
        await actions.removeResourceStatementID(store, 0);
        await store.commit(types.ADD_RESOURCE_ACCOUNT_ID, { account_id: 'account-id' });
        await store.commit(types.ADD_RESOURCE_STATEMENT_ID, { statement_id: 'statement-id' });

        store.commit(types.SET_TRANSACTION_FROM_DATE, 'xxx-invalid-date-xxx');
        store.commit(types.SET_TRANSACTION_TO_DATE, '');
        store.commit(types.SET_CLIENT_ID, '8672384e-9a33-439f-8924-67bb14340d71');
        store.commit(types.SET_CLIENT_SECRET, '2cfb31a3-5443-4e65-b2bc-ef8e00266a77');
        store.commit(types.SET_TOKEN_ENDPOINT, 'https://modelobank2018.o3bank.co.uk:4201/token');
        store.commit(types.SET_AUTHORIZATION_ENDPOINT, 'https://modelobankauth2018.o3bank.co.uk:4101/auth');
        store.commit(types.SET_RESOURCE_BASE_URL, 'https://ob19-rs1.o3bank.co.uk:4501');
        store.commit(types.SET_X_FAPI_FINANCIAL_ID, '0015800001041RHAAY');
        store.commit(types.SET_ISSUER, 'https://modelobankauth2018.o3bank.co.uk:4101');
        store.commit(types.SET_REQUEST_OBJECT_SIGNING_ALG, 'PS256');

        // This will clear out the previous errors, and will result in configurationErrors
        // being empty since they are not any errors.
        expect(await actions.validateConfiguration(store)).toEqual(false);
        const errors1 = [
          'Transaction From Date not ISO 8601 format',
          'Transaction To Date empty',
          'response_type empty',
        ];
        expect(dispatch).toHaveBeenCalledWith('status/setErrors', errors1, { root: true });
      });
    });
  });
});
