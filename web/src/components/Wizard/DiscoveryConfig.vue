<template>
  <div class="d-flex flex-column flex-fill">
    <h3>Configuration: API Discovery</h3>
    <div class="d-flex flex-column editor-container flex-fill">
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
        class="editor"
      />
    </div>
    <div
      v-if="problems"
      class="d-flex flex-column problems">
      <h5>Problems</h5>
      <code>{{ problems }}</code>
    </div>
    <b-button-group>
      <b-button
        variant="danger"
        @click="onReset">Reset</b-button>
      <b-button
        variant="primary"
        @click="onValidate">Validate</b-button>
    </b-button-group>
  </div>
</template>

<style>
.editor {
  margin-top: 8px;
  margin-bottom: 8px;
  width: auto !important;
  height: auto !important;
  flex: 1;
}

.problems code {
  max-height: 30vh;
  overflow: scroll;
}
</style>

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
    ]),
    // Gets called by top-level Wizard component in the validateStep function.
    async validate() {
      await this.validateDiscoveryConfig();
      if (this.problems) {
        return Promise.resolve(false);
      }

      return Promise.resolve(true);
    },
    onReset() {
      this.resetDiscoveryConfig();
      this.resizeEditor();
    },
    async onValidate() {
      await this.validateDiscoveryConfig();
      this.resizeEditor();
    },
    isValidJSON(json) {
      try {
        JSON.parse(json);
      } catch (e) {
        return false;
      }

      return true;
    },
    onChange(discoveryModel) {
      if (!this.isValidJSON(discoveryModel)) {
        return;
      }

      this.setDiscoveryModel(JSON.parse(discoveryModel));
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
