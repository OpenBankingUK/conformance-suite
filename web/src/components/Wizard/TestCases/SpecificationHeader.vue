<template>
  <div>
    <h6>API Specification</h6>
    <b-table
      :items="[ apiSpecificationWithConsentUrl ]"
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
        slot-scope="data">
        <a
          :href="data.item.schemaVersion"
          target="_blank">
          {{ data.item.version }}
        </a>
      </template>
      <!-- format url column as anchor. -->
      <template
        slot="url"
        slot-scope="data">
        <a
          :href="data.value"
          target="_blank">{{ data.value }}</a>
      </template>
      <!-- format schemaVersion column as anchor. -->
      <template
        slot="schemaVersion"
        slot-scope="data">
        <a
          :href="data.value"
          target="_blank">{{ data.value }}</a>
      </template>
      <template
        slot="consentUrl"
        slot-scope="data">
        <a
          :href="data.value"
          target="_blank">{{ data.value }}</a>
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
    /**
     * Fields to display in API Specification Table.
     * See documentation: https://bootstrap-vue.js.org/docs/components/table#fields-column-definitions-
     */
    tableFields: {
      type: Object,
      default: () => ({
        name: {
          tdClass: 'table-data-breakable api-specification-table',
        },
        schema: {
          tdClass: 'table-data-breakable api-specification-table',
        },
        consentUrl: {
          tdClass: 'table-data-breakable api-specification-table',
        },
      }),
    },
  },
  computed: {
    ...mapState([
      'consentUrls',
    ]),
    apiSpecificationWithConsentUrl() {
      const consentUrl = this.consentUrls[this.apiSpecification.name];
      return Object.assign(this.apiSpecification, { consentUrl });
    }
  }
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
