<template>
  <b-form-group
    :id="`#{id}_group`"
    :description="getFormGroupDescription(id)"
    :label="label"
    :label-for="id"
  >
    <b-form-file
      :id="id"
      v-model="certificate"
      :state="isValid"
      :placeholder="placeholder"
      capture
      @input="(file) => { onFileChanged(file, field) }"
    />
  </b-form-group>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';

export default {
  name: 'Configuration',
  props: {
    id: {
      type: String,
      required: true,
      default: null,
    },
    field: {
      type: String,
      required: true,
      default: null,
    },
    label: {
      type: String,
      required: true,
      default: null,
    },
    placeholder: {
      type: String,
      required: true,
      default: 'Choose a file...',
    },
  },
  data() {
    return {
      certificate: null,
    };
  },
  computed: {
    ...mapGetters('config', [
      'configuration',
    ]),
    isValid() {
      return Boolean(this.configuration[this.id]);
    },
  },
  methods: {
    ...mapActions('config', [
      'setConfigurationErrors',
    ]),
    /**
     * Get a description of the file uploaded (when one is selected).
     * Returns the size and last modification date.
     */
    getFormGroupDescription(fileName) {
      const file = this.certificate;
      const contents = this.configuration[fileName];
      if (file) {
        // File (HTML API) contains these fields:
        // lastModified: 1545301720780
        // lastModifiedDate: Thu Dec 20 2018 10:28:40 GMT+0000 (Greenwich Mean Time) {}
        // name: "transport_private.key"
        // size: 891
        // type: "application/x-iwork-keynote-sffkey"
        // webkitRelativePath: ""
        return [
          `Size: ${file.size} bytes`,
          `Last modified: ${file.lastModifiedDate}`,
        ].join('\n');
      } else if (contents) {
        return [
          `Length: ${contents.length}`,
        ].join('\n');
      }

      return '';
    },
    /**
     * readFile turns FileReader API into a Promise-based one,
     * returning a resolved Promise with the contents of the file
     * when it has been loaded.
     */
    readFile(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = evt => resolve(evt.target.result);
        reader.onerror = evt => reject(new Error(`reading ${file.name}: ${evt.target.result}`));

        reader.readAsText(file);
      });
    },
    /**
     * When a file is selected, read its content and set the value in the store.
     * See: https://stackoverflow.com/questions/45179061/file-input-on-change-in-vue-js
     */
    async onFileChanged(file, setterMethodNameSuffix) {
      // Clear previous error.
      this.setConfigurationErrors([]);

      // Compute the method name we need to call in the Vuex store, e.g., could be one of the below:
      // * setConfigurationSigningPrivate
      // * setConfigurationSigningPublic
      const setConfigurationMethodName = `setConfiguration${setterMethodNameSuffix}`;
      const setConfigurationMethod = this[setConfigurationMethodName];

      if (file) {
        // If file is set, read the file then set the value in the store.
        try {
          const data = await this.readFile(file);
          setConfigurationMethod(data);
        } catch (err) {
          this.setConfigurationErrors([err]);
        }
      } else {
        // If no file selected assume they want to clear out the previous file.
        const data = '';
        setConfigurationMethod(data);
      }
    },
  },
};
</script>

<style scoped>
/* See note on `/deep/` selector: https://vue-loader.vuejs.org/guide/scoped-css.html#deep-selectors */

/* Instead of "Browse" we want "Upload" */
.b-form-group.form-group .custom-file-input:lang(en) /deep/ .custom-file-label::after {
  content: "Upload";
}

/* Ensure line breaks (\n) in the form group description are honoured. */
.b-form-group.form-group /deep/ .form-text.text-muted {
  white-space: pre-line;
}
</style>
