<template>
  <div
    id="app"
    class="d-flex flex-column">
    <router-view/>
  </div>
</template>

<style>
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
