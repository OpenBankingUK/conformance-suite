<template>
  <div>
    <h5>{{ this.$options.name }}</h5>
    <TestCaseResult
      v-for="(testCaseResult, index) in testCaseResults.specifications"
      :key="index"
      :test-result="testCaseResult"
    />
  </div>
</template>

<script>
import TestCaseResult from './TestCaseResult.vue';

export default {
  name: 'TestCaseResults',
  components: {
    TestCaseResult,
  },
  props: {
    // {
    //   id: '00000000-0000-0000-0000-000000000000',
    //   specifications: [
    //     {
    //       name: 'spec name',
    //       version: 'spec version',
    //       url: 'url',
    //       schemaVersion: 'spec schema version',
    //       pass: true,
    //       tests: [
    //         {
    //           name: 'test name',
    //           id: 'test id',
    //           endpoint: 'test endpoint',
    //           pass: true,
    //         },
    //       ],
    //     },
    //   ],
    // }
    testCaseResults: {
      type: Object,
      required: true,
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
</style>
