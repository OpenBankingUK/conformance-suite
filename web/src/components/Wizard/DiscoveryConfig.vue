<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="panel w-100" style="height:900px">
        <div class="panel-heading">
          <h5>Discovery {Discovery Name}</h5>
        </div>
        <div class="panel-body">
          <div v-if="problems" class="mb-4">
            <b-alert show>Please fix these problems:
              <ul>
                <li v-for="(problem, index) in discoveryProblems" :key="index">{{ problem.error }}</li>
              </ul>
            </b-alert>
          </div>

          <AceEditor
            :ref="editorName"
            :name="editorName"
            :fontSize="12"
            :showPrintMargin="false"
            :showGutter="true"
            :highlightActiveLine="true"
            :value="discoveryModelString"
            :onChange="onChange"
            :editorProps="{$blockScrolling: Infinity}"
            :focus="true"
            :annotations="problemAnnotations"
            mode="json"
            theme="chrome"
            class="editor panel-body"
            width="100%"
          />
        </div>
      </div>
    </div>
  </div>
</template>


<script>
import "brace";
import "brace/mode/json";
import "brace/theme/chrome";
import { Ace as AceEditor } from "vue2-brace-editor";
import { mapGetters, mapActions } from "vuex";
import discovery from "../../api/discovery";

// Bug in Brace editor using wrong Range function means we need to require Range here:
const AceRange = window.ace.acequire("ace/range").Range;

export default {
  name: "DiscoveryConfig",
  components: {
    AceEditor
  },
  props: {
    editorName: {
      type: String,
      private: true,
      default() {
        return "discovery-config-editor";
      }
    }
  },
  computed: {
    ...mapGetters("config", [
      "discoveryModelString",
      "problems",
      "discoveryProblems"
    ]),
    problemAnnotationAndMarkers() {
      return discovery.annotationsAndMarkers(
        this.discoveryProblems,
        this.discoveryModelString
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
        const errorMarkerIds = keys.filter(
          k => oldMarkers[k].clazz === "ace_error-marker"
        );
        errorMarkerIds.forEach(id => session.removeMarker(id));
      }
      if (markers.length > 0) {
        // Bug in Brace editor using wrong Range function means we need to
        // addMarkers directly here, in order to use correct Range function:
        markers.forEach(
          ({
            startRow,
            startCol,
            endRow,
            endCol,
            className,
            type,
            inFront = false
          }) => {
            const range = new AceRange(startRow, startCol, endRow, endCol);
            session.addMarker(range, className, type, inFront);
          }
        );
      }

      return markers;
    }
  },
  methods: {
    ...mapActions("config", [
      "setDiscoveryModel",
      "validateDiscoveryConfig",
      "setDiscoveryModelProblems"
    ]),
    // Gets called by top-level Wizard component in the validateStep function.
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
    onChange(editorString) {
      this.setDiscoveryModel(editorString);
    },
    // Resize the editor to use available space in the parent container.
    // The editor does not dynamically resize itself to fill up available
    // height so this is necessary.
    resizeEditor() {
      const aceEditorComponent = this.$refs[this.editorName];
      this.$nextTick(() => {
        const force = true;
        aceEditorComponent.editor.resize(force);
      });
    }
  },
  // Prevent user from progressing FORWARD only if the Discovery Config is invalid.
  // They can navigate backwards, however.
  //
  // "The leave guard is usually used to prevent the user from accidentally leaving the route with unsaved edits. The navigation can be canceled by calling next(false)."
  // See documentation: https://router.vuejs.org/guide/advanced/navigation-guards.html#in-component-guards
  async beforeRouteLeave(to, from, next) {
    const isBack =
      from.path === "/wizard/discovery-config" &&
      to.path === "/wizard/continue-or-start";
    const isNext =
      from.path === "/wizard/discovery-config" &&
      to.path !== "/wizard/continue-or-start";

    // Always allow user to go back from this page.
    if (isBack) {
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

    // If we get into this state something is wrong so just log an error, and prevent navigation.
    // Neither isBack or isNext is true.
    // eslint-disable-next-line no-console
    console.error(
      "component=%s, method=beforeRouteLeave: invalid state, vars=%o",
      this.$options.name,
      {
        isBack,
        isNext,
        to,
        from
      }
    );

    return next(false);
  }
};
</script>

<style scoped>
.editor {
  border: 1px solid lightgrey;
  width: auto !important;
  flex: 1;
}
.problems code {
  max-height: 30vh;
  overflow: scroll;
}
</style>
