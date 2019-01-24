<template>
  <div>
    <h5>{{ this.$options.name }}</h5>
    <h5><b>ID</b> {{ testCaseResults.id }}</h5>
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
    //     "id": "008edfbd-497f-4a54-957d-e37b0b774101",
    //     "specifications": [
    //         {
    //             "name": "CustomTest-GetOzoneToken",
    //             "version": "",
    //             "url": "",
    //             "schemaVersion": "",
    //             "pass": false,
    //             "tests": [
    //                 {
    //                     "name": "ClientCredential Grant",
    //                     "id": "#ct0001",
    //                     "endpoint": "POST https://modelobank2018.o3bank.co.uk:4201/token",
    //                     "pass": true
    //                 }
    //             ]
    //         }
    //     ]
    // }
    testCaseResults: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style scoped>
</style>
