<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Overview</h5>
        </div>
        <div class="flex-fill panel-body">
          <TheErrorStatus/>
          <div
            v-if="!headlessConsent"
            class="test-case border p-2 mt-2 mb-2">
            <span
              v-if="wsConnected"
              id="ws-connected"
            />
            <h5>Tokens</h5>
            <b-table
              :items="tokens_acquired"
              :fields="tokenTableFields"
              head-variant="dark"
              caption-top
              hover
              small
              responsive
            >
              <template slot="table-caption">
                <div>
                  <b>Test Cases Completed:</b> {{ test_cases_completed }}
                </div>
                <div>
                  <b>All Token Acquired:</b> {{ tokens_all_acquired }}
                </div>
              </template>
            </b-table>
          </div>
          <TestCases
            :test-cases="testCases"
            if="!hasErrors"/>
          <hr>
          <TheErrorStatus/>
        </div>
        <TheWizardFooter :next-label="computeNextLabel"/>
      </div>
    </div>
  </div>
</template>

<script>
import { createNamespacedHelpers, mapGetters, mapActions } from 'vuex';

import TheErrorStatus from '../../components/TheErrorStatus.vue';
import TestCases from '../../components/Wizard/TestCases/TestCases.vue';
import TheWizardFooter from '../../components/Wizard/TheWizardFooter.vue';

const {
  mapState,
} = createNamespacedHelpers('testcases');

export default {
  name: 'WizardRunOverview',
  components: {
    TheErrorStatus,
    TestCases,
    TheWizardFooter,
  },
  data() {
    return {
      tokenTableFields: {
        type: {
          label: 'Type',
        },
        'value.token_name': {
          label: 'Token Name',
        },
      },
    };
  },
  computed: {
    ...mapGetters('status', [
      'hasErrors',
      'showLoading',
    ]),
    ...mapGetters('config', [
      'tokenAcquisition',
    ]),
    ...mapGetters('testcases', [
      'wsConnected',
    ]),
    ...mapState([
      'consentUrls',
      'testCases',
      'hasRunStarted',
    ]),
    headlessConsent() {
      return this.tokenAcquisition === 'headless';
    },
    pendingPsuConsent() {
      if (this.headlessConsent) {
        return false;
      }
      return !this.tokens_all_acquired;
    },
    computeNextLabel() {
      if (!this.hasRunStarted || !this.test_cases_completed) {
        if (this.pendingPsuConsent) {
          return 'Pending PSU Consent';
        }
        if (this.showLoading) {
          this.setShowLoading(false);
          return 'Pending';
        }
        return 'Run';
      }

      return 'Next Export';
    },
    // Example Value:
    // [
    //     {
    //         "type": "ResultType_AcquiredAccessToken",
    //         "value": {
    //             "token_name": "to1001"
    //         }
    //     }
    // ]
    tokens_acquired: {
      get() {
        return this.$store.state.testcases.tokens.acquired;
      },
    },
    tokens_all_acquired: {
      get() {
        return this.$store.state.testcases.tokens.all_acquired;
      },
    },
    test_cases_completed: {
      get() {
        return this.$store.state.testcases.test_cases_completed;
      },
    },
  },
  methods: {
    ...mapActions('testcases', [
      'computeTestCases',
      'executeTestCases',
    ]),
    ...mapActions('status', [
      'clearErrors',
      'setShowLoading',
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
        await this.executeTestCases(this.setShowLoading);

        return next(false);
      }

      // If tests have not completed, prevent navigation.
      if (!this.test_cases_completed) {
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
