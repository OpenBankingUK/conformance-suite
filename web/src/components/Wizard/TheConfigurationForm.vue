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
    </b-form>
  </div>
</template>

<script>
import isEmpty from 'lodash/isEmpty';

import ConfigurationFormFile from './ConfigurationFormFile.vue';

export default {
  name: 'TheConfigurationForm',
  components: {
    ConfigurationFormFile,
  },
  computed: {
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
  },
  methods: {
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
  },
};
</script>

<style scoped>
</style>
