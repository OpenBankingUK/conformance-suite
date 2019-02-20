<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Export</h5>
        </div>
        <div class="flex-fill panel-body">
          <p>body</p>
        </div>
        <TheWizardFooter/>
      </div>
    </div>
  </div>
</template>

<script>
import TheWizardFooter from '../../components/Wizard/TheWizardFooter.vue';

export default {
  name: 'WizardExport',
  components: {
    TheWizardFooter,
  },
  data() {
    return {};
  },
  computed: {},
  methods: {},
};
</script>

<style scoped>
</style>
