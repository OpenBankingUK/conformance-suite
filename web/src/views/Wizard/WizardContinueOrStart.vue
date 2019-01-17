<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div
        class="panel w-50"
        style="height:300px">
        <div class="panel-heading">
          <h5>Functional Conformance Suite v1</h5>
        </div>
        <div class="panel-body">
          <p>The Functional Conformance Suite is an Open Source test tool provided by Open Banking. The goal of the suite is to provide an easy and comprehensive tool that enables implementers to test interfaces and data endpoints against the Functional API standard.</p>
          <p>This <strong>v0.2.0-alpha</strong> release introduces an example Discovery Template to demonstrate the headless flow for Ozone Model Bank for v3.0 of the OBIE  Accounts and Transactions specifications.</p>
          <p>N.B. This release is not intended to be run in sandbox or production</p>
        </div>
      </div>
      <div
        class="panel w-50"
        style="height:300px">
        <div class="panel-heading">
          <h5>Continue Test</h5>
        </div>
        <div class="panel-body">
          <p>Upload a signed and compatible report to view the results or to rerun.</p>
          <b-button class="mb-4">Unavailable</b-button>
        </div>
      </div>
    </div>

    <div class="d-flex align-items-end">
      <div class="panel w-100">
        <div class="panel-heading">
          <h5>Or choose from a selection of discovery templates to begin.</h5>
        </div>

        <div class="panel-body">
          <b-card-group deck>
            <DiscoveryTemplateCard
              v-for="(template, index) in discoveryTemplates"
              :key="index"
              :discovery-model="template.model.discoveryModel"
              :image="template.image"
            />
          </b-card-group>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import DiscoveryTemplateCard from '../../components/Wizard/DiscoveryTemplateCard.vue';

export default {
  name: 'WizardContinueOrStart',
  components: {
    DiscoveryTemplateCard,
  },
  computed: {
    ...mapGetters('config', ['discoveryTemplates']),
  },
};
</script>

<style scoped>
</style>
