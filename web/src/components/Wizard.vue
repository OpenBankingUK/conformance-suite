<template>
  <div class="d-flex flex-column flex-fill">
    <form-wizard
      :start-index.sync="activeTabIndex"
      title="Conformance Suite"
      subtitle
      color="#6180c3"
      class="p-2 d-flex flex-column flex-fill"
      layout="vertical"
      @on-complete="onComplete"
    >
      <!-- <tab-content
        icon="ti-import"
        title="New/Import">
        <step1/>
      </tab-content>
      <tab-content
        icon="ti-settings"
        title="Configuration">Configuration</tab-content>
      <tab-content
        icon="ti-panel"
        title="Test Overview">Test Overview</tab-content>
      <tab-content
        icon="ti-list"
        title="Summary">Summary</tab-content>
      <tab-content
        icon="ti-export"
        title="Export">Export</tab-content> -->

      <tab-content
        v-for="tab in tabs"
        v-if="!tab.hide"
        :key="tab.title"
        :title="tab.title"
        :icon="tab.icon"
        :before-change="()=>validateStep(tab.component.toLowerCase())">
        <component
          :ref="tab.component.toLowerCase()"
          :is="tab.component"/>
      </tab-content>
    </form-wizard>
  </div>
</template>

<style lang="scss">
.wizard-tab-content {
  width: 100%;
  flex: 1;

}

.wizard-nav,
.wizard-nav-pills,
.wizard-tab-content {
  padding: 0.5rem !important;
  margin: 0.5rem;
  border: 1px solid rgba(0, 0, 0, 0.5);
}

.wizard-navigation {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.wizard-tab-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}
</style>

<script>
import * as _ from 'lodash';
import { FormWizard, TabContent } from 'vue-form-wizard';
import 'vue-form-wizard/dist/vue-form-wizard.min.css';
import 'themify-icons-scss/scss/themify-icons.scss';
import Step1 from './Wizard/Step1.vue';
import Step2 from './Wizard/Step2.vue';
import Step3 from './Wizard/Step3.vue';
import Step4 from './Wizard/Step4.vue';
import Step5 from './Wizard/Step5.vue';

export default {
  name: 'wizard',
  components: {
    FormWizard,
    TabContent,
    Step1,
    Step2,
    Step3,
    Step4,
    Step5,
  },
  data() {
    return {
      activeTabIndex: 1,
      tabs: [
        {
          title: 'New/Import',
          icon: 'ti-import',
          component: 'Step1',
        },
        {
          title: 'Configuration',
          icon: 'ti-settings',
          component: 'Step2',
          hide: false,
        },
        {
          title: 'Test Overview',
          icon: 'ti-panel',
          component: 'Step3',
        },
        {
          title: 'Summary',
          icon: 'ti-list',
          component: 'Step4',
        },
        {
          title: 'Export',
          icon: 'ti-export',
          component: 'Step5',
        },
      ],
    };
  },
  computed: {
  },
  methods: {
    onComplete() {
      // eslint-disable-next-line no-alert
      alert('Yay. Done!');
    },
    async validateStep(name) {
      // this.$refs contains the components Step1-Step5, so we grab the
      // the current component and call the asynchronous function `validate` on it.
      // If it is an Array, we call `validate` on each one and wait for their results
      // to return
      const steps = this.$refs[name];
      if (_.isArray(steps)) {
        const results = await Promise.all(_.map(steps, step => step.validate()));
        // If a single call to `validate` fail
        const result = _.every(results);
        console.info('steps=%o, results=%o, result=%s', steps, results, result);
        return result;
      } else if (_.isObject(steps)) {
        const step = steps[0];
        const result = await step.validate();
        console.info('step=%o, result=%s', step, result);
        return result;
      }

      // do nothing for now
      return false;
    },
  },
};
</script>
