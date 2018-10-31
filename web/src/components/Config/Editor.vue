<template>
  <AceEditor
    :fontSize="14"
    :showPrintMargin="false"
    :showGutter="true"
    :highlightActiveLine="true"
    :name="name"
    :value="JSON.stringify(value, null, 2)"
    :onChange="onChange"
    :editorProps="{$blockScrolling: true}"
    mode="json"
    theme="dracula"
    width="100%"
    class="editor"
    style="height: 300px;"
  />
</template>

<script>
import 'brace';
import 'brace/mode/json';
import 'brace/theme/dracula';
import { Ace as AceEditor } from 'vue2-brace-editor';

export default {
  components: {
    AceEditor,
  },
  props: {
    name: {
      type: String,
      default: 'editor',
    },
    value: {
      type: [Object, Array],
      default: null,
    },
    onChange: {
      type: Function,
      default: () => {},
    },
  },
};
</script>

<style>
.editor {
  margin-top: 8px;
  margin-bottom: 8px;
}
</style>
