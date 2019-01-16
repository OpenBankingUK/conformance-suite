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
      :accept="acceptedExtensions"
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
  },
  data() {
    return {
      file: null,
      data: '',
      fileValidation: {
        valid: false,
        extension: '',
      },
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
        switch (this.id) {
          case 'signing_private':
          case 'transport_private':
            return this.file.name.endsWith('.key');
          case 'signing_public':
          case 'transport_public':
            return this.file.name.endsWith('.pem');
          default:
            return false;
        }
      } else {
        return false;
      }
    },
    acceptedExtensions() {
      switch (this.id) {
        case 'signing_private':
        case 'transport_private':
          return '.key';
        case 'signing_public':
        case 'transport_public':
          return '.pem';
        default:
          break;
      }
      return '';
    },
    /**
         * Description of the file uploaded (when one is selected).
         * Returns the size and last modification date.
         * Else returns contents length from vuex store.
         */
    description() {
      const contents = this.configuration[this.id];
      const extError = extension => ['Invalid file format', `Require file with extension ${extension}`].join('\n');
      const fileInfo = f =>
      // File (HTML API) contains these fields:
      // lastModified: 1545301720780
      // lastModifiedDate: Thu Dec 20 2018 10:28:40 GMT+0000 (Greenwich Mean Time) {}
      // name: "transport_private.key"
      // size: 891
      // type: "application/x-iwork-keynote-sffkey"
      // webkitRelativePath: ""
        [
          `Size: ${f.size} bytes`,
          `Last modified: ${f.lastModifiedDate}`,
        ].join('\n');

      if (this.file && (contents === '' || contents === this.data)) {
        if (this.fileValidation.valid) {
          return fileInfo(this.file);
        }
        return extError(this.fileValidation.extension);
      } else if (contents) {
        return `Size: ${contents.length} bytes`;
      }

      return '';
    },
  },
  methods: {
    ...mapActions('config', [
      'setConfigurationErrors',
    ]),
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
      this.setConfigurationErrors([]);
      // Compute the method name we need to call in the Vuex store, e.g., could be one of the below:
      // * config/setConfigurationSigningPrivate
      // * config/setConfigurationSigningPublic
      const setConfigurationMethodName = `config/setConfiguration${this.setterMethodNameSuffix}`;

      if (this.file) {
        switch (this.id) {
          case 'signing_private':
          case 'transport_private':
            if (!this.file.name.endsWith('.key')) {
              this.fileValidation.valid = false;
              this.fileValidation.extension = '.key';
            } else {
              this.fileValidation.valid = true;
            }
            break;
          case 'signing_public':
          case 'transport_public':
            if (!this.file.name.endsWith('.pem')) {
              this.fileValidation.valid = false;
              this.fileValidation.extension = '.pem';
            } else {
              this.fileValidation.valid = true;
            }
            break;
          default:
            this.fileValidation.valid = true;
            break;
        }

        // If file is set, read the file then set the value in the store.
        try {
          this.data = await this.readFile(this.file);
          this.$store.dispatch(setConfigurationMethodName, this.data);
        } catch (err) {
          this.setConfigurationErrors([err.message]);
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
