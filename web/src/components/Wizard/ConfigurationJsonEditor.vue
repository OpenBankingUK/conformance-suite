<template>
  <div
    id="config-editor"
    class="p-3">
    <JsonEditor
      :jsonString="configurationString"
      :wrapEnabled="true"
      editorName="configuration-editor"
      setChangeFunctionName="config/setConfigurationJSON"
    />
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import JsonEditor from './JsonEditor.vue';

const { mapGetters } = createNamespacedHelpers('config');

export default {
  name: 'ConfigurationJsonEditor',
  components: {
    JsonEditor,
  },
  computed: {
    ...mapGetters([
      'configurationString',
    ]),
  },
};
</script>

<style scoped>
#config-editor {
  height: 650px;
}
</style>
