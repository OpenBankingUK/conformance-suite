<template>
  <div class="d-flex flex-column flex-fill">
    <div class="d-flex flex-fill">
      <div
        :class="collapsed ? 'collapsed' : ''"
        class="collapse-icon"
        @click="collapsed = !collapsed"><i
          :class="collapsed ? 'right' : 'left'"
          class="sidebar-arrow"/></div>
      <div
        :class="collapsed ? 'collapsed-sidebar' : ''"
        class="flex">
        <TheNavBar v-if="!collapsed"/>
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
    return {
      collapsed: false,
    };
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

.sidebar-arrow {
  border: solid #fff;
  border-width: 0 2px 2px 0;
  display: inline-block;
  padding: 4px;
}

.sidebar-arrow.left {
  transform: translate(11px, 3px) rotate(135deg);
  -webkit-transform: translate(11px, 3px) rotate(135deg);
}

.sidebar-arrow.right {
  transform: translate(9px, 3px) rotate(-45deg);
  -webkit-transform: translate(9px, 3px) rotate(-45deg);
}

.collapse-icon {
  width: 30px;
  z-index: 1000;
  height: 30px;
  background-color: #6180c3;
  border-radius: 20px;
  position: absolute;
  top: 60px;
  left: 240px;
}

.collapse-icon:hover {
  cursor: pointer;
}

.collapse-icon.collapsed {
  left: 15px;
}

.collapsed-sidebar {
  width: 50px;
}
</style>
