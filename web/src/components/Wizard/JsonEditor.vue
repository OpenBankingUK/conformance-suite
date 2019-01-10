<template>
  <AceEditor
    :ref="editorName"
    :name="editorName"
    :fontSize="12"
    :showPrintMargin="false"
    :showGutter="true"
    :highlightActiveLine="true"
    :value="jsonString"
    :onChange="onChange"
    :editorProps="{$blockScrolling: Infinity}"
    :focus="true"
    :annotations="problemAnnotations"
    mode="json"
    theme="chrome"
    class="editor panel-body"
    width="100%"
  />
</template>

<script>
import 'brace';
import 'brace/mode/json';
import 'brace/theme/chrome';
import { Ace as AceEditor } from 'vue2-brace-editor';

export default {
  name: 'JsonEditor',
  components: {
    AceEditor,
  },
  props: {
    editorName: {
      type: String,
      required: true,
    },
    jsonString: {
      type: String,
      required: true,
    },
    setChangeFunctionName: {
      type: String,
      required: true,
    },
    problemAnnotations: {
      type: Array,
      required: false,
      default() {
        return null;
      },
    },
  },
  methods: {
    onChange(editorString) {
      this.$store.dispatch(this.setChangeFunctionName, editorString);
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

<style scoped>
.editor {
  border: 1px solid lightgrey;
  width: auto !important;
  flex: 1;
}
</style>
