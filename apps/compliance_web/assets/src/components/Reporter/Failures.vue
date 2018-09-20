<template>
  <a-collapse accordion>
    <a-collapse-panel
      v-for="(item, key) in failures"
      :key="key">
      <template slot="header">
        <span class="error-message">{{ item.message }}</span>
      </template>
      <error-details
        v-if="errors.length"
        v-for="(errors, key) in item.results"
        :errors="errors"
        :type="key"
        :key="key" />
    </a-collapse-panel>
  </a-collapse>
</template>

<script>
import ErrorDetails from './ErrorDetails';

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
