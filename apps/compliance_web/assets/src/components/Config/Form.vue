<template>
  <div>
    <a-form
      :autoFormCreate="(form)=>{this.form = form}"
      style="margin-top: 50px;">
      <a-form-item
        v-for="(row, el) in values.config"
        v-show="active.includes(el)"
        :key="el"
        :label="el.replace(/_/g, ' ')"
        :labelCol="{ span: 8 }"
        :wrapperCol="{ span: 12 }"
        :fieldDecoratorId="el"
        :fieldDecoratorOptions="{
          rules: [{ required: true, message: 'Field required!' }],
          initialValue: row
      }">
        <a-textarea
          v-if="['signing_key', 'transport_cert', 'transport_key'].includes(el)"
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
    <payload
      v-show="current == length - 1"
      :current="current"
      :payload="values.payload"
      :length="length" />
    <div class="validation-button-container">
      <a-button
        v-if="current > 0"
        @click="prev"
        size="large"
      >
        Previous
      </a-button>
      <a-button
        v-if="current < length - 1"
        type="primary"
        @click="next"
        size="large"
        :disabled="hasErrors()"
        class="next"
      >
        Next
      </a-button>
      <a-button
        type="primary"
        size="large"
        v-if="current == length - 1"
        @click="openNotification"
        class="start"
        :disabled="hasErrors()">
        Start validation
      </a-button>
    </div>
  </div>
</template>

<script>
import Payload from './Payload';

export default {
  components: {
    Payload,
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
  methods: {
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
  mounted() {
    this.$nextTick(() => {
      this.form.validateFields();
    });
  },
};
</script>
