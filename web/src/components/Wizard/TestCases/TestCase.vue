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
      <template slot="table-caption">
        <SpecificationHeader
          :apiSpecification="apiSpecification"
        />
      </template>

      <!-- format status column as Bootstrap badge. -->
      <template
        slot="meta.status"
        slot-scope="row">
        <b-badge
          :variant="row.value === 'PASSED' ? 'success' : (row.value === 'FAILED' ? 'danger' : (row.value === 'PENDING' ? 'info' : 'secondary'))"
          :class="row.value === 'FAILED' ? 'clickable' : ''"
          tag="h6"
          @click.stop="toggleError(row)"
        >{{ row.value }}</b-badge>
      </template>

      <template
        slot="row-details"
        slot-scope="row">
        <b-card>
          <strong>Test ID:</strong> {{ row.item.id }}<br>
          <strong>Error:</strong> {{ row.item.error }}<br>
        </b-card>
      </template>
    </b-table>
  </div>
</template>

<script>
import SpecificationHeader from './SpecificationHeader'

export default {
  name: 'TestCase',
  components: {
    SpecificationHeader,
  },
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
          tdClass: 'response-time',
          label: 'Response Time',
        },
        'meta.metrics.responseSize': {
          tdClass: 'response-size',
          label: 'Response Bytes',
        },
      }),
    },
  },
  computed: {
    apiSpecification() {
      return this.testCase.apiSpecification;
    }
  },
  methods: {
    toggleError(row) {
      if (row.item.error) {
        this.$store.commit('testcases/TOGGLE_ROW_DETAILS', row.item);
      }
    },
  },
};
</script>

<style scoped>
/**
 * Don't remove the `/deep/` here.
 *
 * This rule ensures values such as 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0'
 * in the table are broken into separate lines.
 */
.test-case /deep/ .table-data-breakable {
  word-break: break-all;
}

.clickable {
  cursor: pointer;
}

.test-case /deep/ .response-time,
.test-case /deep/ .response-size {
  text-align: center;
}
</style>
