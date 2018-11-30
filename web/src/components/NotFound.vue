<template>
  <div class="centered">
    <h1>404</h1>
    <h2>{{ $route.path }} Not Found</h2>
    <a-divider />
    <router-link to="/">Go to the home page</router-link>
  </div>
</template>

<script>

export default {
  name: 'notfound',
};
</script>
