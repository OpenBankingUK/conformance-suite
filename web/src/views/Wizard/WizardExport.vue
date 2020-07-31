<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Export</h5>
        </div>
        <div class="flex-fill panel-body">
          <b-card bg-variant="light">
            <h5>List of tested APIs</h5>
            <ul id="versions">
              <li
                v-for="{version, name } in api_versions"
                :key="version + name"
                :id="(version+name).replace(/[^a-zA-Z0-9-]/g, '_')"
              >
                {{ version }} {{ name }}
              </li>
            </ul>
            <h5>Export Configuration</h5>
            <b-form>
              <b-form-group
                id="environment_group"
                label-for="environment_group"
                label="Environment"
              >
                <b-form-select
                  id="environment"
                  v-model="environment"
                  :options="available_environments"
                  :state="isNotEmpty(environment)"
                  required
                />
              </b-form-group>
              <b-form-group
                id="implementer_group"
                label-for="implementer"
                label="Implementer/Brand Name"
              >
                <b-form-input
                  id="implementer"
                  v-model="implementer"
                  :state="isNotEmpty(implementer)"
                  required
                  type="text"
                />
              </b-form-group>
              <b-form-group
                id="authorised_by_group"
                label-for="authorised_by"
                label="Authorised by"
              >
                <b-form-input
                  id="authorised_by"
                  v-model="authorised_by"
                  :state="isNotEmpty(authorised_by)"
                  required
                  type="text"
                />
              </b-form-group>
              <b-form-group
                id="job_title_group"
                label-for="job_title"
                label="Job Title">
                <b-form-input
                  id="job_title"
                  v-model="job_title"
                  :state="isNotEmpty(job_title)"
                  required
                  type="text"
                />
              </b-form-group>
              <b-form-group
                id="products_group"
                label-for="products"
                label="Products">
                <b-form-select
                  id="products"
                  v-model="products"
                  :options="available_products"
                  :state="isNotEmpty(products)"
                  multiple />
              </b-form-group>
              <b-form-group
                id="has_agreed_group"
                :label="has_agreed_terms"
                label-for="has_agreed">
                <b-form-checkbox
                  id="has_agreed"
                  v-model="has_agreed"
                  :required="requires_agreement">
                  I agree
                </b-form-checkbox>
              </b-form-group>
              <b-form-group
                id="add_digital_signature_group"
                label-for="add_digital_signature"
                label="Sign this report with your private key"
              >
                <b-form-checkbox
                  id="add_digital_signature"
                  v-model="add_digital_signature"
                >Add Digital Signature?</b-form-checkbox>
              </b-form-group>
            </b-form>
          </b-card>
          <br >
          <a
            v-if="export_results_blob"
            :href="export_results_download"
            :download="export_results_filename"
            :filename="export_results_filename"
            class="download-report-link"
            target="_blank"
          >
            <b-button
              block
              variant="primary">Download {{ export_results_filename }}</b-button>
          </a>
          <br >
          <b-card
            v-if="hasErrors"
            bg-variant="light">
            <TheErrorStatus />
          </b-card>
        </div>

        <TheWizardFooter
          :is-next-enabled="canExport"
          :next-label="computeFooterNextLabel"
          :on-next="exportResults"
        />
      </div>
    </div>
  </div>
</template>

<script>
import isEmpty from 'lodash/isEmpty';
import every from 'lodash/every';
import { mapGetters, mapActions } from 'vuex';
import TheErrorStatus from '../../components/TheErrorStatus.vue';
import TheWizardFooter from '../../components/Wizard/TheWizardFooter.vue';

export default {
  name: 'WizardExport',
  components: {
    TheErrorStatus,
    TheWizardFooter,
  },
  data() {
    return {
      has_agreed_terms: 'Certification: Implementer has tested the Deployment and verified that it conforms to the OBIE Standard, and hereby certifies to the Open Banking Implementation Entity and the public that the Deployment conforms to the OBIE Standard as set forth above.',
      available_products: [
        'Business',
        'Personal',
        'Cards',
      ],
    };
  },
  computed: {
    ...mapGetters('status', [
      'hasErrors',
    ]),
    environment: {
      get() {
        return this.$store.state.exporter.environment;
      },
      set(value) {
        this.$store.commit('exporter/SET_ENVIRONMENT', value);
      },
    },
    implementer: {
      get() {
        return this.$store.state.exporter.implementer;
      },
      set(value) {
        this.$store.commit('exporter/SET_IMPLEMENTER', value);
      },
    },
    authorised_by: {
      get() {
        return this.$store.state.exporter.authorised_by;
      },
      set(value) {
        this.$store.commit('exporter/SET_AUTHORISED_BY', value);
      },
    },
    api_versions: {
      get() {
        const versions = [];
        /* eslint-disable */
        for (const item of this.$store.state.config.discoveryModel
          .discoveryModel.discoveryItems) {
          versions.push({
            name: item.apiSpecification.name,
            version: item.apiSpecification.version,
          });
        }
        return versions;
        /* eslint-enable */
      },
    },
    job_title: {
      get() {
        return this.$store.state.exporter.job_title;
      },
      set(value) {
        this.$store.commit('exporter/SET_JOB_TITLE', value);
      },
    },
    products: {
      get() {
        return this.$store.state.exporter.products;
      },
      set(value) {
        this.$store.commit('exporter/SET_PRODUCTS', value);
      },
    },
    requires_agreement: {
      get() {
        return this.environment === 'production';
      },
    },
    has_agreed: {
      get() {
        return this.$store.state.exporter.has_agreed;
      },
      set(value) {
        this.$store.commit('exporter/SET_HAS_AGREED', value);
      },
    },
    add_digital_signature: {
      get() {
        return this.$store.state.exporter.add_digital_signature;
      },
      set(value) {
        this.$store.commit('exporter/SET_ADD_DIGITAL_SIGNATURE', value);
      },
    },
    export_results_blob: {
      get() {
        return this.$store.state.exporter.export_results_blob;
      },
    },
    export_results_filename: {
      get() {
        return this.$store.state.exporter.export_results_filename;
      },
    },
    available_environments: {
      get() {
        return [
          'testing',
          'sandbox',
          'production',
        ];
      },
    },
    export_results_download() {
      if (this.export_results_blob) {
        // TODO(mbana): Remember to call `window.URL.revokeObjectURL()`. No big deal.
        return window.URL.createObjectURL(this.export_results_blob);
      }
      return '';
    },
    canExport() {
      const conditions = [
        this.isNotEmpty(this.implementer),
        this.isNotEmpty(this.authorised_by),
        this.isNotEmpty(this.job_title),
        this.has_agreed || !this.requires_agreement,
        this.isNotEmpty(this.products),
      ];
      const canExport = every(conditions);
      return canExport;
    },
    computeFooterNextLabel() {
      return 'Export Conformance Report';
    },
  },
  methods: {
    ...mapActions('exporter', [
      'exportResults',
    ]),
    isNotEmpty(value) {
      return !isEmpty(value);
    },
  },
};
</script>

<style scoped>
</style>
