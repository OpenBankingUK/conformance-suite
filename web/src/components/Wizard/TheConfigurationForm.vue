<template>
  <div class="p-3">
    <b-form>
      <b-card bg-variant="light">
        <b-form-group
          label="Client"
          label-size="lg" />
        <ConfigurationFormFile
          id="signing_private"
          setter-method-name-suffix="SigningPrivate"
          label="Private Signing Key (.key):"
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
          label="Private Transport Key (.key):"
          validExtension=".key"
        />
        <ConfigurationFormFile
          id="transport_public"
          setter-method-name-suffix="TransportPublic"
          label="Public Transport Certificate (.pem):"
          validExtension=".pem"
        />

        <b-form-group
          id="tpp_signature_kid"
          label-for="tpp_signature_kid"
          label="Client (TPP) Signature KID">
          <b-form-input
            id="tpp_signature_kid"
            v-model="tpp_signature_kid"
            :state="isNotEmpty(tpp_signature_kid)"
            required
            type="text"
          />
        </b-form-group>

        <b-form-group
          id="tpp_signature_issuer"
          label-for="tpp_signature_issuer"
          label="Client (TPP) Signature Issuer">
          <b-form-input
            id="tpp_signature_issuer"
            v-model="tpp_signature_issuer"
            :state="isNotEmpty(tpp_signature_issuer)"
            required
            type="text"
          />
        </b-form-group>

        <b-form-group
          id="tpp_signature_tan"
          label-for="tpp_signature_tan"
          label="Client (TPP) Signature Trust Anchor">
          <b-form-input
            id="tpp_signature_tan"
            v-model="tpp_signature_tan"
            :state="isNotEmpty(tpp_signature_tan)"
            required
            type="text"
          />
        </b-form-group>

        <b-form-group
          id="resource_account_id_group"
          label-for="resource_account_ids"
          label="Account IDs"
        >
          <b-input-group
            v-for="(item, index) in resourceAccountIds"
            :key="index"
            class="mt-3">
            <b-input-group-prepend v-if="resourceAccountIds.length > 1">
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
            <b-input-group-append v-if="index == resourceAccountIds.length -1">
              <b-button
                variant="success"
                @click="addResourceAccountIDField('')">+</b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>

        <b-form-group
          id="resource_statement_id_group"
          label-for="resource_statement_ids"
          label="Statement IDs"
        >
          <b-input-group
            v-for="(item, index) in resourceStatementIds"
            :key="index"
            class="mt-3">
            <b-input-group-prepend v-if="resourceStatementIds.length > 1">
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
            <b-input-group-append v-if="index == resourceStatementIds.length -1">
              <b-button
                variant="success"
                @click="addResourceStatementIDField('')">+</b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>

        <b-form-group
          id="transaction_from_date_group"
          label-for="transaction_from_date"
          label="Transaction From Date"
        >
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
          label="Transaction To Date"
        >
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
          v-if="client_secret_visible()"
          id="client_secret_group"
          label-for="client_secret"
          label="Client Secret"
        >
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
        <b-form-group
          id="send_x_fapi_customer_ip_address_group"
          label-for="send_x_fapi_customer_ip_address"
          label="Send x-fapi-customer-ip-address header"
        >
          <b-form-checkbox
            id="send_x_fapi_customer_ip_address"
            v-model="send_x_fapi_customer_ip_address"
          />
        </b-form-group>

        <b-form-group
          v-if="send_x_fapi_customer_ip_address"
          id="x_fapi_customer_ip_address_group"
          label-for="x_fapi_customer_ip_address"
          label="x-fapi-customer-ip-address"
          description="The IP address of the logged in PSU. Providing this HTTP header infers that the PSU is present during the interaction."
        >
          <b-form-input
            id="x_fapi_customer_ip_address"
            v-model="x_fapi_customer_ip_address"
            placeholder="x-fapi-customer-ip-address"
            type="text"
          />
        </b-form-group>
      </b-card>
      <br >
      <b-card bg-variant="light">
        <b-form-group
          label="Well-Known"
          label-size="lg" />
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
          label="Authorization Endpoint"
        >
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
          label="Resource Base URL"
        >
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

      <br >

      <b-card bg-variant="light">
        <b-form-group
          label="Payments"
          label-size="lg" />

        <DateTimeISO8601
          id="first_payment_date_time"
          field="first_payment_date_time"
          label="First Payment Date Time"
          description="First Payment Date Time"
          mutation="SET_FIRST_PAYMENT_DATE_TIME"
        />
        <DateTimeISO8601
          id="requested_execution_date_time"
          field="requested_execution_date_time"
          label="Requested Execution Date Time"
          description="Requested Execution Date Time"
          mutation="SET_REQUESTED_EXECUTION_DATE_TIME"
        />

        <b-form-group
          id="creditor_account_group"
          label-for="creditor_account"
          label="CreditorAccount"
          description="OBCashAccount5"
        >
          <SchemeName creditorAccountType="Local" />
          <b-form-group
            id="creditor_account_identification_group"
            label-for="creditor_account_identification"
            label="Identification"
            description="Beneficiary account identification"
          >
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
            description="Name of the account, as assigned by the account servicing institution.\nUsage: The account name is the name or names of the account owner(s) represented at an account level. The account name is not the product name or the nickname of the account."
          >
            <b-form-input
              id="creditor_account_name"
              v-model="creditor_account.name"
              :state="isNotEmpty(creditor_account.name)"
              required
            />
          </b-form-group>

          <SchemeName creditorAccountType="International" />
          <b-form-group
            id="international_creditor_account_identification_group"
            label-for="international_creditor_account_identification"
            label="International Identification"
            description="International beneficiary account identification"
          >
            <b-form-input
              id="international_creditor_account_identification"
              v-model="international_creditor_account.identification"
              :state="isNotEmpty(international_creditor_account.identification)"
              required
            />
          </b-form-group>
          <b-form-group
            id="international_creditor_account_name_group"
            label-for="international_creditor_account_name"
            label="International Name"
            description="International name of the account, as assigned by the account servicing institution.\nUsage: The account name is the name or names of the account owner(s) represented at an account level. The account name is not the product name or the nickname of the account."
          >
            <b-form-input
              id="international_creditor_account_name"
              v-model="international_creditor_account.name"
              :state="isNotEmpty(international_creditor_account.name)"
              required
            />
          </b-form-group>

          <b-form-group
            id="instructed_amount_value_group"
            label-for="instructed_amount_value"
            label="Instructed Amount Value (Capped at 1.00)"
            description="Value of the instructed amount (^\d{1,13}\.\d{1,5}$)."
          >
            <b-form-input
              id="instructed_amount_value"
              v-model="instructed_amount.value"
              :state="isNotEmpty(instructed_amount.value)"
              required
            />
          </b-form-group>
          <b-form-group
            id="instructed_amount_currency_group"
            label-for="instructed_amount_currency"
            label="Instructed Amount Currency"
            description="Instructed amount currency (^[A-Z]{3,3}$)."
          >
            <b-form-select
              id="instructed_amount_currency"
              v-model="instructed_amount.currency"
              :options="top_20_currencies"
              required
            />
          </b-form-group>
          <b-form-group
            id="currency_of_transfer_group"
            label-for="currency_of_transfer"
            label="Currency Of Transfer For International Payments"
            description="Currency Of Transfer."
          >
            <b-form-select
              id="currency_of_transfer"
              v-model="currency_of_transfer"
              :options="top_20_currencies"
              required
            />
          </b-form-group>
        </b-form-group>

        <PaymentFrequency />
      </b-card>

      <br >

      <b-card bg-variant="light">
        <b-form-group
          label="Confirmation Of Funds"
          label-size="lg" />

        <SchemeName creditorAccountType="CBPII" />

        <b-form-group
          id="cbpii_debtor_account_identification_group"
          label-for="cbpii_debtor_account_identification"
          label="Debtor Account Identification"
          description="Debtor Account Identification"
        >
          <b-form-input
            id="cbpii_debtor_account_identification"
            v-model="cbpii_debtor_account.identification"
            :state="isNotEmpty(cbpii_debtor_account.identification)"
            required
          />
        </b-form-group>
        <b-form-group
          id="cbpii_debtor_account_name_group"
          label-for="cbpii_debtor_account_name"
          label="Debtor Account Name"
          description="Name of the account, as assigned by the account servicing institution"
        >
          <b-form-input
            id="cbpii_debtor_account_name"
            v-model="cbpii_debtor_account.name"
            :state="isNotEmpty(cbpii_debtor_account.name)"
            required
          />
        </b-form-group>

      </b-card>

      <br >

      <b-card
        v-if="conditional_properties && conditional_properties.length > 0"
        bg-variant="light">
        <b-form-group
          label="Conditional Properties"
          label-size="lg" />
        <b-card bg-variant="default">
          <b-form-group
            v-for="(property, propertyKey) in conditional_properties"
            :key="property.name"
            :label="property.name"
            label-size="lg" >
            <b-card bg-variant="light">
              <b-form-group
                v-for="(endpoint, endpointKey) in property.endpoints"
                :key="endpoint.name"
                :label="`${endpoint.method} ${endpoint.path}`"
                label-size="lg" >
                <b-form-group>
                  <div>
                    <b-row
                      v-for="(conditionalProperty, conditionalPropertyKey) in endpoint.conditionalProperties"
                      :key="conditionalProperty.name"
                      striped
                      hover
                      class="my-1 p-3">
                      <b-col sm="12">
                        <label><b>Schema:</b> {{ conditionalProperty.schema }}</label>
                      </b-col>
                      <b-col sm="12">
                        <label><b>Name:</b> {{ conditionalProperty.name }}</label>
                      </b-col>
                      <b-col sm="12">
                        <label><b>Path:</b> {{ conditionalProperty.path }}</label>
                      </b-col>
                      <b-col sm="12">
                        <b-form-input
                          :placeholder="getConditionalPropertyPlaceholderValue(conditional_properties[propertyKey].endpoints[endpointKey].conditionalProperties[conditionalPropertyKey].required)"
                          v-model="conditional_properties[propertyKey].endpoints[endpointKey].conditionalProperties[conditionalPropertyKey].value"/>
                      </b-col>
                      <br >
                    </b-row>
                  </div>
                </b-form-group>
              </b-form-group>
            </b-card>
          </b-form-group>
        </b-card>
      </b-card>
    </b-form>
  </div>
</template>

<script>
import { createNamespacedHelpers, mapActions } from 'vuex';
import isEmpty from 'lodash/isEmpty';
import ConfigurationFormFile from './ConfigurationFormFile.vue';
import PaymentFrequency from '../config/PaymentFrequency.vue';
import SchemeName from '../config/SchemeName.vue';
import DateTimeISO8601 from '../config/DateTimeISO8601.vue';
import api from '../../api/apiUtil';

const { mapGetters } = createNamespacedHelpers('config');

export default {
  name: 'TheConfigurationForm',
  components: {
    ConfigurationFormFile,
    PaymentFrequency,
    SchemeName,
    DateTimeISO8601,
  },
  data() {
    api.get('/api/config/conditional-property').then((res) => {
      res.json().then((body) => {
        if (this.$store.state.config.configuration.conditional_properties && this.$store.state.config.configuration.conditional_properties.length === 0) {
          this.$store.commit('config/SET_CONDITIONAL_PROPERTIES', body);
        }
      });
    });

    return {};
  },
  computed: {
    ...mapGetters(['resourceAccountIds', 'resourceStatementIds']),
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
    tpp_signature_kid: {
      get() {
        return this.$store.state.config.configuration.tpp_signature_kid;
      },
      set(value) {
        this.$store.commit('config/SET_TPP_SIGNATURE_KID', value);
      },
    },
    tpp_signature_issuer: {
      get() {
        return this.$store.state.config.configuration.tpp_signature_issuer;
      },
      set(value) {
        this.$store.commit('config/SET_TPP_SIGNATURE_ISSUER', value);
      },
    },
    tpp_signature_tan: {
      get() {
        return this.$store.state.config.configuration.tpp_signature_tan;
      },
      set(value) {
        this.$store.commit('config/SET_TPP_SIGNATURE_TAN', value);
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
    acr_values_supported: {
      get() {
        return this.$store.state.config.acr_values_supported;
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
        return this.$store.state.config.configuration
          .token_endpoint_auth_method;
      },
      set(value) {
        this.$store.commit('config/SET_TOKEN_ENDPOINT_AUTH_METHOD', value);
      },
    },
    request_object_signing_alg_values_supported: {
      get() {
        return this.$store.state.config
          .request_object_signing_alg_values_supported;
      },
    },
    request_object_signing_alg: {
      get() {
        return this.$store.state.config.configuration
          .request_object_signing_alg;
      },
      set(value) {
        this.$store.commit('config/SET_REQUEST_OBJECT_SIGNING_ALG', value);
      },
    },
    token_endpoint_auth_signing_alg_values_supported: {
      get() {
        return this.$store.state.config.configuration
          .token_endpoint_auth_signing_alg_values_supported;
      },
    },
    token_endpoint_auth_signing_alg: {
      get() {
        return this.$store.state.config.configuration
          .token_endpoint_auth_signing_alg;
      },
      set(value) {
        this.$store.commit('config/SET_TOKEN_ENDPOINT_AUTH_SIGNING_ALG', value);
      },
    },
    id_token_signing_alg_values_supported: {
      get() {
        return this.$store.state.config.configuration
          .id_token_signing_alg_values_supported;
      },
    },
    conditional_properties: {
      get() {
        return this.$store.state.config.configuration.conditional_properties;
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
    send_x_fapi_customer_ip_address: {
      get() {
        return this.$store.state.config.configuration
          .send_x_fapi_customer_ip_address;
      },
      set(value) {
        this.$store.commit('config/SET_SEND_X_FAPI_CUSTOMER_IP_ADDRESS', value);
      },
    },
    x_fapi_customer_ip_address: {
      get() {
        return this.$store.state.config.configuration
          .x_fapi_customer_ip_address;
      },
      set(value) {
        this.$store.commit('config/SET_X_FAPI_CUSTOMER_IP_ADDRESS', value);
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
          get identification() {
            return self.$store.state.config.configuration.creditor_account
              .identification;
          },
          set identification(value) {
            self.$store.commit(
              'config/SET_CREDITOR_ACCOUNT_IDENTIFICATION',
              value,
            );
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
    international_creditor_account: {
      get() {
        const self = this;
        return {
          get identification() {
            return self.$store.state.config.configuration
              .international_creditor_account.identification;
          },
          set identification(value) {
            self.$store.commit(
              'config/SET_INTERNATIONAL_CREDITOR_ACCOUNT_IDENTIFICATION',
              value,
            );
          },
          get name() {
            return self.$store.state.config.configuration
              .international_creditor_account.name;
          },
          set name(value) {
            self.$store.commit(
              'config/SET_INTERNATIONAL_CREDITOR_ACCOUNT_NAME',
              value,
            );
          },
        };
      },
    },
    cbpii_debtor_account: {
      get() {
        const self = this;
        return {
          get identification() {
            return self.$store.state.config.configuration.cbpii_debtor_account.identification;
          },
          set identification(value) {
            self.$store.commit(
              'config/SET_CBPII_DEBTOR_ACCOUNT_IDENTIFICATION',
              value,
            );
          },
          get scheme_name() {
            return self.$store.state.config.configuration.cbpii_debtor_account.scheme_name;
          },
          set scheme_name(value) {
            self.$store.commit(
              'config/SET_CBPII_DEBTOR_ACCOUNT_SCHEME_NAME',
              value,
            );
          },
          get name() {
            return self.$store.state.config.configuration.cbpii_debtor_account.name;
          },
          set name(value) {
            self.$store.commit(
              'config/SET_CBPII_DEBTOR_ACCOUNT_NAME',
              value,
            );
          },
        };
      },
    },
    instructed_amount: {
      get() {
        const self = this;
        return {
          get value() {
            return self.$store.state.config.configuration.instructed_amount
              .value;
          },
          set value(value) {
            self.$store.commit('config/SET_INSTRUCTED_AMOUNT_VALUE', value);
          },
          get currency() {
            return self.$store.state.config.configuration.instructed_amount
              .currency;
          },
          set currency(currency) {
            self.$store.commit(
              'config/SET_INSTRUCTED_AMOUNT_CURRENCY',
              currency,
            );
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
    top_20_currencies: {
      get() {
        return [
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
          'BRL',
        ];
      },
    },
    payment_frequency: {
      get() {
        return this.$store.state.config.configuration.payment_frequency;
      },
      set(value) {
        this.$store.commit('config/SET_PAYMENT_FREQUENCY', value);
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
    client_secret_visible() {
      return (
        this.$store.state.config.configuration.token_endpoint_auth_method
        === 'client_secret_basic'
      );
    },
    addResourceAccountIDField(value) {
      this.$store.commit('config/ADD_RESOURCE_ACCOUNT_ID', {
        account_id: value,
      });
    },
    removeResourceAccountIDField(index) {
      this.removeResourceAccountID(index);
    },
    addResourceStatementIDField(value) {
      this.$store.commit('config/ADD_RESOURCE_STATEMENT_ID', {
        statement_id: value,
      });
    },
    removeResourceStatementIDField(index) {
      this.removeResourceStatementID(index);
    },
    getConditionalPropertyPlaceholderValue(required) {
      return `Value ${required ? '(Required)' : '(Optional)'}`;
    },
  },
};
</script>

<style scoped>
</style>
