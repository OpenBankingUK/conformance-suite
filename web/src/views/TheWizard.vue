<template>
  <div class="d-flex flex-column flex-fill">
    <div class="d-flex flex-fill">
      <div class="flex">
        <TheNavBar/>
      </div>
      <div class="flex-fill">
        <router-view class="d-flex flex-column flex-fill p-3 router-section"/>
      </div>
    </div>
  </div>
</template>

<script>
import TheNavBar from '../components/TheNavBar.vue';

export default {
  name: 'TheWizard',
  components: {
    TheNavBar,
  },
  data() {
    return {};
  },
  computed: {},
  methods: {},
};
</script>

<style scoped>
.router-section {
  background: #f6f6f6;
  overflow-y: auto;
}

.router-section /deep/ .wizard-step.panel {
  min-height: 850px;
  padding-bottom: 0;
}
</style>
