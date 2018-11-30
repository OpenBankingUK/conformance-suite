<template>
  <a-collapse accordion>
    <a-collapse-panel
      v-for="(item, key) in failures"
      :key="key">
      <template slot="header">
        <span class="error-message">{{ item.message }}</span>
      </template>
      <error-details
        v-for="(errors, key) in item.results"
        v-if="errors.length"
        :errors="errors"
        :type="key"
        :key="key" />
    </a-collapse-panel>
  </a-collapse>
</template>

<script>
import ErrorDetails from './ErrorDetails.vue';

export default {
  components: { ErrorDetails },
  props: {
    failures: {
      type: Array,
      default: () => ([]),
    },
  },
};
</script>

<style>
.error-message {
  color: #cc0000;
}
</style>
