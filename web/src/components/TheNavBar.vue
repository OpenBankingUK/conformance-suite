<template>
  <div class="d-flex flex-column align-items-stretch h-500 navbar">
    <div class="nav-section px-1 mb-5">
      <b-nav
        vertical
        pills>
        <TheNavBarItem
          :no="1"
          route="/wizard/continue-or-start"
          label="Start/Load Test"/>
        <TheNavBarItem
          :no="2"
          route="/wizard/discovery-config"
          label="Discovery"
          exact/>
        <TheNavBarItem
          :no="3"
          route="/wizard/configuration"
          label="Configuration"/>
        <TheNavBarItem
          :no="4"
          route="/wizard/overview-run"
          label="Run/Overview"/>
        <TheNavBarItem
          :no="5"
          route="/wizard/export"
          label="Export"/>
      </b-nav>
    </div>

    <div class="nav-section specifications mb-5">
      <h6
        class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted"
      >
        <span>Specifications</span>
      </h6>
      <b-nav
        vertical
        class="mb-2">
        <b-nav-item
          v-for="(specification, index) in specifications"
          :key="index"
          :href="specification.swaggerUIURL"
          target="_blank"
        >
          <i class="swagger-icon icon-class"/>{{ specification.label }}
        </b-nav-item>
      </b-nav>
    </div>

    <div class="nav-section tools flex-fill">
      <h6
        class="sidebar-heading d-flex justify-content-between align-items-center px-3 mb-1 text-muted"
      >
        <span>Tools</span>
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
        <b-nav-item
          href="https://bitbucket.org/openbankingteam/conformance-suite/src/develop/README.md"
          target="_blank">
          {{ suiteVersion }}
        </b-nav-item>
      </b-nav>
    </div>
  </div>
</template>

<script>
import { PlusCircleIcon, FileTextIcon } from 'vue-feather-icons';
import { mapGetters } from 'vuex';
import TheNavBarItem from './TheNavBarItem.vue';
import Specifications from '../../../pkg/model/testdata/spec-config.golden.json';

export default {
  name: 'TheNavBar',
  components: {
    TheNavBarItem,
    PlusCircleIcon,
    FileTextIcon,
  },
  data() {
    const specifications = Specifications
      .filter(specification => specification.Version === 'v3.1.6')
      .map(specification => ({
        label: `${specification.Name} (${specification.Version})`,
        swaggerUIURL: `/swagger/${specification.Identifier}/${specification.Version}/docs`,
        wikiURL: `${specification.URL.Scheme}://${specification.URL.Host}${specification.URL.Path}`,
        specificationURL: `${specification.SchemaVersion.Scheme}://${specification.SchemaVersion.Host}${specification.SchemaVersion.Path}`,
        _specification: specification,
      }));

    return {
      specifications,
    };
  },
  computed: {
    ...mapGetters('status', ['suiteVersion']),
  },
  methods: {
  },
};
</script>

<style scoped>
.navbar {
  box-shadow: 6px 0 25px 0 rgba(38, 50, 56, 0.2);
  padding: 30px 0 0 0;
  width: 256px;
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

.nav-section.specifications,
.nav-section.tools {
  font-size: 12px;
}

.icon-class {
  height: 16px;
  width: 16px;
}

.swagger-icon {
  background-image: url('~@/assets/images/swagger/favicon-16x16.png');
  display: inline-block;
  vertical-align: middle;
}
</style>
