<template>
  <div>
    <h2>Discovery Model</h2>
    <editor
      :value="discoveryModelValue"
      :onChange="handleSetDiscoveryModel"
      name="discoveryModel"/>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import Editor from './Editor.vue';

export default {
  components: {
    Editor,
  },
  props: {
    current: {
      type: Number,
      default: 0,
    },
    length: {
      type: Number,
      default: 0,
    },
    discoveryModelValue: {
      type: Object,
      default: () => { },
    },
  },
  computed: {
    ...mapGetters('config', {
      discoveryModel: 'getDiscoveryModel',
    }),
  },
  methods: {
    ...mapActions('config', ['setDiscoveryModel']),
    isValidJSON(json) {
      try {
        JSON.parse(json);
      } catch (e) {
        return false;
      }
      return true;
    },
    handleSetDiscoveryModel(discoveryModel) {
      if (!this.isValidJSON(discoveryModel)) return;
      this.setDiscoveryModel(JSON.parse(discoveryModel));
    },
  },
};
</script>

<style>
</style>
