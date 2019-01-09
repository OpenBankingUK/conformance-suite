<template>
  <div class="test-case border p-2 mt-2">
    <b-table
      :items="testResult.tests"
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
          :items="[ testResult ]"
          :fields="fieldsSpecification"
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
  name: 'TestCaseResult',
  components: {},
  props: {
    //   {
    //     "name": "spec name",
    //     "version": "spec version",
    //     "url": "url",
    //     "schemaVersion": "spec schema version",
    //     "pass": true,
    //     "tests": [
    //       {
    //         "name": "test name",
    //         "id": "test id",
    //         "endpoint": "test endpoint",
    //         "pass": true
    //       }
    //     ]
    //   }
    testResult: {
      type: Object,
      required: true,
      private: true,
    },
    // Fields to display in the table.
    fields: {
      type: Array,
      private: true,
      default: () => [
        'id',
        'name',
        'pass',
      ],
    },
    // Fields to display in API Specification Table
    fieldsSpecification: {
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
        pass: {
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
