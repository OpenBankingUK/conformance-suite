<template>
  <div class="d-flex flex-column flex-fill">
    <b>name:</b> {{ this.$options.name }}
  </div>
</template>

<style scoped>
</style>

<script>
export default {
  name: 'Step5',
  methods: {
    // Gets called by top-level Wizard component in the validateStep function.
    async validate() {
      const valid = true;
      // TODO: remove stub.
      return Promise.resolve(valid);
    },
  },
};
</script>
