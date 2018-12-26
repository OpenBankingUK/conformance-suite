<template>
  <div class="align-items-start">
    <h2 class="pt-3 pb-2 mb-3">Continue or Select Disvovery Template</h2>

    <p>Upload a signed and compatible report to view the results or to rerun.</p>
    <b-button class="mb-4">Choose file</b-button>

    <hr class="border-bottom w-100">
    <p>Or start a new test. Choose from a selection of discovery templates to begin.</p>

    <b-card-group deck>
      <DiscoveryTemplateCard
        v-for="(template, index) in discoveryTemplates"
        :key="index"
        :discoveryModel="template.model.discoveryModel"
        :image="template.image"
      />
    </b-card-group>
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
    ...mapGetters('config', [
      'discoveryTemplates',
    ]),
  },
};
</script>

<style>
</style>
