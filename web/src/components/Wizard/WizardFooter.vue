<template>
  <div class="wizard-footer-section d-flex flex-row justify-content-between p-3">
    <b-btn
      :disabled="isBackDisabled"
      variant="primary"
      @click="onBack()">Back</b-btn>
    <b-btn
      :disabled="isNextDisabled"
      variant="success"
      @click="onNext()">Next</b-btn>
  </div>
</template>

<style scoped>
.wizard-footer-section {
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.125);
}
</style>

<script>
import * as _ from 'lodash';

export default {
  name: 'WizardFooter',
  components: {
  },
  data() {
    // normal navigation
    const ROUTES = {
      '/wizard/step1': '/wizard/discovery-config',
      '/wizard/discovery-config': '/wizard/configuration',
      '/wizard/configuration': '/wizard/run-overview',
      '/wizard/run-overview': '/wizard/summary',
      '/wizard/summary': '/wizard/export',
    };
    // invert the ROUTES map, could do this manually but there is no point
    // as it is error-prone.
    //
    // Example output:
    // const ROUTES_BACK = {
    //   '/wizard/discovery-config': '/wizard/step1',
    //   '/wizard/configuration': '/wizard/discovery-config',
    //   '/wizard/run-overview': '/wizard/configuration',
    //   '/wizard/summary': '/wizard/run-overview',
    //   '/wizard/export': '/wizard/summary',
    // };
    const ROUTES_INVERTED = _.invert(ROUTES);

    return {
      ROUTES_BACK: ROUTES_INVERTED,
      ROUTES_NEXT: ROUTES,
    };
  },
  computed: {
    // disable back button if we are in the first step of the wizard
    isBackDisabled() {
      const { path } = this.$route;
      if (path === '/wizard/step1') {
        return true;
      }

      return false;
    },
    // disable next button if we are the last step in the wizard
    isNextDisabled() {
      const { path } = this.$route;
      if (path === '/wizard/export') {
        return true;
      }

      return false;
    },
  },
  methods: {
    onBack() {
      const { ROUTES_BACK: ROUTES } = this;
      const router = this.$router;
      const { path } = this.$route;

      // find new route
      const location = ROUTES[path];
      // go to new route
      router.push(location);
    },
    onNext() {
      const { ROUTES_NEXT: ROUTES } = this;
      const router = this.$router;
      const { path } = this.$route;

      // find new route
      const location = ROUTES[path];
      // go to new route
      router.push(location);
    },
  },
};
</script>
