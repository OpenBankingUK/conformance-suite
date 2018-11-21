<template>
  <div class="d-flex flex-column">
    <h3>Start a new Test or Load Previous</h3>
    <div class="d-flex align-items-start flex-column">
      <p>Start a new test of v3.0 of the Open Banking Specifications</p>
      <b-button>New v3.0 Test</b-button>
    </div>
    <div class="d-flex align-items-start flex-column">
      <p>Import results from a v3.0 test</p>
      <b-form-file
        v-model="file"
        :state="Boolean(file)"/>
    </div>
  </div>
</template>

<style>
.custom-file-input:lang(en)~.custom-file-label::after {
    content: "Upload";
}
</style>

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
      console.warn('%s.validate=%s', this.$options.name, valid);
      // TODO: remove stub.
      return Promise.resolve(true);
    },
  },
};
</script>
