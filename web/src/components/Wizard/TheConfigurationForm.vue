<template>
  <div class="p-3">
    <b-form>
      <b-card bg-variant="light">
        <b-form-group
          label="Client"
          label-size="lg"/>
        <ConfigurationFormFile
          id="signing_private"
          setter-method-name-suffix="SigningPrivate"
          label="Private Signing Certificate (.key):"
          validExtension=".key"
        />
        <ConfigurationFormFile
          id="signing_public"
          setter-method-name-suffix="SigningPublic"
          label="Public Signing Certificate (.pem):"
          validExtension=".pem"
        />
        <ConfigurationFormFile
          id="transport_private"
          setter-method-name-suffix="TransportPrivate"
          label="Private Transport Certificate (.key):"
          validExtension=".key"
        />
        <ConfigurationFormFile
          id="transport_public"
          setter-method-name-suffix="TransportPublic"
          label="Public Transport Certificate (.pem):"
          validExtension=".pem"
        />

        <b-form-group
          id="resource_account_id_group"
          label-for="resource_account_ids"
          label="Account IDs">
          <b-input-group
            v-for="(item, index) in resourceAccountIds"
            :key="index"
            class="mt-3">
            <b-input-group-prepend
              v-if="resourceAccountIds.length > 1">
              <b-button
                variant="danger"
                @click="removeResourceAccountIDField(index)">-</b-button>
            </b-input-group-prepend>
            <b-form-input
              :id="`resource_account_ids-${index}`"
              :value="item.account_id"
              :state="isNotEmpty(item.account_id)"
              label="Resource - Account IDs"
              placeholder="At least one Account ID is required"
              required
              type="text"
              @input="(value) => { updateAccountId(index, value) }"
            />
            <b-input-group-append
              v-if="index == resourceAccountIds.length -1">
              <b-button
                variant="success"
                @click="addResourceAccountIDField('')">+</b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>

        <b-form-group
          id="resource_statement_id_group"
          label-for="resource_statement_ids"
          label="Statement IDs">
          <b-input-group
            v-for="(item, index) in resourceStatementIds"
            :key="index"
            class="mt-3">
            <b-input-group-prepend
              v-if="resourceStatementIds.length > 1">
              <b-button
                variant="danger"
                @click="removeResourceStatementIDField(index)">-</b-button>
            </b-input-group-prepend>
            <b-form-input
              :id="`resource_statement_ids-${index}`"
              :value="item.statement_id"
              :state="isNotEmpty(item.statement_id)"
              label="Resource - Statement IDs"
              placeholder="At least one Statement ID is required"
              required
              type="text"
              @input="(value) => { updateStatementId(index, value) }"
            />
            <b-input-group-append
              v-if="index == resourceStatementIds.length -1">
              <b-button
                variant="success"
                @click="addResourceStatementIDField('')">+</b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>

        <b-form-group
          id="transaction_from_date_group"
          label-for="transaction_from_date"
          label="Transaction From Date">
          <b-form-input
            id="transaction_from_date"
            v-model="transaction_from_date"
            :state="isNotEmpty(transaction_from_date)"
            placeholder="e.g. 2006-01-02T15:04:05Z07:00"
            required
            type="text"
          />
        </b-form-group>

        <b-form-group
          id="transaction_to_date_group"
          label-for="transaction_to_date"
          label="Transaction To Date">
          <b-form-input
            id="transaction_to_date"
            v-model="transaction_to_date"
            :state="isNotEmpty(transaction_to_date)"
            placeholder="e.g. 2006-01-02T15:04:05Z07:00"
            required
            type="text"
          />
        </b-form-group>

        <b-form-group
          id="client_id_group"
          label-for="client_id"
          label="Client ID">
          <b-form-input
            id="client_id"
            v-model="client_id"
            :state="isNotEmpty(client_id)"
            required
            type="text"
            placeholder="Enter your Client ID"
          />
        </b-form-group>

        <b-form-group
          id="client_secret_group"
          label-for="client_secret"
          label="Client Secret">
          <b-form-input
            id="client_secret"
            v-model="client_secret"
            :state="isNotEmpty(client_secret)"
            required
            type="text"
            placeholder="Enter your Client Secret"
          />
        </b-form-group>
        <b-form-group
          id="x_fapi_financial_id_group"
          label-for="x_fapi_financial_id"
          label="x-fapi-financial-id"
          description="The unique id of the ASPSP to which the request is issued. The unique id will be issued by OB."
        >
          <b-form-input
            id="x_fapi_financial_id"
            v-model="x_fapi_financial_id"
            :state="isNotEmpty(x_fapi_financial_id)"
            placeholder="Enter your x-fapi-financial-id"
            required
            type="text"
          />
        </b-form-group>
      </b-card>
      <br>
      <b-card bg-variant="light">
        <b-form-group
          label="Well-Known"
          label-size="lg"/>
        <b-form-group
          id="token_endpoint_group"
          label-for="token_endpoint"
          label="Token Endpoint">
          <b-form-input
            id="token_endpoint"
            v-model="token_endpoint"
            :state="isValidUrl(token_endpoint)"
            required
            type="url"
          />
        </b-form-group>

        <b-form-group
          id="response_type_group"
          label-for="response_type"
          label="OAuth 2.0 response_type"
          description="REQUIRED. JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values"
        >
          <b-form-select
            id="response_type"
            v-model="response_type"
            :options="response_types_supported"
            :state="isNotEmpty(response_type)"
            required
          />
        </b-form-group>

        <b-form-group
          id="token_endpoint_auth_method_group"
          label-for="token_endpoint_auth_method"
          label="Token Endpoint Auth Method"
          description="Registered client authentication method, e.g client_secret_basic"
        >
          <b-form-select
            id="token_endpoint_auth_method"
            v-model="token_endpoint_auth_method"
            :options="token_endpoint_auth_methods"
            :state="true"
            required
          />
        </b-form-group>

        <b-form-group
          id="request_object_signing_alg_group"
          label-for="request_object_signing_alg"
          label="Request object signing algorithm"
          description="Algorithm used to sign requests objects"
        >
          <b-form-select
            id="request_object_signing_alg"
            v-model="request_object_signing_alg"
            :options="request_object_signing_alg_values_supported"
            :state="isNotEmpty(request_object_signing_alg)"
            required
          />
        </b-form-group>

        <b-form-group
          id="authorization_endpoint_group"
          label-for="authorization_endpoint"
          label="Authorization Endpoint">
          <b-form-input
            id="authorization_endpoint"
            v-model="authorization_endpoint"
            :state="isValidUrl(authorization_endpoint)"
            required
            type="url"
          />
        </b-form-group>

        <b-form-group
          id="resource_base_url_group"
          label-for="resource_base_url"
          label="Resource Base URL">
          <b-form-input
            id="resource_base_url"
            v-model="resource_base_url"
            :state="isValidUrl(resource_base_url)"
            required
            type="url"
          />
        </b-form-group>

        <b-form-group
          id="issuer_group"
          label-for="issuer"
          label="Issuer">
          <b-form-input
            id="issuer"
            v-model="issuer"
            :state="isValidUrl(issuer)"
            required
            type="url"
          />
        </b-form-group>

        <b-form-group
          id="redirect_url_group"
          label-for="redirect_url"
          label="Redirect URL">
          <b-form-input
            id="redirect_url"
            v-model="redirect_url"
            :state="isValidUrl(redirect_url)"
            required
            type="url"
          />
        </b-form-group>
      </b-card>
      <br>
      <b-card bg-variant="light">
        <b-form-group
          id="creditor_account_group"
          label-for="creditor_account"
          label="CreditorAccount"
          description="OBCashAccount5">
          <b-form-group
            id="creditor_account_scheme_name_group"
            label-for="creditor_account_scheme_name"
            label="SchemeName"
            description="OBExternalAccountIdentification4Code">
            <b-form-select
              id="creditor_account_scheme_name"
              v-model="creditor_account.scheme_name"
              :options="[
                'UK.OBIE.BBAN',
                'UK.OBIE.IBAN',
                'UK.OBIE.PAN' ,
                'UK.OBIE.Paym',
                'UK.OBIE.SortCodeAccountNumber'
              ]"
              required/>
          </b-form-group>
          <b-form-group
            id="creditor_account_identification_group"
            label-for="creditor_account_identification"
            label="Identification"
            description="Beneficiary account identification">
            <b-form-input
              id="creditor_account_identification"
              v-model="creditor_account.identification"
              :state="isNotEmpty(creditor_account.identification)"
              required
            />
          </b-form-group>
          <b-form-group
            id="creditor_account_name_group"
            label-for="creditor_account_name"
            label="Name"
            description="Name of the account, as assigned by the account servicing institution.\nUsage: The account name is the name or names of the account owner(s) represented at an account level. The account name is not the product name or the nickname of the account.">
            <b-form-input
              id="creditor_account_name"
              v-model="creditor_account.name"
              :state="isNotEmpty(creditor_account.name)"
              required
            />
          </b-form-group>
          <b-form-group
            id="instructed_amount_value"
            label-for="instructed_amount_value"
            label="Instructed Amount Value"
            description="Value of the instructed amount.">
            <b-form-input
              id="instructed_amount_value"
              v-model="instructed_amount.value"
              :state="isNotEmpty(instructed_amount.value)"
              required
            />
          </b-form-group>
          <b-form-group
            id="instructed_amount_currency"
            label-for="instructed_amount_currency"
            label="Instructed Amount Currency"
            description="Instructed amount currency.">
            <b-form-select
              id="instructed_amount_currency"
              v-model="instructed_amount.currency"
              :options="[
                'USD',
                'EUR',
                'JPY',
                'GBP',
                'AUD',
                'CAD',
                'CHF',
                'CNY',
                'SEK',
                'NZD',
                'MXN',
                'SGD',
                'HKD',
                'NOK',
                'KRW',
                'TRY',
                'RUB',
                'INR',
                'BRL'
              ]"
              required/>
          </b-form-group>
          <b-form-group
            id="currency_of_transfer"
            label-for="currency_of_transfer"
            label="Currency Of Transfer"
            description="Currency Of Transfer.">
            <b-form-select
              id="currency_of_transfer"
              v-model="currency_of_transfer"
              :options="[
                'USD',
                'EUR',
                'JPY',
                'GBP',
                'AUD',
                'CAD',
                'CHF',
                'CNY',
                'SEK',
                'NZD',
                'MXN',
                'SGD',
                'HKD',
                'NOK',
                'KRW',
                'TRY',
                'RUB',
                'INR',
                'BRL'
              ]"
              required/>
          </b-form-group>
        </b-form-group>
      </b-card>
    </b-form>
  </div>
</template>

<script>
import { createNamespacedHelpers, mapActions } from 'vuex';
import isEmpty from 'lodash/isEmpty';

import ConfigurationFormFile from './ConfigurationFormFile.vue';

const { mapGetters } = createNamespacedHelpers('config');

export default {
  name: 'TheConfigurationForm',
  components: {
    ConfigurationFormFile,
  },
  computed: {
    ...mapGetters([
      'resourceAccountIds',
      'resourceStatementIds',
    ]),
    token_endpoint_auth_methods() {
      const authMethods = this.$store.state.config.token_endpoint_auth_methods;
      return authMethods.map(m => ({
        value: m,
        text: m,
      }));
    },
    transaction_from_date: {
      get() {
        return this.$store.state.config.configuration.transaction_from_date;
      },
      set(value) {
        this.$store.commit('config/SET_TRANSACTION_FROM_DATE', value);
      },
    },
    transaction_to_date: {
      get() {
        return this.$store.state.config.configuration.transaction_to_date;
      },
      set(value) {
        this.$store.commit('config/SET_TRANSACTION_TO_DATE', value);
      },
    },
    // For an explanation on how these work. See:
    // * https://stackoverflow.com/a/45841419/241993
    // * http://shzhangji.com/blog/2018/04/17/form-handling-in-vuex-strict-mode/#Computed-Property
    client_id: {
      get() {
        return this.$store.state.config.configuration.client_id;
      },
      set(value) {
        this.$store.commit('config/SET_CLIENT_ID', value);
      },
    },

    client_secret: {
      get() {
        return this.$store.state.config.configuration.client_secret;
      },
      set(value) {
        this.$store.commit('config/SET_CLIENT_SECRET', value);
      },
    },

    token_endpoint: {
      get() {
        return this.$store.state.config.configuration.token_endpoint;
      },
      set(value) {
        this.$store.commit('config/SET_TOKEN_ENDPOINT', value);
      },
    },

    response_types_supported: {
      get() {
        return this.$store.state.config.response_types_supported;
      },
    },
    response_type: {
      get() {
        return this.$store.state.config.configuration.response_type;
      },
      set(value) {
        this.$store.commit('config/SET_RESPONSE_TYPE', value);
      },
    },
    token_endpoint_auth_method: {
      get() {
        return this.$store.state.config.configuration.token_endpoint_auth_method;
      },
      set(value) {
        this.$store.commit('config/SET_TOKEN_ENDPOINT_AUTH_METHOD', value);
      },
    },
    request_object_signing_alg_values_supported: {
      get() {
        return this.$store.state.config.request_object_signing_alg_values_supported;
      },
    },
    request_object_signing_alg: {
      get() {
        return this.$store.state.config.configuration.request_object_signing_alg;
      },
      set(value) {
        this.$store.commit('config/SET_REQUEST_OBJECT_SIGNING_ALG', value);
      },
    },
    token_endpoint_auth_signing_alg_values_supported: {
      get() {
        return this.$store.state.config.configuration.token_endpoint_auth_signing_alg_values_supported;
      },
    },
    token_endpoint_auth_signing_alg: {
      get() {
        return this.$store.state.config.configuration.token_endpoint_auth_signing_alg;
      },
      set(value) {
        this.$store.commit('config/SET_TOKEN_ENDPOINT_AUTH_SIGNING_ALG', value);
      },
    },
    id_token_signing_alg_values_supported: {
      get() {
        return this.$store.state.config.configuration.id_token_signing_alg_values_supported;
      },
    },
    id_token_signing_alg: {
      get() {
        return this.$store.state.config.configuration.id_token_signing_alg;
      },
      set(value) {
        this.$store.commit('config/SET_ID_TOKEN_SIGNING_ALG', value);
      },
    },
    authorization_endpoint: {
      get() {
        return this.$store.state.config.configuration.authorization_endpoint;
      },
      set(value) {
        this.$store.commit('config/SET_AUTHORIZATION_ENDPOINT', value);
      },
    },

    resource_base_url: {
      get() {
        return this.$store.state.config.configuration.resource_base_url;
      },
      set(value) {
        this.$store.commit('config/SET_RESOURCE_BASE_URL', value);
      },
    },

    x_fapi_financial_id: {
      get() {
        return this.$store.state.config.configuration.x_fapi_financial_id;
      },
      set(value) {
        this.$store.commit('config/SET_X_FAPI_FINANCIAL_ID', value);
      },
    },

    issuer: {
      get() {
        return this.$store.state.config.configuration.issuer;
      },
      set(value) {
        this.$store.commit('config/SET_ISSUER', value);
      },
    },

    redirect_url: {
      get() {
        return this.$store.state.config.configuration.redirect_url;
      },
      set(value) {
        this.$store.commit('config/SET_REDIRECT_URL', value);
      },
    },

    creditor_account: {
      get() {
        const self = this;
        return {
          get scheme_name() {
            return self.$store.state.config.configuration.creditor_account.scheme_name;
          },
          set scheme_name(value) {
            self.$store.commit('config/SET_CREDITOR_ACCOUNT_NAME_SCHEME_NAME', value);
          },
          get identification() {
            return self.$store.state.config.configuration.creditor_account.identification;
          },
          set identification(value) {
            self.$store.commit('config/SET_CREDITOR_ACCOUNT_IDENTIFICATION', value);
          },
          get name() {
            return self.$store.state.config.configuration.creditor_account.name;
          },
          set name(value) {
            self.$store.commit('config/SET_CREDITOR_ACCOUNT_NAME', value);
          },
        };
      },
    },

    instructed_amount: {
      get() {
        const self = this;
        return {
          get value() {
            return self.$store.state.config.configuration.instructed_amount.value;
          },
          set value(value) {
            self.$store.commit('config/SET_INSTRUCTED_AMOUNT_VALUE', value);
          },
          get currency() {
            return self.$store.state.config.configuration.instructed_amount.currency;
          },
          set currency(currency) {
            self.$store.commit('config/SET_INSTRUCTED_AMOUNT_CURRENCY', currency);
          },
        };
      },
    },

    currency_of_transfer: {
      get() {
        return this.$store.state.config.configuration.currency_of_transfer;
      },
      set(value) {
        this.$store.commit('config/SET_CURRENCY_OF_TRANSFER', value);
      },
    },
  },
  methods: {
    ...mapActions('config', [
      'removeResourceAccountID',
      'removeResourceStatementID',
      'setResourceAccountID',
      'setResourceStatementID',
    ]),
    updateAccountId(index, value) {
      this.setResourceAccountID({ index, value });
    },
    updateStatementId(index, value) {
      this.setResourceStatementID({ index, value });
    },
    isNotEmpty(value) {
      return !isEmpty(value);
    },
    isValidUrl(value) {
      try {
        return Boolean(new URL(value));
      } catch (e) {
        return false;
      }
    },
    addResourceAccountIDField(value) {
      this.$store.commit('config/ADD_RESOURCE_ACCOUNT_ID', { account_id: value });
    },
    removeResourceAccountIDField(index) {
      this.removeResourceAccountID(index);
    },
    addResourceStatementIDField(value) {
      this.$store.commit('config/ADD_RESOURCE_STATEMENT_ID', { statement_id: value });
    },
    removeResourceStatementIDField(index) {
      this.removeResourceStatementID(index);
    },
  },
};
</script>

<style scoped>
</style>
