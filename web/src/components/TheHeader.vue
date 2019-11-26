<template>
  <b-navbar>
    <b-navbar-brand to="/">Conformance Suite</b-navbar-brand>
    <TheNotificationBell />
  </b-navbar>
</template>

<script>
import TheNotificationBell from './TheNotificationBell.vue';

export default {
  name: 'TheHeader',
  components: {
    TheNotificationBell,
  },
};
</script>

<style scoped>
.navbar {
  background: linear-gradient(90deg, #6180c3, #6180c3);
  box-shadow: 0 6px 25px 0 rgba(38, 50, 56, 0.2);
}

.navbar-brand {
  color: #ffffff;
  font-size: 1rem;
}
</style>
