<template>
  <div class="test-case border p-2 mt-2">
    <b-table
      :items="testCase.testCases"
      :fields="tableFields"
      head-variant="dark"
      caption-top
      hover
      small
      responsive
    >
      <!-- format status column as Bootstrap badge. -->
      <template
        slot="meta.status"
        slot-scope="data">
        <b-badge
          :variant="data.value === 'PASSED' ? 'success' : (data.value === 'FAILED' ? 'danger' : (data.value === 'PENDING' ? 'info' : 'secondary'))"
          tag="h6"
        >{{ data.value }}</b-badge>
      </template>
      <template
        slot="show_details"
        slot-scope="row">
        <b-button
          v-if="row.item.error"
          v-model="row.detailsShowing"
          @click.stop="row.toggleDetails">
          <i
            v-if="row.item._showDetails"
            class="fa fa-minus-square">-</i>
          <i
            v-else
            class="fa fa-plus-square">+</i>
        </b-button>
      </template>

      <template
        slot="row-details"
        slot-scope="row">
        <b-card>
          The following error occurred during test execution:<br>
          {{ row.item.error }}
        </b-card>
      </template>

      <template slot="table-caption">
        <h6>API Specification</h6>
        <b-table
          :items="[ testCase.apiSpecification ]"
          :fields="apiSpecificationTableFields"
          small
          fixed
          stacked
        >
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
        </b-table>
      </template>
    </b-table>
  </div>
</template>

<script>


export default {
  name: 'TestCase',
  components: {},
  props: {
    // Example value for `testCase`.
    // {
    //   "apiSpecification": {
    //     "name": "Account and Transaction API Specification",
    //     "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
    //     "version": "v3.0",
    //     "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
    //   },
    //   "testCases": [
    //     {
    //       "@id": "#t1000",
    //       "name": "Create Account Access Consents",
    //       "input": {
    //         "method": "POST",
    //         "endpoint": "/account-access-consents",
    //         "contextGet": {}
    //       },
    //       "expect": {
    //         "status-code": 201,
    //         "schema-validation": true,
    //         "contextPut": {}
    //       }
    //     }
    //   ]
    // }
    testCase: {
      type: Object,
      required: true,
    },
    /**
             * Fields to display in the table.
             * See documentation: https://bootstrap-vue.js.org/docs/components/table#fields-column-definitions-
             */
    tableFields: {
      type: Object,
      default: () => ({
        show_details: {
          label: '',
          tdClass: 'table-data-breakable',
        },
        '@id': {},
        name: {
          tdClass: 'table-data-breakable',
        },
        'input.method': {
          tdClass: 'table-data-breakable',
        },
        'input.endpoint': {
          tdClass: 'table-data-breakable',
        },
        'expect.status-code': {
          tdClass: 'table-data-breakable',
        },
        'meta.status': {
          label: 'Status',
        },
        'meta.metrics.responseTime': {
          label: 'Response Time',
        },
        'meta.metrics.responseSize': {
          label: 'Response Body Size',
        },
      }),
    },
    /**
             * Fields to display in API Specification Table.
             * See documentation: https://bootstrap-vue.js.org/docs/components/table#fields-column-definitions-
             */
    apiSpecificationTableFields: {
      type: Object,
      default: () => ({
        name: {
          tdClass: 'table-data-breakable api-specification-table',
        },
        url: {
          tdClass: 'table-data-breakable api-specification-table',
        },
        version: {
          tdClass: 'table-data-breakable api-specification-table',
        },
        schemaVersion: {
          tdClass: 'table-data-breakable api-specification-table',
        },
      }),
    },
  },
  data() {
    return {};
  },
  computed: {},
  methods: {},
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
