<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Configuration</h5>
        </div>
        <div class="flex-fill panel-body">
          <div class="d-flex flex-column flex-fill">
            <b-tabs v-model="tabIndex">
              <b-tab
                id="form-view"
                title="Form view">
                <TheConfigurationForm v-if="tabIndex == 0"/>
              </b-tab>
              <b-tab
                id="json-view"
                title="JSON view">
                <!-- We use v-if to ensure editor renders fresh each toggle to
                render updates. Otherwise updates do not show on editor.-->
                <TheConfigurationJsonEditor v-if="tabIndex == 1"/>
              </b-tab>
            </b-tabs>
            <TheErrorStatus/>
          </div>
        </div>
        <TheWizardFooter/>
      </div>
    </div>
  </div>
</template>

<script>
import * as _ from 'lodash';
import { mapGetters, mapActions } from 'vuex';

import TheConfigurationForm from '../../components/Wizard/TheConfigurationForm.vue';
import TheConfigurationJsonEditor from '../../components/Wizard/TheConfigurationJsonEditor.vue';
import TheErrorStatus from '../../components/TheErrorStatus.vue';
import TheWizardFooter from '../../components/Wizard/TheWizardFooter.vue';

export default {
  name: 'WizardConfigurationTabs',
  components: {
    TheConfigurationForm,
    TheConfigurationJsonEditor,
    TheErrorStatus,
    TheWizardFooter,
  },
  data() {
    return {
      tabIndex: 0,
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
      'validateConfiguration',
    ]),
    ...mapActions('status', [
      'clearErrors',
    ]),
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
      '/wizard/overview-run',
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
    this.clearErrors();
    return next();
  },
};
</script>

<style scoped>
</style>
