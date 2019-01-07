<template>
  <div>
    <div class="d-flex">
      <div class="flex">
        <NavBar/>
      </div>
      <div class="flex-fill">
        <router-view class="d-flex flex-column flex-fill p-3 router-section"/>
      </div>
    </div>
  </div>
</template>

<script>
import NavBar from './NavBar/NavBar.vue';

export default {
  name: 'Wizard',
  components: {
    NavBar,
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
  min-height: 900px;
  padding-bottom: 0;
}
</style>
