<template>
  <div
    id="app"
    class="d-flex flex-column">
    <router-view/>
  </div>
</template>

<style>
div {
  border: 1px solid rgba(0, 0, 0, 0.09);
}
</style>

<script>
export default {
  name: 'app',
  components: {
  },
  computed: {
    layout() {
      return `${this.$route.meta.layout || 'default'}-layout`;
    },
  },
};
</script>
