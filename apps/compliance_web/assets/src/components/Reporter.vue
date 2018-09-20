<template>
  <div class="reporter_container">
    <h2>Validation Run Report</h2>
    <div v-if="validationRunId && tests">
      <endpoint
        v-for="(value, key) in tests"
        :value="value"
        :endpoint="key"
        :key="key" />
      <div class="buttons_container">
        <a-button
          class="disconnect_button"
          type="danger"
          :disabled="connectionState !== 'CONNECTED'"
          @click="doDisconnect()">Disconnect</a-button>
      </div>
    </div>
    <spinner v-else/>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import Spinner from './Spinner';
import Endpoint from './Reporter/Endpoint';

const namespacedReporter = createNamespacedHelpers('reporter');
const namespacedValidations = createNamespacedHelpers('validations');

export default {
  name: 'reporter',
  components: {
    Spinner,
    Endpoint,
  },
  computed: {
    ...namespacedValidations.mapGetters([
      'validationRunId',
    ]),
    ...namespacedReporter.mapGetters([
      'connectionState',
      'lastUpdate',
      'tests',
    ]),
  },
  watch: {
    async validationRunId() {
      await this.connect();
      return this.subscribeToChannel(this.validationRunId);
    },
  },
  methods: {
    ...namespacedReporter.mapActions([
      'connect',
      'disconnect',
      'subscribeToChannel',
    ]),
    async doDisconnect() {
      return this.disconnect(this.validationRunId);
    },
  },
  async beforeDestroy() {
    return this.disconnect(this.validationRunId);
  },
  beforeRouteEnter(to, from, next) {
    if (from.name === 'Config') return next();
    return next('/config');
  },
};
</script>

<style>
.buttons_container {
  margin-top: 20px;
  text-align: right;
}
</style>
