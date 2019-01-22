<template>
  <div class="test-case border p-2 mt-2">
    <b-table
      :items="testResult.tests"
      :fields="fields"
      head-variant="dark"
      caption-top
      hover
      small
      responsive
    >
      <template slot="table-caption">
        <h6>API Specification</h6>
        <b-table
          :items="[ testResult ]"
          :fields="fieldsSpecification"
          small
          fixed
          stacked>
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
  name: 'TestCaseResult',
  components: {},
  props: {
    // {
    //   name: 'spec name',
    //   version: 'spec version',
    //   url: 'url',
    //   schemaVersion: 'spec schema version',
    //   pass: true,
    //   tests: [
    //     {
    //       name: 'test name',
    //       id: 'test id',
    //       endpoint: 'test endpoint',
    //       pass: true,
    //     },
    //   ],
    // }
    testResult: {
      type: Object,
      required: true,
    },
    /**
     * Fields to display in the table.
     * See documentation: https://bootstrap-vue.js.org/docs/components/table#fields-column-definitions-
     */
    fields: {
      type: Object,
      default: () => ({
        id: {
          tdClass: 'table-data-breakable',
        },
        name: {
          tdClass: 'table-data-breakable',
        },
        pass: {
          tdClass: 'table-data-breakable',
        },
      }),
    },
    /**
     * Fields to display in API Specification Table.
     * See documentation: https://bootstrap-vue.js.org/docs/components/table#fields-column-definitions-
     */
    fieldsSpecification: {
      type: Object,
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
/**
 * Don't remove the `/deep/` here.
 * This rule ensures values such as 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0'
 * in the table are broken into separate lines.
 */
.test-case /deep/ .table-data-breakable {
  word-break: break-all;
}
</style>
