<template>
  <!-- eslint-disable vue/this-in-template -->
  <a-form
    :autoFormCreate="(form)=>{this.form = form}"
    layout="vertical"
    @submit="handleSubmit">
    <a-list>
      <a-list-item>
        <a-list-item-meta>
          <div slot="title">
            <h4>Add new {{ isAccount ? 'account' : 'payment' }}</h4>
          </div>
          <div slot="description">
            <a-form-item
              :fieldDecoratorOptions="{
              rules: [{ required: true, message: 'Api Version is required!' }]}"
              fieldDecoratorId="api_version">
              <a-select
                placeholder="Select api version..."
                style="width: 100%;">
                <a-select-option
                  v-for="val in apiVersions"
                  :key="val"
                  :value="val">
                  {{ val }}
                </a-select-option>
              </a-select>
            </a-form-item>
            <div v-if="!isAccount">
              <a-form-item
                v-for="field in fields"
                :key="field.key"
                :fieldDecoratorId="field.key"
                :label="field.label"
                :fieldDecoratorOptions="{
                rules: [{ required: true, message: `${field.label} is required!` }]}">
                <a-input :name="field.key" />
              </a-form-item>
            </div>
          </div>
        </a-list-item-meta>
        <a-button
          type="primary"
          htmlType="submit"
          size="small"
          shape="circle"
          icon="plus" />
      </a-list-item>
    </a-list>
  </a-form>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  props: {
    type: {
      type: String,
      default: 'accounts',
    },
    addNew: {
      type: Function,
      default: () => {},
    },
  },
  data() {
    return {
      fields: [
        { key: 'name', label: 'Full Name' },
        { key: 'sort_code', label: 'Sort Code' },
        { key: 'account_number', label: 'Account Number' },
        { key: 'amount', label: 'Amount' },
      ],
    };
  },
  computed: {
    isAccount() {
      return this.type === 'accounts';
    },
    apiVersions() {
      return this.isAccount
        ? ['1.1', '2.0', '3.0.0']
        : ['1.1', '3.0.0'];
    },
  },
  methods: {
    ...mapActions('config', ['updatePayload']),
    handleSubmit(e) {
      e.preventDefault();
      this.form.validateFields((err, values) => {
        if (!err) {
          this.updatePayload({
            ...values,
            type: this.type,
          });
          this.form.resetFields();
        }
      });
    },
  },
};
</script>
