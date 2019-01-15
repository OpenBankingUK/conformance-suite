<template>
  <div
    id="config-editor"
    class="p-3">
    <TheJsonEditor
      :json-string="configurationString"
      :wrap-enabled="true"
      editor-name="configuration-editor"
      set-change-function-name="config/setConfigurationJSON"
    />
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import TheJsonEditor from './TheJsonEditor.vue';

const { mapGetters } = createNamespacedHelpers('config');

export default {
  name: 'TheConfigurationJsonEditor',
  components: {
    TheJsonEditor,
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
