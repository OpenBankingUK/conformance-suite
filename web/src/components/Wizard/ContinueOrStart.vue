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
          <p>To start, upload a compatible report to view the results or select a discovery template to run a new test.</p>
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
          <b-button class="mb-4">Choose file</b-button>
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
              :discoveryModel="template.model.discoveryModel"
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
import DiscoveryTemplateCard from './DiscoveryTemplateCard.vue';

export default {
  name: 'ContinueOrStart',
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
