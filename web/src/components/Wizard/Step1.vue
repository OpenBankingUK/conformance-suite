<template>
  <div class="d-flex flex-column align-items-start">
    <h3 class="mb-4">Start a new Test or Load Previous</h3>

    <p>Start a new test of v3.0 of the Open Banking Specifications</p>
    <b-button class="mb-4">New v3.0 Test</b-button>

    <p>Import results from a v3.0 test</p>
    <b-form-file
      v-model="file"
      :state="Boolean(file)"/>
  </div>
</template>

<script>
export default {
  name: 'Step1',
  data() {
    return {
      file: null,
    };
  },
  created() {
    // console.log(this.$options);
  },
  methods: {
    // Gets called by top-level Wizard component in the validateStep function.
    async validate() {
      const valid = true;
      // TODO: remove stub.
      return Promise.resolve(valid);
    },
  },
};
</script>

<style>
.custom-file-input:lang(en)~.custom-file-label::after {
    content: "Upload";
}
</style>
