<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Overview</h5>
        </div>
        <div class="panel-body">
          <div v-if="hasErrors">
            <h2 class="pt-3 pb-2 mb-3">Errors</h2>
            <b-alert
              v-for="(err, index) in errors"
              :key="index"
              show
              variant="danger">{{ err }}</b-alert>
          </div>

          <TestCases
            v-else-if="!hasErrors"
            :test-cases="testCases"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import * as _ from 'lodash';
import TestCases from '../TestCases/TestCases.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('config');

export default {
  name: 'RunOverview',
  components: {
    TestCases,
  },
  data() {
    return {};
  },
  computed: {
    ...mapGetters([
      'testCases',
    ]),
    errors() {
      return this.$store.getters['config/errors'].testCases;
    },
    hasErrors() {
      return this.errors && this.errors.length > 0;
    },
  },
  methods: {
    ...mapActions([
      'computeTestCases',
    ]),
    async onCompute() {
      await this.computeTestCases();
    },
  },
  /**
   * Fetch all the test cases when we navigate to this route.
   *
   * Docs: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
   */
  beforeRouteEnter(to, from, next) {
    next(async (vm) => {
      await vm.computeTestCases();
    });
  },
  /**
   * Prevent user from going forward if there is an error with test case generation.
   *
   * Docs: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
   */
  async beforeRouteLeave(to, from, next) {
    const nextRoutes = [
      '/wizard/summary',
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
.wizard-step {
  min-height: 900px;
}
</style>
