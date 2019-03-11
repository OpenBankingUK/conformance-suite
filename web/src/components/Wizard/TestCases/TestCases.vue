<template>
  <div>
    <h5>{{ this.$options.name }}</h5>
    <TestCase
      v-for="(testGroup, index) in testCases"
      :key="index"
      :test-group="testGroup"
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
