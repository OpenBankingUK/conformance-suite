<template>
  <div class="d-flex flex-row justify-content-between p-3">
    <b-btn
      id="back"
      :disabled="isBackDisabled"
      variant="primary"
      @click="onBack()">Back</b-btn>
    <b-btn
      id="next"
      :disabled="isNextDisabled"
      variant="success"
      @click="onNext()">
      <b-spinner
        v-if="showLoading"
        small/>
      <span class="next-label">{{ nextLabel }}</span>
    </b-btn>
  </div>
</template>

<script>
import invert from 'lodash/invert';
import { createNamespacedHelpers } from 'vuex';
import BSpinner from '../BSpinner';

const { mapGetters } = createNamespacedHelpers('status');

export default {
  name: 'TheWizardFooter',
  components: {
    BSpinner,
  },
  props: {
    nextLabel: {
      type: String,
      required: false,
      default: () => 'Next',
    },
    isNextEnabled: {
      type: Boolean,
      required: false,
      default: () => true,
    },

    nextRoute: {
      type: Function,
      required: false,
      default() {
        const { path } = this.$route;
        return this.ROUTES_NEXT[path];
      },
    },
    previousRoute: {
      type: Function,
      required: false,
      default() {
        const { path } = this.$route;
        return this.ROUTES_BACK[path];
      },
    },

    onBack: {
      type: Function,
      required: false,
      default() {
        this.$router.push(this.previousRoute());
      },
    },
    onNext: {
      type: Function,
      required: false,
      default() {
        this.$router.push(this.nextRoute());
      },
    },
  },
  data() {
    // normal navigation
    const ROUTES = {
      '/wizard/continue-or-start': '/wizard/discovery-config',
      '/wizard/discovery-config': '/wizard/configuration',
      '/wizard/configuration': '/wizard/overview-run',
      '/wizard/overview-run': '/wizard/export',
    };
    // invert the ROUTES map, could do this manually but there is no point
    // as it is error-prone.
    const ROUTES_INVERTED = invert(ROUTES);

    return {
      ROUTES_BACK: ROUTES_INVERTED,
      ROUTES_NEXT: ROUTES,
    };
  },
  computed: {
    ...mapGetters([
      'showLoading',
    ]),
    // disable back button if we are in the first step of the wizard
    isBackDisabled() {
      const { path } = this.$route;
      if (path === '/wizard/continue-or-start') {
        return true;
      }

      return false;
    },
    // disable next button when showLoading is true.
    // disable next button when we are on discovery template selection step.
    // disable next button when property `isNextEnabled` is false.
    isNextDisabled() {
      const { path } = this.$route;
      if (this.showLoading) {
        return true;
      }
      if (path === '/wizard/continue-or-start') {
        return true;
      }
      if (!this.isNextEnabled) {
        return true;
      }

      return false;
    },
  },
  methods: {
  },
};
</script>

<style scoped>
.next-label {
  margin-left: 7px;
  margin-right: 7px;
}
</style>
