<template>
  <div>
    <div class="d-flex">
      <div class="flex">
        <NavBar/>
      </div>
      <div class="flex-fill">
        <router-view class="d-flex flex-column flex-fill p-3 router-section"/>
        <WizardFooter/>
      </div>
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
import ContinueOrStart from './Wizard/ContinueOrStart.vue';
import DiscoveryConfig from './Wizard/DiscoveryConfig.vue';
import Configuration from './Wizard/Configuration.vue';
import RunOverview from './Wizard/RunOverview.vue';
import Summary from './Wizard/Summary.vue';
import Export from './Wizard/Export.vue';
import WizardFooter from './Wizard/WizardFooter.vue';

export default {
  name: 'Wizard',
  components: {
    ContinueOrStart,
    DiscoveryConfig,
    Configuration,
    RunOverview,
    Summary,
    Export,
    NavBar,
    WizardFooter,
  },
  data() {
    return {};
  },
  computed: {},
  methods: {},
};
</script>
