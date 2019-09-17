<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="d-flex flex-column panel w-100 wizard-step">
        <div class="panel-heading">
          <h5>Discovery {{ name }}</h5>
        </div>
        <div class="d-flex flex-column flex-fill panel-body">
          <div>
            When authoring your own discovery file, more information about the Discovery Model can be found
            <a
              href="https://bitbucket.org/openbankingteam/conformance-suite/src/develop/docs/discovery.md"
              target="_blank">here.
            </a>
          </div>
          <br>
          <TheErrorStatus />
          <TheJsonEditor
            :problem-annotations="problemAnnotations"
            :json-string="discoveryModelString"
            editor-name="discovery-config-editor"
            set-change-function-name="config/setDiscoveryModel"
          />
        </div>
        <TheWizardFooter/>
      </div>
    </div>
  </div>
</template>

<script>
import * as _ from 'lodash';
import { mapGetters, mapActions } from 'vuex';

import TheErrorStatus from '../../components/TheErrorStatus.vue';
import TheJsonEditor from '../../components/Wizard/TheJsonEditor.vue';
import TheWizardFooter from '../../components/Wizard/TheWizardFooter.vue';
import discovery from '../../api/discovery';

// Bug in Brace editor using wrong Range function means we need to require Range here:
const AceRange = window.ace.acequire('ace/range').Range;

export default {
  name: 'WizardDiscoveryConfig',
  components: {
    TheErrorStatus,
    TheWizardFooter,
    TheJsonEditor,
  },
  computed: {
    ...mapGetters('config', [
      'discoveryModel',
      'discoveryModelString',
      'problems',
      'discoveryProblems',
    ]),
    /**
     * The name of the discovery template selected.
     */
    name() {
      return _.get(this, 'discoveryModel.discoveryModel.name', '');
    },
    problemAnnotationAndMarkers() {
      return discovery.annotationsAndMarkers(
        this.discoveryProblems,
        this.discoveryModelString,
      );
    },
    problemAnnotations() {
      // Trigger recalculation of problemMarkers
      this.problemMarkers; // eslint-disable-line
      return this.problemAnnotationAndMarkers.annotations;
    },
    problemMarkers() {
      const { markers } = this.problemAnnotationAndMarkers;
      const editorComponent = this.$children.filter(c => c.editor)[0];
      if (!editorComponent) {
        return markers;
      }

      const { editor } = editorComponent;
      const session = editor.getSession();
      const oldMarkers = session.getMarkers();
      if (oldMarkers) {
        // Bug in Brace editor using wrong Range function means we need to
        // removeMarkers directly here.
        const keys = Object.keys(oldMarkers);
        const errorMarkerIds = keys.filter(k => oldMarkers[k].clazz === 'ace_error-marker');
        errorMarkerIds.forEach(id => session.removeMarker(id));
      }
      if (markers.length > 0) {
        // Bug in Brace editor using wrong Range function means we need to
        // addMarkers directly here, in order to use correct Range function:
        markers.forEach(({
          startRow,
          startCol,
          endRow,
          endCol,
          className,
          type,
          inFront = false,
        }) => {
          const range = new AceRange(startRow, startCol, endRow, endCol);
          session.addMarker(range, className, type, inFront);
        });
      }

      return markers;
    },
  },
  methods: {
    ...mapActions('config', [
      'setDiscoveryModel',
      'validateDiscoveryConfig',
    ]),
    ...mapActions('status', [
      'clearErrors',
    ]),
    /**
     * Validates the Discovery Config.
     */
    async validate() {
      if (this.problems) {
        return false;
      }
      await this.validateDiscoveryConfig();
      if (this.problems) {
        return false;
      }
      return true;
    },
  },
  // Prevent user from progressing FORWARD only if the Discovery Config is invalid.
  // They can navigate backwards, however.
  //
  // "The leave guard is usually used to prevent the user from accidentally leaving the route with unsaved edits. The navigation can be canceled by calling next(false)."
  // See documentation: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
  async beforeRouteLeave(to, from, next) {
    const isBack = from.path === '/wizard/discovery-config'
      && to.path === '/wizard/continue-or-start';
    const isNext = from.path === '/wizard/discovery-config'
      && to.path !== '/wizard/continue-or-start';

    // Always allow user to go back from this page.
    if (isBack) {
      this.clearErrors();
      return next();
    }

    // Allow the user to only go forward if the discovery config is valid
    if (isNext) {
      const valid = await this.validate();
      if (valid) {
        return next();
      }

      return next(false);
    }

    // Neither isBack or isNext is true: If we get into this state something is wrong so just log an error, and prevent navigation.
    // eslint-disable-next-line no-console
    console.error(
      'component=%s, method=beforeRouteLeave: invalid state, vars=%o',
      this.$options.name,
      {
        isBack,
        isNext,
        to,
        from,
      },
    );

    return next(false);
  },
};
</script>

<style scoped>
</style>
