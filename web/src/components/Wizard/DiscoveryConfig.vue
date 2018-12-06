<template>
  <div class="d-flex flex-column flex-fill">
    <h3 class="mb-4">Configuration: API Discovery</h3>
    <div
      v-if="problems"
      class="mb-4">
      <b-alert show>
        Please fix these problems:
        <ul>
          <li
            v-for="(problem, index) in problems"
            :key="index"
          >
            {{ problem }}
          </li>
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
      :value="JSON.stringify(discoveryModel, null, 2)"
      :onChange="onChange"
      :editorProps="{$blockScrolling: Infinity}"
      :focus="true"
      mode="json"
      theme="chrome"
      class="editor mb-4"
      width="100%"
    />
    <b-button-group>
      <b-button
        variant="danger"
        @click="onReset">Reset</b-button>
    </b-button-group>
  </div>
</template>

<script>
import 'brace';
import 'brace/mode/json';
import 'brace/theme/chrome';
import { Ace as AceEditor } from 'vue2-brace-editor';
import { mapGetters, mapActions } from 'vuex';

export default {
  name: 'DiscoveryConfig',
  components: {
    AceEditor,
  },
  props: {
    editorName: {
      type: String,
      private: true,
      default() {
        return 'discovery-config-editor';
      },
    },
  },
  computed: {
    ...mapGetters('config', {
      discoveryModel: 'getDiscoveryModel',
      problems: 'problems',
    }),
  },
  methods: {
    ...mapActions('config', [
      'setDiscoveryModel',
      'resetDiscoveryConfig',
      'validateDiscoveryConfig',
      'setDiscoveryModelProblems',
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
    onReset() {
      this.resetDiscoveryConfig();
      this.resizeEditor();
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
    },
  },
};
</script>

<style>
.editor {
  border: 1px solid lightgrey;
  width: auto !important;
  height: auto !important;
  flex: 1;
}
.problems code {
  max-height: 30vh;
  overflow: scroll;
}
</style>
