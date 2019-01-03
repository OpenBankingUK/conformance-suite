<template>
  <div class="d-flex flex-column h-100 navbar">
    <div class="nav-section">
      <b-nav
        vertical
        pills>
        <NavBarItem
          :no="1"
          route="/wizard/continue-or-start"
          label="Start/Load Test"/>
        <NavBarItem
          :no="2"
          route="/wizard/discovery-config"
          label="Discovery"
          exact/>
        <NavBarItem
          :no="3"
          route="/wizard/configuration"
          label="Configuration"/>
        <NavBarItem
          :no="4"
          route="/wizard/run-overview"
          label="Run/Overview"/>
        <NavBarItem
          :no="5"
          route="/wizard/summary"
          label="Summary"/>
        <NavBarItem
          :no="6"
          route="/wizard/export"
          label="Export"/>
      </b-nav>
    </div>
    <div class="nav-section">
      <h6
        class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted"
      >
        <span>Tools</span>
        <plus-circle-icon/>
      </h6>
      <b-nav
        vertical
        class="mb-2">
        <b-nav-item
          href="https://bitbucket.org/openbankingteam/conformance-suite/src/develop/README.md"
          target="_blank"
        >
          <file-text-icon class="icon-class"/>Documentation
        </b-nav-item>
        <b-nav-item
          href="https://bitbucket.org/openbankingteam/conformance-suite/issues"
          target="_blank"
        >
          <file-text-icon class="icon-class"/>Bug Tracker
        </b-nav-item>
        <b-nav-item
          href="https://bitbucket.org/openbankingteam/conformance-suite"
          target="_blank">
          <file-text-icon class="icon-class"/>Website
        </b-nav-item>
        <b-nav-item
          disabled
          target="_blank">
          <file-text-icon class="icon-class"/>Integrations
        </b-nav-item>
      </b-nav>
    </div>
  </div>
</template>

<script>
import { PlusCircleIcon, FileTextIcon } from 'vue-feather-icons';
import NavBarItem from './NavBarItem.vue';

export default {
  name: 'NavBar',
  components: {
    NavBarItem,
    PlusCircleIcon,
    FileTextIcon,
  },
  // inheritAttrs: false,
  methods: {},
};
</script>

<style scoped>
.navbar {
  box-shadow: 6px 0 25px 0 rgba(38, 50, 56, 0.2);
  padding: 30px 0 0 0;
  width: 190px;
}
.nav-section {
  background: #ffffff;
}

.sidebar-heading {
  font-size: 0.75rem;
  text-transform: uppercase;
}

.nav-item a {
  color: #9e9e9e;
}

.icon-class {
  margin-right: 4px;
}
</style>
