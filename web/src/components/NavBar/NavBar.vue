<template>
  <div class="d-flex flex-column p-3">
    <div class="nav-section">
      <b-nav
        vertical
        pills>
        <NavBarItem
          :no="1"
          route="/wizard/step1"
          label="Start/Load Test" />
        <NavBarItem
          :no="2"
          route="/wizard/discovery-config"
          label="Discovery"
          exact/>
        <NavBarItem
          :no="3"
          route="/wizard/configuration"
          label="Configuration" />
        <NavBarItem
          :no="4"
          route="/wizard/run-overview"
          label="Run/Overview" />
        <NavBarItem
          :no="5"
          route="/wizard/summary"
          label="Summary" />
        <NavBarItem
          :no="6"
          route="/wizard/export"
          label="Export" />
      </b-nav>
    </div>
    <div class="nav-section">
      <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
        <span>Tools</span>
      </h6>
      <b-nav
        vertical
        class="mb-2">
        <b-nav-item>Documentation</b-nav-item>
        <b-nav-item>Bug Tracker</b-nav-item>
        <b-nav-item>Website</b-nav-item>
        <b-nav-item>Integrations</b-nav-item>
      </b-nav>
    </div>
  </div>
</template>

<style scoped>
.nav-section {
  background: #FFFFFF;
}
</style>

<script>
import NavBarItem from './NavBarItem.vue';

export default {
  name: 'NavBar',
  components: {
    NavBarItem,
  },
  // inheritAttrs: false,
  methods: {
  },
};
</script>
