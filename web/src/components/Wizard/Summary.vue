<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Summary</h5>
        </div>
        <div class="flex-fill panel-body">
          <p>body</p>
        </div>
        <WizardFooter/>
      </div>
    </div>
  </div>
</template>

<script>
import WizardFooter from './WizardFooter.vue';

export default {
  name: 'Summary',
  components: {
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

<style scoped>
</style>
