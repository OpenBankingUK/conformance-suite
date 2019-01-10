<template>
  <div class="p-3">
    <JsonEditor
      jsonString=""
      editorName="configuration-editor"
      setChangeFunctionName=""
    />
  </div>
</template>

<script>

import JsonEditor from './JsonEditor.vue';

export default {
  name: 'ConfigurationJsonEditor',
  components: {
    JsonEditor,
  },
};
</script>

<style scoped>
</style>
