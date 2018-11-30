<template>
  <div class="config">
    <div class="config-header">
      <h2>Config</h2>
      <a-button @click="handleSwitch">Use {{ useJson ? 'Web Form' : 'JSON' }}</a-button>
    </div>
    <a-divider/>
    <div v-if="useJson">
      <div class="validation-button-container">
        <a-button
          type="primary"
          size="large"
          class="start_validation"
          @click="startValidation"
        >Start validation</a-button>
      </div>
      <editor
        :value="config"
        :onChange="handleSetConfig"
        name="editor"/>
      <h2>Discovery Model</h2>
      <editor
        :value="discoveryModel"
        :onChange="handleSetDiscoveryModel"
        name="discoveryModel"/>
    </div>
    <div v-else>
      <a-steps :current="current">
        <a-step
          v-for="item in steps"
          :key="item.title"
          :title="item.title"/>
      </a-steps>
      <Form
        :handleOnChange="handleOnChange"
        :handleSubmit="handleSubmit"
        :current="current"
        :length="steps.length"
        :next="next"
        :prev="prev"
        :values="{ config, discoveryModel }"
        :active="steps[current].filter"
      />
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
        title: 'Configuration',
        filter: [
          'accountAccessToken',
          'certificateSigning',
          'certificateTransport',
          'clientScopes',
          'keyId',
          'privateKeySigning',
          'privateKeyTransport',
          'softwareStatementId',
          'targetHost',
        ],
      }, {
        title: 'DiscoveryModel',
      }],
    };
  },
  computed: {
    ...mapGetters('config', {
      config: 'getConfig',
      discoveryModel: 'getDiscoveryModel',
    }),
  },
  methods: {
    ...mapActions('config', ['setConfig', 'setDiscoveryModel', 'startValidation']),
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
    handleSetDiscoveryModel(discoveryModel) {
      if (!this.isValidJSON(discoveryModel)) return;
      this.setDiscoveryModel(JSON.parse(discoveryModel));
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
