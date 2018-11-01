<template>
  <div>
    <!-- eslint-disable vue/this-in-template -->
    <a-form
      :autoFormCreate="(form)=>{this.form = form}"
      style="margin-top: 50px;">
      <a-form-item
        v-for="(row, el) in values.config"
        v-show="active.includes(el)"
        :key="el"
        :label="formatLabel(el)"
        :labelCol="{ span: 8 }"
        :wrapperCol="{ span: 12 }"
        :fieldDecoratorId="el"
        :fieldDecoratorOptions="{
          rules: [{ required: true, message: 'Field required!' }],
          initialValue: row
      }">
        <a-textarea
          v-if="['certificateSigning',
                 'certificateTransport',
                 'privateKeySigning',
                 'privateKeyTransport'].includes(el)"
          :name="el"
          :autosize="{ minRows: 2, maxRows: 6 }"
          @change="handleOnChange" />
        <a-input
          v-else
          :name="el"
          size="large"
          @change="handleOnChange" />
      </a-form-item>
    </a-form>
    <DiscoveryModel
      v-show="current == length - 1"
      :current="current"
      :discoveryModelValue="values.discoveryModel"
      :length="length" />
    <div class="validation-button-container">
      <a-button
        v-if="current > 0"
        size="large"
        @click="prev"
      >
        Previous
      </a-button>
      <a-button
        v-if="current < length - 1"
        :disabled="hasErrors()"
        type="primary"
        size="large"
        class="next"
        @click="next"
      >
        Next
      </a-button>
      <a-button
        v-if="current == length - 1"
        :disabled="hasErrors()"
        type="primary"
        size="large"
        class="start"
        @click="openNotification">
        Start validation
      </a-button>
    </div>
  </div>
</template>

<script>
import * as _ from 'lodash';
import DiscoveryModel from './DiscoveryModel.vue';

export default {
  components: {
    DiscoveryModel,
  },
  props: {
    handleOnChange: {
      type: Function,
      default: () => {},
    },
    handleSubmit: {
      type: Function,
      default: () => {},
    },
    values: {
      type: Object,
      default: null,
    },
    active: {
      type: Array,
      default: () => [],
    },
    current: {
      type: Number,
      default: 0,
    },
    length: {
      type: Number,
      default: 0,
    },
    next: {
      type: Function,
      default: () => {},
    },
    prev: {
      type: Function,
      default: () => {},
    },
  },
  mounted() {
    this.$nextTick(() => {
      this.form.validateFields();
    });
  },
  methods: {
    formatLabel(label) {
      return _.startCase(label);
      // return label.replace(/^[a-z]|[A-Z]/g, (v, i) =>
      //   (i === 0 ? v.toUpperCase() : ` ${v.toLowerCase()}`));
    },
    openNotification() {
      this.$notification.success({
        message: 'Valid config',
        description: 'Starting validation...',
        duration: 2,
      });
      this.handleSubmit();
    },
    hasErrors() {
      if (this.form) {
        const fieldsError = this.form.getFieldsError();
        return Object.keys(fieldsError).some(field => fieldsError[field]);
      }
      return true;
    },
  },
};
</script>
