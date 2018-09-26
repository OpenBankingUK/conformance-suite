<template>
  <div class="config">
    <div class="config-header">
      <h2>Config</h2>
      <a-button @click="handleSwitch">Use {{ useJson ? 'Web Form' : 'JSON' }}</a-button>
    </div>
    <a-divider />
    <div v-if="useJson">
      <editor
        name="editor"
        :value="config"
        :onChange="handleSetConfig" />
      <h2>Payload</h2>
      <editor
        name="payload"
        :value="payload"
        :onChange="handleSetPayload" />
      <div class="validation-button-container">
        <a-button
          type="primary"
          size="large"
          @click="startValidation"
          class="start_validation">
          Start validation
        </a-button>
      </div>
    </div>
    <div v-else>
      <a-steps :current="current">
        <a-step
          v-for="item in steps"
          :key="item.title"
          :title="item.title" />
      </a-steps>
      <Form
        :handleOnChange="handleOnChange"
        :handleSubmit="handleSubmit"
        :current="current"
        :length="steps.length"
        :next="next"
        :prev="prev"
        :values="{ config, payload }"
        :active="steps[current].filter" />
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import Editor from './Config/Editor';
import Form from './Config/Form';

export default {
  data() {
    return {
      useJson: 0,
      current: 0,
      steps: [{
        title: 'ASPSP',
        filter: ['authorization_endpoint', 'fapi_financial_id', 'issuer', 'resource_endpoint', 'token_endpoint'],
      }, {
        title: 'TPP Registered',
        filter: ['token_endpoint_auth_method', 'client_id', 'client_secret', 'redirect_uri'],
      }, {
        title: 'TPP Keys',
        filter: ['signing_key', 'signing_kid', 'transport_cert', 'transport_key'],
      }, {
        title: 'Payload',
      }],
    };
  },
  components: {
    Editor,
    Form,
  },
  computed: {
    ...mapGetters('config', {
      config: 'getConfig',
      payload: 'getPayload',
    }),
  },
  methods: {
    ...mapActions('config', ['setConfig', 'setPayload', 'startValidation']),
    next() {
      this.current += 1;
    },
    prev() {
      this.current -= 1;
    },
    handleSwitch() {
      this.useJson = !this.useJson;
      if (!this.useJson && this.form) {
        this.form.validateFields();
      }
    },
    isValidJSON(json) {
      try {
        JSON.parse(json);
      } catch (e) {
        return false;
      }
      return true;
    },
    handleSubmit() {
      return this.startValidation();
    },
    handleOnChange({ target }) {
      this.setConfig({ ...this.config, [target.name]: target.value });
    },
    handleSetConfig(config) {
      if (!this.isValidJSON(config)) return;
      this.setConfig(JSON.parse(config));
    },
    handleSetPayload(payload) {
      if (!this.isValidJSON(payload)) return;
      this.setPayload(JSON.parse(payload));
    },
  },
};
</script>

<style>
.validation-button-container {
  text-align: right;
}
.config-header {
  display: flex;
  justify-content: space-between;
}
.config label {
  text-transform: capitalize;
}
</style>
