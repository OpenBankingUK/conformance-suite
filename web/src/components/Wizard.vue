<template>
  <div class="d-flex flex-row flex-fill">
    <NavBar/>
    <div class="d-flex flex-column flex-fill">
      <router-view class="p-3 router-section"/>
      <WizardFooter/>
    </div>
  </div>
</template>

<style scoped>
.router-section {
  background: #f6f6f6;
}
</style>

<script>
import NavBar from './NavBar/NavBar.vue';
import Step1 from './Wizard/Step1.vue';
import DiscoveryConfig from './Wizard/DiscoveryConfig.vue';
import Configuration from './Wizard/Configuration.vue';
import RunOverview from './Wizard/RunOverview.vue';
import Summary from './Wizard/Summary.vue';
import Export from './Wizard/Export.vue';
import WizardFooter from './Wizard/WizardFooter.vue';

export default {
  name: 'Wizard',
  components: {
    Step1,
    DiscoveryConfig,
    Configuration,
    RunOverview,
    Summary,
    Export,
    NavBar,
    WizardFooter,
  },
  data() {
    return {
    };
  },
  computed: {
  },
  methods: {
  },
};
</script>
