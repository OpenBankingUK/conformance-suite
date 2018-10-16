<template>
  <div id="app">
    <component :is="layout">
      <router-view/>
    </component>
  </div>
</template>

<script>
import NavBar from './components/NavBar.vue';

export default {
  name: 'app',
  components: {
    NavBar,
  },
  computed: {
    layout() {
      return `${this.$route.meta.layout || 'default'}-layout`;
    },
  },
};
</script>
