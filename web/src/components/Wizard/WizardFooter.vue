<template>
  <div class="wizard-footer-section d-flex flex-row justify-content-between p-3">
    <b-btn :disabled="isBackDisabled" variant="primary" @click="onBack()">Back</b-btn>
    <b-btn :disabled="isNextDisabled" variant="success" @click="onNext()">Next</b-btn>
  </div>
</template>

<style scoped>
.wizard-footer-section {
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.125);
}
</style>

<script>
import * as _ from "lodash";

export default {
  name: "WizardFooter",
  components: {},
  data() {
    // normal navigation
    const ROUTES = {
      "/wizard/continue-or-start": "/wizard/discovery-config",
      "/wizard/discovery-config": "/wizard/configuration",
      "/wizard/configuration": "/wizard/run-overview",
      "/wizard/run-overview": "/wizard/summary",
      "/wizard/summary": "/wizard/export"
    };
    // invert the ROUTES map, could do this manually but there is no point
    // as it is error-prone.
    const ROUTES_INVERTED = _.invert(ROUTES);

    return {
      ROUTES_BACK: ROUTES_INVERTED,
      ROUTES_NEXT: ROUTES
    };
  },
  computed: {
    // disable back button if we are in the first step of the wizard
    isBackDisabled() {
      const { path } = this.$route;
      if (path === "/wizard/continue-or-start") {
        return true;
      }

      return false;
    },
    // disable next button when we are the last step in the wizard
    // disable next button when we are on discovery template selection step
    isNextDisabled() {
      const { path } = this.$route;
      if (path === "/wizard/export" || path === "/wizard/continue-or-start") {
        return true;
      }

      return false;
    }
  },
  methods: {
    nextRoute() {
      const { path } = this.$route;
      return this.ROUTES_NEXT[path];
    },
    previousRoute() {
      const { path } = this.$route;
      return this.ROUTES_BACK[path];
    },
    onBack() {
      this.$router.push(this.previousRoute());
    },
    onNext() {
      this.$router.push(this.nextRoute());
    }
  }
};
</script>
