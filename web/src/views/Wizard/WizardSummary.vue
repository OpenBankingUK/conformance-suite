<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Summary</h5>
        </div>
        <div class="flex-fill panel-body">
          <div v-if="hasErrors">
            <h2 class="pt-3 pb-2 mb-3">Errors</h2>
            <b-alert
              v-for="(err, index) in errors"
              :key="index"
              show
              variant="danger">{{ err }}
            </b-alert>
          </div>

          <TestCaseResults
            v-else-if="!hasErrors"
            :test-case-results="testCaseResults"/>

          <p>{{ JSON.stringify(execution) }}</p>
        </div>
        <TheWizardFooter/>
      </div>
    </div>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import * as _ from 'lodash';

import TheWizardFooter from '../../components/Wizard/TheWizardFooter.vue';
import TestCaseResults from '../../components/Wizard/TestCaseResults/TestCaseResults.vue';

const { mapGetters } = createNamespacedHelpers('testcases');
const { mapActions, mapState } = createNamespacedHelpers('config');

export default {
  name: 'WizardSummary',
  components: {
    TheWizardFooter,
    TestCaseResults,
  },
  data() {
    return {};
  },
  computed: {
    ...mapGetters(['execution']),
    ...mapState({
      testCaseResults: 'testCaseResults',
    }),
    errors() {
      return this.$store.getters['config/errors'].testCases;
    },
    hasErrors() {
      return this.errors && this.errors.length > 0;
    },
  },
  methods: {
    ...mapActions([
      'computeTestCaseResults',
    ]),
    async onCompute() {
      await this.computeTestCaseResults();
    },
  },

  /**
  * Fetch all the test cases when we navigate to this route.
  *
  * Docs: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
  */
  beforeRouteEnter(to, from, next) {
    next(async (vm) => {
      await vm.computeTestCaseResults();
    });
  },

  /**
  * Prevent user from going forward if there is an error with test case generation.
  *
  * Docs: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
  */
  async beforeRouteLeave(to, from, next) {
    const nextRoutes = [
      '/wizard/export',
    ];
    const isNext = _.includes(nextRoutes, to.path);

    if (isNext) {
      if (this.hasErrors) {
        // prevent going forward if there is an error
        return next(false);
      }

      return next();
    }

    return next();
  },
};
</script>

<style scoped>
</style>
