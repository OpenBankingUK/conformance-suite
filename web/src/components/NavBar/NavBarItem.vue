<template>
  <b-nav-item :to="route">
    <b-badge
      :variant="isActiveRoute(route)"
      pill>{{ no }}</b-badge>
    <span>&nbsp;</span>
    <span class="label">{{ label }}</span>
  </b-nav-item>
</template>

<style scoped>
/* Make it match branding guidelines:
"A simple guide to the
Open Banking brand
Version 1  July 2017" */
.badge-primary,
.nav-link.active {
  background-color: #6180c3 !important;
}

.nav-link.active {
  opacity: 0.7;
}

.nav-link:not(.active) > .label {
  color: #6180c3;
}
</style>

<script>
export default {
  name: 'NavBarItem',
  // inheritAttrs: false,
  props: {
    route: {
      type: String,
      required: true,
      private: true,
    },
    no: {
      type: Number,
      required: true,
      private: true,
    },
    label: {
      type: String,
      required: true,
      private: true,
    },
  },
  methods: {
    // see: https://getbootstrap.com/docs/4.0/components/badge/#pill-badges
    // this method determines what colour badge to display
    // when it is the active route.
    isActiveRoute(routeName) {
      if (this.$route.path === routeName) {
        return 'success';
      }

      return 'primary';
    },
  },
};
</script>
