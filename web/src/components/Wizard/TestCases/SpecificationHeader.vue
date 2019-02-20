<template>
  <div>
    <h6>API Specification</h6>
    <b-table
      :items="[ apiSpecificationWithConsentUrls ]"
      :fields="tableFields"
      small
      fixed
      stacked
    >
      <template
        slot="name"
        slot-scope="data">
        <a
          :href="data.item.url"
          target="_blank">
          {{ data.item.name }}
        </a>
      </template>
      <template
        slot="schema"
        slot-scope="data"
      >
        <a
          :href="data.item.schemaVersion"
          target="_blank">
          {{ data.item.version }}
        </a>
      </template>
      <template
        slot="consentUrls"
        slot-scope="data">
        <template v-for="(url, index) in data.value">
          <a
            :href="url"
            :key="url"
            target="_blank">
            Start PSU Consent
          </a>
          <br :key="index">
        </template>
      </template>
    </b-table>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const {
  mapState,
} = createNamespacedHelpers('testcases');

export default {
  name: 'SpecificationHeader',
  props: {
    // Example value for `apiSpecification`.
    // {
    //     "name": "Account and Transaction API Specification",
    //     "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
    //     "version": "v3.0",
    //     "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
    // }
    apiSpecification: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...mapState([
      'consentUrls',
    ]),
    specConsentUrls() {
      return this.consentUrls[this.apiSpecification.name];
      // Uncomment below and comment line above to test before backend consent URL changes finished:
      // return [`http://example.com/${this.apiSpecification.name}/1`, `http://example.com/${this.apiSpecification.name}/2`];
    },
    apiSpecificationWithConsentUrls() {
      const consentUrls = this.specConsentUrls;
      return Object.assign({ consentUrls }, this.apiSpecification);
    },
    hasConsentUrls() {
      return this.specConsentUrls && this.specConsentUrls.length > 0;
    },
    /**
     * Fields to display in API Specification Table.
     * See documentation: https://bootstrap-vue.js.org/docs/components/table#fields-column-definitions-
     */
    tableFields() {
      const fields = {
        name: {
          tdClass: 'table-data-breakable api-specification-table',
        },
        schema: {
          tdClass: 'table-data-breakable api-specification-table',
        },
      };
      if (this.hasConsentUrls) {
        Object.assign(fields, {
          consentUrls: {
            tdClass: 'table-data-breakable api-specification-table',
          },
        });
      }
      return fields;
    },
  },
};
</script>

<style scoped>
  /*
   * Don't remove the `/deep/` here.
   *
   * This rule ensures values such as 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0'
   * in the table are broken into separate lines.
   */
  .test-case /deep/ .table-data-breakable {
    word-break: break-all;
  }

  .test-case /deep/ .api-specification-table {
    grid-template-columns: 20% auto !important;
  }
</style>
