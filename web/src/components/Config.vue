<template>
  <div class="config">
    <div class="config-header">
      <h2>Config</h2>
      <a-button @click="handleSwitch">Use {{ useJson ? 'Web Form' : 'JSON' }}</a-button>
    </div>
    <a-divider />
    <div v-if="useJson">
      <div class="validation-button-container">
        <a-button
          type="primary"
          size="large"
          class="start_validation"
          @click="startValidation">
          Start validation
        </a-button>
      </div>
      <editor
        :value="config"
        :onChange="handleSetConfig"
        name="editor" />
      <h2>Payload</h2>
      <editor
        :value="payload"
        :onChange="handleSetPayload"
        name="payload" />
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
import Editor from './Config/Editor.vue';
import Form from './Config/Form.vue';

export default {
  components: {
    Editor,
    Form,
  },
  data() {
    return {
      useJson: 1,
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
