<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div
        class="panel w-100"
        style="height: 900px">
        <div class="panel-heading">
          <h5>{{ this.$options.name }}</h5>
        </div>
        <div class="panel-body">
          <div class="d-flex flex-column flex-fill">
            <b-form>
              <!-- Signing certificates -->
              <b-form-group
                id="signing_private_group"
                :description="getFormGroupDescription('signing_private')"
                label="Private Signing Certificate (.key):"
                label-for="signing_private"
              >
                <!--
                maybe limit file selection to these file types:
                * .key: application/x-iwork-keynote-sffkey
                -->
                <b-form-file
                  id="signing_private"
                  v-model="files.signing_private"
                  :state="Boolean(configuration.signing_private)"
                  placeholder="Choose a file..."
                  capture
                  @input="(file) => { onFileChanged(file, 'SigningPrivate') }"
                />
              </b-form-group>
              <b-form-group
                id="signing_public_group"
                :description="getFormGroupDescription('signing_public')"
                label="Public Signing Certificate (.pem):"
                label-for="signing_public"
              >
                <!--
                maybe limit file selection to these file types:
                * .pem: application/x-x509-ca-cert
                -->
                <b-form-file
                  id="signing_public"
                  v-model="files.signing_public"
                  :state="Boolean(configuration.signing_public)"
                  placeholder="Choose a file..."
                  capture
                  @input="(file) => { onFileChanged(file, 'SigningPublic') }"
                />
              </b-form-group>

              <!-- Transport certificates -->
              <b-form-group
                id="transport_private_group"
                :description="getFormGroupDescription('transport_private')"
                label="Private Transport Certificate (.key):"
                label-for="transport_private"
              >
                <!--
                maybe limit file selection to these file types:
                * .key: application/x-iwork-keynote-sffkey
                -->
                <b-form-file
                  id="transport_private"
                  v-model="files.transport_private"
                  :state="Boolean(configuration.transport_private)"
                  placeholder="Choose a file..."
                  capture
                  @input="(file) => { onFileChanged(file, 'TransportPrivate') }"
                />
              </b-form-group>
              <b-form-group
                id="transport_public_group"
                :description="getFormGroupDescription('transport_public')"
                label="Public Transport Certificate (.pem):"
                label-for="transport_public"
              >
                <!--
                maybe limit file selection to these file types:
                * .pem: application/x-x509-ca-cert
                -->
                <b-form-file
                  id="transport_public"
                  v-model="files.transport_public"
                  :state="Boolean(configuration.transport_public)"
                  placeholder="Choose a file..."
                  capture
                  @input="(file) => { onFileChanged(file, 'TransportPublic') }"
                />
                <!-- <div class="mt-3">Selected file: {{ files.transport_public && files.transport_public.name }}</div> -->
              </b-form-group>
            </b-form>
            <div v-if="error || configurationErrors.length > 0">
              <h2 class="pt-3 pb-2 mb-3">Errors</h2>

              <b-alert
                v-if="error"
                show
                variant="danger">{{ error }}</b-alert>
              <b-alert
                v-for="(err, index) in configurationErrors"
                v-else-if="configurationErrors.length > 0"
                :key="index"
                show
                variant="danger"
              >{{ err }}</b-alert>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as _ from 'lodash';
import { mapGetters, mapActions } from 'vuex';

export default {
  name: 'Configuration',
  components: {},
  data() {
    return {
      files: {
        signing_private: null,
        signing_public: null,
        transport_private: null,
        transport_public: null,
      },
      error: null,
    };
  },
  computed: {
    ...mapGetters('config', [
      'configuration',
      'configurationErrors',
    ]),
  },
  methods: {
    ...mapActions('config', [
      'setConfigurationSigningPrivate',
      'setConfigurationSigningPublic',
      'setConfigurationTransportPrivate',
      'setConfigurationTransportPublic',
      'validateConfiguration',
    ]),
    /**
     * Get a description of the file uploaded (when one is selected).
     * Returns the size and last modification date.
     */
    getFormGroupDescription(fileName) {
      const file = this.files[fileName];
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
      this.error = null;

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
          this.error = err.toString();
        }
      } else {
        // If no file selected assume they want to clear out the previous file.
        const data = '';
        setConfigurationMethod(data);
      }
    },
    /**
     * Validates the configuration.
     */
    async validate() {
      const valid = await this.validateConfiguration();
      return valid;
    },
  },
  /**
   * Prevent user from progressing FORWARD only if the Configuration is invalid.
   *
   * "The leave guard is usually used to prevent the user from accidentally leaving the route with unsaved edits. The navigation can be canceled by calling next(false)."
   * See documentation: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
   */
  async beforeRouteLeave(to, from, next) {
    const nextRoutes = [
      '/wizard/run-overview',
      '/wizard/summary',
      '/wizard/export',
    ];
    const isNext = from.path === '/wizard/configuration' && _.includes(nextRoutes, to.path);

    // Allow the user to only go forward if the configuration is valid
    if (isNext) {
      const valid = await this.validate();
      if (valid) {
        return next();
      }

      return next(false);
    }

    return next();
  },
};
</script>

<!-- The `style` cannot have the `scoped` attribute/property set. When `scoped` is set, the Bootstrap component do not have the styles specified applied. -->
<style>
.custom-file-input:lang(en) ~ .custom-file-label::after {
  content: "Upload";
}

/* Ensure line breaks (\n) in the form group description are honoured. */
.b-form-group.form-group .form-text.text-muted {
  white-space: pre-line;
}
</style>

