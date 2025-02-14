<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Overview</h5>
        </div>
        <div class="flex-fill panel-body">
          <TheErrorStatus/>
          <PublicTestCases
            :test-cases="publicTestCases"
            if="!hasErrors"/>
          <hr>
          <TheErrorStatus/>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { createNamespacedHelpers, mapGetters, mapActions } from 'vuex';

import TheErrorStatus from '../../components/TheErrorStatus.vue';
import PublicTestCases from '../../components/Wizard/TestCases/PublicTestCases.vue';
import TheWizardFooter from '../../components/Wizard/TheWizardFooter.vue';

const {
  mapState,
} = createNamespacedHelpers('testcases');

export default {
  name: 'PublicOverviewRun',
  components: {
    TheErrorStatus,
    PublicTestCases,
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
      'publicTestCases',
    ]),
    headlessConsent() {
      return this.tokenAcquisition === 'headless';
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
      try {
        await vm.computePublicTestCases();
      } catch (err) {
        console.error('Failed to load test cases:', err);
      }
    });
  },
};
</script>

  <style scoped>
  </style>
