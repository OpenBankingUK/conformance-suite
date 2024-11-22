<template>
    <div class="d-flex flex-row flex-fill">
      <div class="d-flex align-items-start">
        <div class="d-flex flex-column panel w-100 wizard-step">
          <div class="panel-heading">
            <h5>{{ componentHeading }}</h5>
          </div>
          <div class="panel-body">
            <table class="table table-striped">
                <thead>
                    <tr>
                        <th scope="col">Test ID</th>
                        <th scope="col">Created Date</th>
                        <th scope="col">API Version</th>
                        <th scope="col"></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="test in tests" :key="test.id">
                        <td>{{ test.ID }}</td>
                        <td>{{ test.CreatedAt }}</td>
                        <td>{{ test.DiscoveryModel.discoveryModel.discoveryItems[0].apiSpecification.version }}</td>
                        <td>
                            <b-btn variant="success" @click="handleRerun(test.DiscoveryModel.discoveryModel)">
                                Re-Run
                            </b-btn>
                        </td>
                    </tr>
                </tbody>
            </table>
          </div>
          <div
            v-if="error"
            class="panel-body text-danger">
            Error: {{ error }}
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import isEmpty from 'lodash/isEmpty';
  import axios from 'axios';
  import { mapActions } from 'vuex';
  
  export const MODES = {
    REVIEW: 'REVIEW',
    RERUN: 'RERUN',
  };
  
  export default {
    name: 'WizardHistorical',
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
        error: null,
        tests: [],
      };
    },
    async created() {
        try {
            const response = await axios.get('/api/runs/historical');
            this.tests = response.data;
        } catch (err) {
            this.error = 'Failed to fetch test data';
        }
    },
    computed: {
      componentHeading() {
        return 'Historical Test Runs';
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
      ...mapActions('config', ['setDiscoveryModel']),
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
          reader.onloadend = () => resolve(reader.result);
          reader.readAsDataURL(file);
        });
      },
      async handleRerun(discoveryModel) {
        try {
            this.setDiscoveryModel(JSON.stringify({ discoveryModel: discoveryModel }));
            this.$router.push('/wizard/discovery-config');
        } catch (err) {
            this.error = err.error || 'An error occurred during import';
        }
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
        try {
          const results = await this.doImport();
          this.setDiscoveryModel(JSON.stringify({ discoveryModel: results.discoveryModel }));
          this.$router.push('/wizard/discovery-config');
        } catch (err) {
          this.error = err.error || 'An error occurred during import';
        }
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
          const err = new Error(`WizardHistorical: invalid mode=${vm.mode}`);
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