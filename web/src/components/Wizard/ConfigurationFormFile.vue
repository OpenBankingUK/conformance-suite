<template>
  <b-form-group
    :id="groupId"
    :description="description"
    :label="label"
    :label-for="id"
  >
    <b-form-file
      :id="id"
      :ref="id"
      v-model="file"
      :state="isValid"
      :placeholder="placeholder"
      :accept="validExtension"
      capture
      @input="() => { onFileChanged() }"
    />
  </b-form-group>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';

export default {
  name: 'ConfigurationFormFile',
  props: {
    id: {
      type: String,
      required: true,
    },
    setterMethodNameSuffix: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
    placeholder: {
      type: String,
      required: false,
      default: 'Choose a file...',
    },
    validExtension: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      file: null,
      data: '',
      validFile: false,
    };
  },
  computed: {
    ...mapGetters('config', [
      'configuration',
    ]),
    groupId() {
      return `${this.id}_group`;
    },
    isValid() {
      const contents = this.configuration[this.id];
      if (contents !== this.data) {
        // Clear file, as JSON editor has changed contents
        this.clearFile();
      }

      if (this.file) {
        return this.file.name.endsWith(this.validExtension);
      }
      return Boolean(contents);
    },
    /**
     * Description of the file uploaded (when one is selected).
     * Returns the size and last modification date.
     * Else returns contents length from vuex store.
     */
    description() {
      const contents = this.configuration[this.id];
      if (this.file && (contents === '' || contents === this.data)) {
        if (this.validFile) {
          return this.fileInfo();
        }
        return this.extError(this.validExtension);
      }
      if (contents) {
        return `Size: ${contents.length} bytes`;
      }

      return '';
    },
  },
  methods: {
    ...mapActions('status', [
      'clearErrors',
      'setErrors',
    ]),
    extError(extension) {
      return ['Invalid file format', `Require file with extension ${extension}`].join('\n');
    },
    fileInfo() {
      const size = `Size: ${this.file.size} bytes`;
      if (this.file.lastModifiedDate) {
        return [
          size,
          `Last modified: ${this.file.lastModifiedDate}`,
        ].join('\n');
      }
      return size;
    },
    clearFile() {
      if (this.$refs[this.id] && this.$refs[this.id].reset) {
        this.$refs[this.id].reset();
      }
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
    async onFileChanged() {
      // Clear previous error.
      this.clearErrors();
      // Compute the method name we need to call in the Vuex store, e.g., could be one of the below:
      // * config/setConfigurationSigningPrivate
      // * config/setConfigurationSigningPublic
      const setConfigurationMethodName = `config/setConfiguration${this.setterMethodNameSuffix}`;

      if (this.file) {
        this.validFile = this.file.name.endsWith(this.validExtension);

        // If file is set, read the file then set the value in the store.
        try {
          this.data = await this.readFile(this.file);
          this.$store.dispatch(setConfigurationMethodName, this.data);
        } catch (err) {
          this.setErrors([err.message]);
        }
      } else {
        // If no file selected assume they want to clear out the previous file.
        this.data = '';
        this.$store.dispatch(setConfigurationMethodName, this.data);
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
