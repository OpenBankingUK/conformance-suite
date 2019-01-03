<template>
  <div>
    <h5>{{ this.$options.name }}</h5>
    <TestCase
      v-for="(testCase, index) in testCases"
      :key="index"
      :test-case="testCase"
    />
  </div>
</template>

<script>
import TestCase from './TestCase.vue';

export default {
  name: 'TestCases',
  components: {
    TestCase,
  },
  props: {
    testCases: {
      type: Array,
      required: true,
      private: true,
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
