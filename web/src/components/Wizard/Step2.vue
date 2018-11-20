<template>
  <div class="d-flex flex-column flex-fill">
    <h3>Configuration: API Discovery</h3>
    <div class="p-2 flex-fill">
      contents
    </div>
    <div class="p-2 d-flex flex-column">
      <h3>Status</h3>
    </div>
  </div>
</template>

<style scoped>
</style>

<script>
export default {
  name: 'Step2',
  methods: {
    validate() {
      console.log('%s.validate=', this.$options.name);
    },
  },
};
</script>
