<template>
  <div class="test-case border p-2 mt-2">
    <b-table
      :items="testCase.testCases"
      :fields="fields"
      head-variant="dark"
      caption-top
      striped
      hover
      small
      fixed
      responsive
    >
      <template slot="table-caption">
        <h6>API Specification</h6>
        <b-table
          :items="[ testCase.apiSpecification ]"
          :fields="fieldsApiSpecification"
          small
          fixed
          stacked
        />
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
      private: true,
    },
    // Fields to display in the table.
    fields: {
      type: Array,
      private: true,
      default: () => [
        '@id',
        'name',
      ],
    },
    // Fields to display in API Specification Table
    fieldsApiSpecification: {
      type: Object,
      private: true,
      default: () => ({
        name: {
          tdClass: 'table-data-breakable',
        },
        url: {
          tdClass: 'table-data-breakable',
        },
        version: {
          tdClass: 'table-data-breakable',
        },
        schemaVersion: {
          tdClass: 'table-data-breakable',
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
.b-table /deep/ .table-data-breakable {
  word-break: break-all;
}
</style>
