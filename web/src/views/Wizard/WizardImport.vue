<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>{{ componentHeading }} <b-badge variant="danger">WIP</b-badge></h5>
        </div>
        <div class="panel-body">
          <b-form @submit="onSubmit">
            <b-form-group
              id="report_group"
              label-for="report"
              description="Report ZIP archive"
              label="report.zip"
            >
              <b-form-file
                id="report"
                v-model="file"
                :state="report_zip_archive_valid"
                placeholder="Choose a file..."
                accept=".zip"
                capture
                @input="() => { onFileChanged() }"
              />
            </b-form-group>
            <b-button
              id="report-submit-btn"
              :disabled="!report_zip_archive_valid"
              type="submit"
              variant="primary"
            >Import</b-button>
          </b-form>
        </div>
        <div class="panel-heading">
          <h5>Import Response</h5>
        </div>
        <div class="panel-body breakable">
          {{ import_response }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import isEmpty from 'lodash/isEmpty';
import { mapActions } from 'vuex';

export const MODES = {
  REVIEW: 'REVIEW',
  RERUN: 'RERUN',
};

export default {
  name: 'WizardImport',
  components: {
  },
  props: {
    mode: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      file: null,
    };
  },
  computed: {
    componentHeading() {
      if (this.is_review) {
        return 'Import Review';
      } if (this.is_rerun) {
        return 'Import Rerun';
      }
      return 'Import';
    },
    report_zip_archive_valid() {
      return this.isNotEmpty(this.report_zip_archive);
    },
    is_review: {
      get() {
        return this.$store.state.importer.is_review;
      },
      set(value) {
        return this.$store.commit('importer/SET_IS_REVIEW', value);
      },
    },
    is_rerun: {
      get() {
        return this.$store.state.importer.is_rerun;
      },
      set(value) {
        return this.$store.commit('importer/SET_IS_RERUN', value);
      },
    },
    report_zip_archive: {
      get() {
        return this.$store.state.importer.report_zip_archive;
      },
      set(value) {
        return this.$store.commit('importer/SET_REPORT_ZIP_ARCHIVE', value);
      },
    },
    import_response: {
      get() {
        return this.$store.state.importer.import_response;
      },
    },
  },
  methods: {
    ...mapActions('importer', [
      'doImport',
    ]),
    isNotEmpty(value) {
      return !isEmpty(value);
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
      if (this.file) {
        // If file is set, read the file then set the value in the store.
        try {
          this.report_zip_archive = await this.readFile(this.file);
        } catch (err) {
          // TODO(mbana): ignoring errors for now just clear out the previously
          // selected file so that they have to re-upload.
          this.report_zip_archive = '';
        }
      } else {
        // If no file selected assume they want to clear out the previous file.
        this.report_zip_archive = '';
      }
    },
    /**
     * When form is submitted.
     */
    async onSubmit(evt) {
      evt.preventDefault();
      await this.doImport();
    },
  },
  /**
   * Set the mode, review or reun, we are running in before we enter route.
   */
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      // Just calling ES6 setters so disable linting rules here.
      /* eslint-disable no-param-reassign */
      if (vm.mode === MODES.REVIEW) {
        vm.is_review = true;
        vm.is_rerun = false;
      } else if (vm.mode === MODES.RERUN) {
        vm.is_review = false;
        vm.is_rerun = true;
      } else {
        const err = new Error(`WizardImport: invalid mode=${vm.mode}`);
        next(err);
      }
      /* eslint-enable no-param-reassign */
    });
  },
};
</script>

<style scoped>
/* Make sure the import response doesn't overflow the screen */
.breakable {
  word-break: break-all;
}
</style>
