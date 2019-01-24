<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Overview</h5>
        </div>
        <div class="flex-fill panel-body">
          <TheErrorStatus />
          <TestCases
            :test-cases="testCases"
            if="!hasErrors"/>
          <hr>
          <TestCaseResults
            v-if="hasTestCaseResults"
            :test-case-results="execution"/>
          <TheErrorStatus />
        </div>
        <TheWizardFooter :next-label="computeNextLabel"/>
      </div>
    </div>
  </div>
</template>

<script>
import * as _ from 'lodash';

import { createNamespacedHelpers, mapGetters, mapActions } from 'vuex';

import TheWizardFooter from '@/components/Wizard/TheWizardFooter.vue';
import TestCases from '@/components/Wizard/TestCases/TestCases.vue';
import TestCaseResults from '@/components/Wizard/TestCaseResults/TestCaseResults.vue';
import TheErrorStatus from '@/components/TheErrorStatus.vue';

const {
  mapState,
} = createNamespacedHelpers('testcases');

export default {
  name: 'WizardRunOverview',
  components: {
    TheWizardFooter,
    TestCases,
    TestCaseResults,
    TheErrorStatus,
  },
  computed: {
    ...mapGetters('testcases', [
      'testCases',
    ]),
    ...mapGetters('status', [
      'hasErrors',
    ]),
    ...mapState({
      execution: 'execution',
      hasRunStarted: 'hasRunStarted',
    }),
    hasTestCaseResults() {
      return !_.isEmpty(this.execution);
    },
    computeNextLabel() {
      if (!this.hasRunStarted || !this.hasTestCaseResults) {
        return 'Run';
      }

      return 'Next Export';
    },
  },
  methods: {
    ...mapActions('testcases', [
      'computeTestCases',
      'executeTestCases',
    ]),
    ...mapActions('status', [
      'clearErrors',
    ]),
  },
  /**
   * Fetch all the test cases when we navigate to this route.
   * Docs: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
   */
  beforeRouteEnter(to, from, next) {
    next(async (vm) => {
      await vm.computeTestCases();
    });
  },
  /**
   * Prevent user from going forward if there is an error with test case generation, and
   * execute the test cases if the route being navigated to is `/wizard/export`.
   * Docs: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
   */
  async beforeRouteLeave(to, from, next) {
    const isNext = to.path === '/wizard/export' && from.path === '/wizard/overview-run';

    if (isNext) {
      // prevent going forward if there is an error
      if (this.hasErrors) {
        return next(false);
      }

      // Execute and compute results once.
      if (!this.hasRunStarted) {
        await this.executeTestCases();

        return next(false);
      }

      // If there are no results prevent navigation.
      if (!this.hasTestCaseResults) {
        return next(false);
      }

      // We have executed and computed the results and we have results.
      return next();
    }

    // Clear errors before going to a prior step
    this.clearErrors();
    return next();
  },
};
</script>

<style scoped>
</style>
