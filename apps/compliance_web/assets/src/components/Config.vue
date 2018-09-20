<template>
  <div>
    <h2>Config</h2>
    <AceEditor
      :fontSize="14"
      :showPrintMargin="false"
      :showGutter="true"
      :highlightActiveLine="true"
      mode="json"
      theme="dracula"
      width="100%"
      name="editor"
      :value="value"
      :onChange="setConfig"
      :editorProps="{$blockScrolling: true}"
      class="editor"
      style="height: 300px;"
    />
    <h2>Payload</h2>
    <AceEditor
      :fontSize="14"
      :showPrintMargin="false"
      :showGutter="true"
      :highlightActiveLine="true"
      mode="json"
      theme="dracula"
      width="100%"
      name="payload"
      :value="payload"
      :onChange="setPayload"
      :editorProps="{$blockScrolling: true}"
      class="editor"
      style="height: 350px;"
    />
    <a-button
      type="primary"
      @click="startValidation"
      class="start_validation"
    >
      Start validation
    </a-button>
    <a-divider />
  </div>
</template>

<script>
import 'brace';
import 'brace/mode/json';
import 'brace/theme/dracula';
import { Ace as AceEditor } from 'vue2-brace-editor';
import { mapState, mapActions } from 'vuex';

export default {
  components: {
    AceEditor,
  },
  computed: {
    ...mapState({
      value: state => state.config.raw,
      payload: state => state.config.payload.raw,
    }),
  },
  methods: {
    ...mapActions('config', ['setConfig', 'setPayload', 'startValidation']),
  },
};
</script>

<style>
.editor {
  margin-bottom: 15px;
}
</style>
