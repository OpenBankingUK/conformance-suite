<template>
  <div
    v-if="hasErrors"
    id="error-status">
    <b-alert
      v-for="(err, index) in errorMessages"
      :key="index"
      show
      variant="danger">{{ err }}</b-alert>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('status');

export default {
  name: 'TheErrorStatus',
  computed: {
    ...mapGetters([
      'errorMessages',
      'hasErrors',
    ]),
  },
};
</script>

<style scoped>
</style>
