<template>
  <b-container>
    <b-row>
      <b-col cols="12">
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
      </b-col>
    </b-row>
  </b-container>
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
    isValidJSON(json) {
      try {
        JSON.parse(json);
        this.setDiscoveryModelProblems(null);
      } catch (e) {
        this.setDiscoveryModelProblems([e.message]);
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

<style>
.editor {
  border: 1px solid lightgrey;
}
.problems code {
  max-height: 30vh;
  overflow: scroll;
}
</style>
