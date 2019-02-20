<template>
  <div class="d-flex flex-column flex-fill p-3">
    <div class="centered">
      <h1>404</h1>
      <h2>{{ $route.path }} Not Found</h2>
      <hr>
      <router-link to="/">Go to the home page</router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NotFound',
};
</script>

<style scoped>
</style>
