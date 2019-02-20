<template>
  <div
    id="app"
    class="d-flex flex-column app">
    <div class="header-section">
      <TheHeader/>
    </div>
    <div class="d-flex flex-column flex-fill">
      <router-view class="d-flex flex-column flex-fill"/>
    </div>
  </div>
</template>

<script>
import TheHeader from './components/TheHeader.vue';

export default {
  name: 'app',
  components: {
    TheHeader,
  },
};
</script>

<style scoped>
@import "~@/assets/css/app.css";
@import "~@/assets/css/fonts.css";

.app {
  font-family: "Metropolis", "Arial", sans-serif;
}

.header-section {
  height: 50px;
}
</style>
