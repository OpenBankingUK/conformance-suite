<template>
  <div>
    <h5>{{ this.$options.name }}</h5>
    <PublicTestCase
      v-for="(testGroup, index) in testCases"
      :key="index"
      :test-group="testGroup"
    />
  </div>
</template>

<script>
import PublicTestCase from './PublicTestCase.vue';

export default {
  name: 'PublicTestCases',
  components: {
    PublicTestCase,
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
