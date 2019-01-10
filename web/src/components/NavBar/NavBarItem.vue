<template>
  <b-nav-item
    :to="route"
    :disabled="isRouteDisabled"
    exact>
    <b-badge
      :variant="computeVariant"
      pill>{{ no }}</b-badge>
    <span>&nbsp;</span>
    <span class="label">{{ label }}</span>
  </b-nav-item>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('config');

export default {
  name: 'NavBarItem',
  props: {
    route: {
      type: String,
      required: true,
    },
    no: {
      type: Number,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
  },
  computed: {
    ...mapGetters([
      'navigation',
    ]),
    /**
     * Determines what colour badge to display when it is the active route.
     * See: https://getbootstrap.com/docs/4.0/components/badge/#pill-badges
     */
    computeVariant() {
      if (this.$route.path === this.route) {
        return 'success';
      }

      return 'primary';
    },
    /**
     * Determines if this route is disabled.
     */
    isRouteDisabled() {
      const disabled = !this.navigation[this.route];
      return disabled;
    },
  },
};
</script>

<style scoped>
/**
 * Make it match branding guidelines: "A simple guide to the Open Banking brand Version 1 July 2017"
 */
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

.nav-link.disabled {
  opacity: 0.35;
}
</style>
