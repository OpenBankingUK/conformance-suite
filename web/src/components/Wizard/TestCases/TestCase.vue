<template>
  <div class="test-case border p-2 mt-2">
    <b-table
      :items="testGroup.testCases"
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
      <template
        slot="name"
        slot-scope="row">
        <truncate
          :text="row.value"
          :length="60"
          clamp="..."
          less="Show Less" />
      </template>

      <template
        slot="input.endpoint"
        slot-scope="row">
        <truncate
          :text="row.value"
          :length="30"
          clamp="..."
          less="Show Less" />
      </template>

      <!-- format status column as Bootstrap badge. -->
      <template
        slot="meta.status"
        slot-scope="row">
        <b-badge
          v-if="row.value !== ''"
          :variant="row.value === 'PASSED' ? 'success' : (row.value === 'FAILED' ? 'danger' : (row.value === 'PENDING' ? 'info' : 'secondary'))"
          :class="row.value === 'FAILED' ? 'clickable' : ''"
          :id="statusIdSelector(row)"
          tag="h6"
          @click.stop="toggleError(row)"
        >{{ row.value }} <i
          v-if="row.value === 'FAILED'"
          class="arrow down"/></b-badge>
      </template>

      <template
        slot="row-details"
        slot-scope="row">
        <b-card>
          <b-card-text><strong>Test ID:</strong> {{ row.item.id }}</b-card-text>
          <b-card-text><strong>Ref URI:</strong> <a
            ref="noreferrer"
            :href="row.item.refURI"
            target="_blank"> {{ row.item.refURI }}</a></b-card-text>
          <b-card-text><strong>Detail:</strong> {{ row.item.detail }}</b-card-text>
          <b-card-text><strong>Errors:</strong>
            <ol>
              <li
                v-for="error in row.item.error"
                :key="error">
                <ul>
                  <li><strong>Test Case message:</strong> {{ JSON.parse(error).testCaseMessage }}</li>
                  <li><strong>Endpoint response:</strong> {{ JSON.parse(error).endpointResponse }}</li>
                </ul>
              </li>
            </ol>
          </b-card-text>
        </b-card>
      </template>
    </b-table>
  </div>
</template>

<script>
import truncate from 'vue-truncate-collapsed';
// https://github.com/kavalcante/vue-truncate-collapsed
import SpecificationHeader from './SpecificationHeader.vue';

export default {
  name: 'TestCase',
  components: {
    SpecificationHeader,
    truncate,
  },
  props: {
    // Example value for `testGroup`.
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
    testGroup: {
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
          fixed: true,
        },
        '@id': {},
        name: {
          tdClass: 'table-data-breakable',
        },
        'input.method': {
          tdClass: 'table-data-breakable',
          label: 'Method',
        },
        'input.endpoint': {
          tdClass: 'table-data-breakable',
          label: 'Endpoint',
        },
        'expect.status-code': {
          sortable: true,
          label: 'Expect',
          formatter(value, key, item) {
            if (item.expect['status-code'] > 0) {
              return item.expect['status-code'];
            }

            return item
              .expect_one_of
              .map(expect => expect['status-code'])
              .filter(statusCode => statusCode > 0)
              .join(' or ');
          },
        },
        'meta.status': {
          label: 'Status',
        },
        'meta.metrics.responseTime': {
          tdClass: 'response-time',
          sortable: true,
          label: 'Time',
        },
        'meta.metrics.responseSize': {
          tdClass: 'response-size',
          sortable: true,
          label: 'Bytes',
        },
      }),
    },
  },
  computed: {
    apiSpecification() {
      return this.testGroup.apiSpecification;
    },
  },
  methods: {
    statusIdSelector(row) {
      return row.item['@id'].replace('#', '');
    },
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
