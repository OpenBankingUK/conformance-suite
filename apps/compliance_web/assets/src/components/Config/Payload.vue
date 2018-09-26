<template>
  <div>
    <h2>Accounts</h2>
    <div
      v-for="(item, key) in accounts"
      :key="item.api_version + key"
      class="payload-item">
      <a-list>
        <a-list-item>
          <a-list-item-meta>
            <h3 slot="title">Api version: {{ item.api_version }}</h3>
          </a-list-item-meta>
          <a-button
            type="danger"
            size="small"
            shape="circle"
            @click="handleDelete(item)"
            :data-item="`account-${key}`"
            icon="minus" />
        </a-list-item>
      </a-list>
    </div>
    <div class="payload-item add-account">
      <add-new
        type="accounts" />
    </div>
    <h2>Payments</h2>
    <div
      class="payload-item"
      v-for="(item, key) in payments"
      :key="item.account_number + key">
      <a-list>
        <a-list-item>
          <a-list-item-meta>
            <h3 slot="title">Api version: {{ item.api_version }}</h3>
            <ul slot="description">
              <li>Name: <strong>{{ item.name }}</strong></li>
              <li>Account Number: <strong>{{ item.account_number }}</strong></li>
              <li>Sort Code: <strong>{{ item.sort_code }}</strong></li>
              <li>Amount: <strong>{{ item.amount }}</strong></li>
            </ul>
          </a-list-item-meta>
          <a-button
            type="danger"
            size="small"
            shape="circle"
            :data-item="`payment-${key}`"
            @click="handleDelete(item)"
            icon="minus" />
        </a-list-item>
      </a-list>
    </div>
    <div class="payload-item add-payment">
      <add-new
        type="payments" />
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import AddNew from './PayloadAddNew';

export default {
  components: {
    AddNew,
  },
  props: {
    current: {
      type: Number,
      default: 0,
    },
    length: {
      type: Number,
      default: 0,
    },
    payload: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    accounts() {
      return this.payload.filter(({ type }) => type === 'accounts');
    },
    payments() {
      return this.payload.filter(({ type }) => type === 'payments');
    },
  },
  methods: {
    ...mapActions('config', ['deletePayload']),
    handleDelete(item) {
      this.deletePayload(item);
    },
  },
};
</script>

<style>
.payload-item {
  margin-bottom: 10px;
  border-radius: 3px;
  background: #fafafa;
  padding: 5px 10px;
}
.payload-item .ant-list-item {
  align-items: flex-start;
}
.payload-item ul {
  margin: 0;
}
</style>
