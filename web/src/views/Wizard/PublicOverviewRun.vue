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
                    <b>Token Acquisition mode:</b> {{ tokenAcquisition }}
                  </div>
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
              :test-cases="publicTestCases"
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
    name: 'PublicOverviewRun',
    components: {
      TheErrorStatus,
      TestCases,
      TheWizardFooter,
    },
    data() {
      return {
        tokenTableFields: [
          { key: 'type', label: 'Type' },
          { key: 'value.token_name', label: 'Token Name' },
        ],
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
        'computePublicTestCases',
      ]),
    },
    /**
     * Fetch all the test cases when we navigate to this route.
     * Docs: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
     */
    beforeRouteEnter(to, from, next) {
      next(async (vm) => {
        await vm.computePublicTestCases();
      });
    },
  };
  </script>
  
  <style scoped>
  </style>
  