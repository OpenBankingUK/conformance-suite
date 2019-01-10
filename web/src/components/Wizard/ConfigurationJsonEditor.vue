<template>
  <div class="p-3">
    <JsonEditor
      :jsonString="configurationJsonString"
      editorName="configuration-editor"
      setChangeFunctionName=""
    />
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

import JsonEditor from './JsonEditor.vue';

export default {
  name: 'ConfigurationJsonEditor',
  components: {
    JsonEditor,
  },
  computed: {
    ...mapGetters('config', [
      'configurationString',
    ]),
    configurationJsonString() {
      return this.configurationString;
    },
  },
};
</script>

<style scoped>
</style>
